package util

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gookit/color"
	"github.com/pkg/errors"
)

// CloseDB closes the db connection.
func CloseDB() {
	if err := db.Close(); err != nil {
		color.Red.Println("db.Close err:", err)
	}
}

// getDB returns the opened db connection.
func getDB(c *CMDConfig) (*sql.DB, error) {
	if db != nil {
		return db, nil
	}

	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Database)

	db, err = sql.Open(MySQLDriverName, dsn)
	if err != nil {
		return nil, errors.WithMessage(err, "sql.Open err")
	}

	return db, nil
}

// getTableComment returns the comment of table.
func getTableComment(c *CMDConfig) (string, error) {
	db, err := getDB(c)
	if err != nil {
		return "", err
	}

	querySQL := "SELECT TABLE_COMMENT " +
		"FROM INFORMATION_SCHEMA.TABLES " +
		"WHERE TABLE_NAME = ?"

	rows, err := db.Query(querySQL, c.Table)
	if err != nil {
		return "", errors.WithMessage(err, "db.Query err")
	}

	if rows == nil {
		return "", errors.New("no rows returned")
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			color.Red.Println("rows.Close err:", err)
		}
	}(rows)

	var comment string
	for rows.Next() {
		err = rows.Scan(&comment)
		if err != nil {
			return "", errors.WithMessage(err, "rows.Scan err")
		}
	}

	return comment, nil
}

// getColumnInfos returns the details of columns.
func getColumnInfos(c *CMDConfig) ([]*ColumnInfo, error) {
	db, err := getDB(c)
	if err != nil {
		return nil, err
	}

	querySQL := "SELECT COLUMN_NAME, ORDINAL_POSITION, COLUMN_DEFAULT, IS_NULLABLE, " +
		"DATA_TYPE, CHARACTER_MAXIMUM_LENGTH, NUMERIC_PRECISION, NUMERIC_SCALE, " +
		"COLUMN_TYPE, COLUMN_KEY, EXTRA, COLUMN_COMMENT " +
		"FROM INFORMATION_SCHEMA.COLUMNS " +
		"WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? " +
		"ORDER BY ORDINAL_POSITION"

	rows, err := db.Query(querySQL, c.Database, c.Table)
	if err != nil {
		return nil, errors.WithMessage(err, "db.Query err")
	}

	if rows == nil {
		return nil, errors.New("no rows returned")
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			color.Red.Println("rows.Close err:", err)
		}
	}(rows)

	columnInfos := make([]*ColumnInfo, 0)
	indexInfos, err := getIndexInfos(c)
	if err != nil {
		return nil, errors.WithMessage(err, "getIndexInfos err")
	}

	for rows.Next() {
		var (
			// COLUMN_NAME, IS_NULLABLE, DATA_TYPE, COLUMN_TYPE, COLUMN_KEY, EXTRA, COLUMN_COMMENT
			cn, in, dt, ct, ck, e, cc string
			// ORDINAL_POSITION
			op int
			// COLUMN_DEFAULT
			cd sql.NullString
			// CHARACTER_MAXIMUM_LENGTH, NUMERIC_PRECISION, NUMERIC_SCALE
			cml, np, nc sql.NullInt64
		)

		err = rows.Scan(&cn, &op, &cd, &in, &dt, &cml, &np, &nc, &ct, &ck, &e, &cc)
		if err != nil {
			return nil, errors.WithMessage(err, "rows.Scan err")
		}

		ci := ColumnInfo{
			Name: cn, DataType: dt, Type: ct, Default: strings.TrimSpace(cd.String), Comment: strings.TrimSpace(cc),
			Length: cml.Int64, Precision: np.Int64, Scale: nc.Int64, Position: op,
			IsPrimaryKey: ck == "PRI", IsAutoIncrement: strings.Contains(e, "auto_increment"),
			IsUnsigned: strings.Contains(ct, "unsigned") && !c.DisableUnsigned, IsNullable: in == "YES",
		}

		ci.Indexes, ci.UniqueIndexes = getColumnIndexInfos(indexInfos, ci.Name)
		columnInfos = append(columnInfos, &ci)
	}

	if c.EnableBeegoTag {
		tableIndexes, tableUniques = getTableIndexes(indexInfos, c.EnableInitialism)
	}

	return columnInfos, nil
}

// getIndexInfos returns the details of indexes.
func getIndexInfos(c *CMDConfig) ([]*IndexInfo, error) {
	db, err := getDB(c)
	if err != nil {
		return nil, err
	}

	querySQL := "SELECT NON_UNIQUE, INDEX_NAME, SEQ_IN_INDEX, COLUMN_NAME, INDEX_COMMENT " +
		"FROM INFORMATION_SCHEMA.STATISTICS " +
		"WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? " +
		"ORDER BY INDEX_NAME, SEQ_IN_INDEX"

	rows, err := db.Query(querySQL, c.Database, c.Table)
	if err != nil {
		return nil, errors.WithMessage(err, "db.Query err")
	}

	indexInfos := make([]*IndexInfo, 0)

	if rows == nil {
		return indexInfos, nil
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			color.Red.Println("rows.Close err:", err)
		}
	}(rows)

	for rows.Next() {
		var (
			// NON_UNIQUE, SEQ_IN_INDEX
			nu, sii int
			// INDEX_NAME, COLUMN_NAME, INDEX_COMMENT
			in, cn, ic string
		)

		err = rows.Scan(&nu, &in, &sii, &cn, &ic)
		if err != nil {
			return nil, errors.WithMessage(err, "rows.Scan err")
		}

		if in == "PRIMARY" {
			continue
		}

		ii := IndexInfo{
			Name: in, ColumnName: cn, Comment: ic, Sequence: sii, IsUnique: nu == indexUnique,
		}

		indexInfos = append(indexInfos, &ii)
	}

	return indexInfos, nil
}

// getColumnIndexInfos returns the details of column indexes and column unique indexes.
func getColumnIndexInfos(indexInfos []*IndexInfo, columnName string) (columnIndexes, columnUniques []*IndexInfo) {
	for i := range indexInfos {
		indexInfo := indexInfos[i]
		if indexInfo.ColumnName == columnName {
			if indexInfo.IsUnique {
				columnUniques = append(columnUniques, indexInfo)
			} else {
				columnIndexes = append(columnIndexes, indexInfo)
			}
		}
	}

	return
}

// getTableIndexes returns the details of table indexes and table unique indexes.
func getTableIndexes(indexInfos []*IndexInfo, enableInitialism ...bool) (tableIndexes, tableUniques []string) {
	tableIndexMap, tableUniqueMap := make(map[string][]string), make(map[string][]string)

	for i := range indexInfos {
		indexInfo := indexInfos[i]
		columnName := fmt.Sprintf("%q", convertName(indexInfo.ColumnName, enableInitialism...))
		if indexInfo.IsUnique {
			uniqueIndexes := tableUniqueMap[indexInfo.Name]
			uniqueIndexes = append(uniqueIndexes, columnName)
			tableUniqueMap[indexInfo.Name] = uniqueIndexes
		} else {
			normalIndexes := tableIndexMap[indexInfo.Name]
			normalIndexes = append(normalIndexes, columnName)
			tableIndexMap[indexInfo.Name] = normalIndexes
		}
	}

	if len(tableUniqueMap) != 0 {
		for indexName := range tableUniqueMap {
			columns := tableUniqueMap[indexName]
			tableUniques = append(tableUniques, strings.Join(columns, ","))
		}
	}
	if len(tableIndexMap) != 0 {
		for indexName := range tableIndexMap {
			columns := tableIndexMap[indexName]
			tableIndexes = append(tableIndexes, strings.Join(columns, ","))
		}
	}

	return
}
