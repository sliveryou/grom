package util

import (
	"bytes"
	_ "embed"
	"go/format"
	"log"
	"strings"
	"text/template"

	"github.com/gookit/color"
	"github.com/pkg/errors"
)

var (
	generator *template.Template

	//go:embed tpl/out.tpl
	outTemplate string
	//go:embed tpl/gorm.tpl
	gormTemplate string
	//go:embed tpl/xorm.tpl
	xormTemplate string
	//go:embed tpl/beego.tpl
	beegoTemplate string
	//go:embed tpl/gormv2.tpl
	gormV2Template string
)

func init() {
	var err error
	generator, err = template.New("out").Parse(outTemplate)
	if err != nil {
		log.Fatalln(color.Red.Render("parse out template err:", err))
	}
	generator, err = generator.New("gorm").Parse(gormTemplate)
	if err != nil {
		log.Fatalln(color.Red.Render("parse gorm template err:", err))
	}
	generator, err = generator.New("xorm").Parse(xormTemplate)
	if err != nil {
		log.Fatalln(color.Red.Render("parse xorm template err:", err))
	}
	generator, err = generator.New("beego").Funcs(
		template.FuncMap{"getBeegoType": getBeegoType}).Parse(beegoTemplate)
	if err != nil {
		log.Fatalln(color.Red.Render("parse beego orm template err:", err))
	}
	generator, err = generator.New("gormV2").Parse(gormV2Template)
	if err != nil {
		log.Fatalln(color.Red.Render("parse gormV2 template err:", err))
	}
}

// generateCode generates the output code by command config and structure fields.
func generateCode(cc *CMDConfig, fields []*StructField) (string, error) {
	buffer := &bytes.Buffer{}
	err := generator.ExecuteTemplate(buffer, "out", struct {
		Table              string
		TableComment       string
		PackageName        string
		StructName         string
		ShortStructName    string
		StructFields       []*StructField
		TableIndexes       []string
		TableUniques       []string
		EnableFieldComment bool
		NeedImport         bool
		EnableGoTime       bool
		EnableSqlNull      bool
		EnableGureguNull   bool
		EnableTableName    bool
		EnableTableIndex   bool
		EnableTableUnique  bool
	}{
		Table:              cc.Table,
		TableComment:       cc.TableComment,
		PackageName:        cc.PackageName,
		StructName:         cc.StructName,
		ShortStructName:    strings.ToLower(cc.StructName[0:1]),
		StructFields:       fields,
		TableIndexes:       uniqueStrings(tableIndexes),
		TableUniques:       uniqueStrings(tableUniques),
		EnableFieldComment: cc.EnableFieldComment,
		NeedImport:         cc.EnableGoTime || cc.EnableSqlNull || cc.EnableGureguNull,
		EnableGoTime:       cc.EnableGoTime,
		EnableSqlNull:      cc.EnableSqlNull,
		EnableGureguNull:   cc.EnableGureguNull,
		EnableTableName:    cc.EnableGormTag || cc.EnableXormTag || cc.EnableBeegoTag || cc.EnableGoroseTag || cc.EnableGormV2Tag,
		EnableTableIndex:   cc.EnableBeegoTag && len(tableIndexes) != 0,
		EnableTableUnique:  cc.EnableBeegoTag && len(tableUniques) != 0,
	})
	if err != nil {
		return "", errors.WithMessage(err, "generator.ExecuteTemplate err")
	}

	code, err := format.Source(buffer.Bytes())
	if err != nil {
		return "", errors.WithMessage(err, "format.Source err")
	}

	return string(code[:len(code)-1]), nil
}

// generateTag generates the tag string by column information and tag name.
func generateTag(ci *ColumnInfo, tag string) string {
	buffer := &bytes.Buffer{}
	err := generator.ExecuteTemplate(buffer, tag, ci)
	if err != nil {
		// err just print
		color.Red.Printf("generateTag err: %v, tag: %s, column: %+v\n", err, tag, *ci)
		return ""
	}

	return strings.TrimSpace(buffer.String())
}

// uniqueStrings returns the unique string slice.
func uniqueStrings(slice []string) []string {
	result := make([]string, 0, len(slice))
	uniqueMap := make(map[string]struct{})

	for _, value := range slice {
		if _, ok := uniqueMap[value]; !ok {
			uniqueMap[value] = struct{}{}
			result = append(result, value)
		}
	}

	return result
}
