#!/bin/sh
host="localhost"
port="3306"
user="user"
password="password"
database="database"
table="table"

mysql -h${host} -P${port} -u${user} -p${password} -e \
  "SELECT COLUMN_NAME, ORDINAL_POSITION, COLUMN_DEFAULT, IS_NULLABLE, DATA_TYPE, CHARACTER_MAXIMUM_LENGTH,
    NUMERIC_PRECISION, NUMERIC_SCALE, COLUMN_TYPE, COLUMN_KEY, EXTRA, COLUMN_COMMENT
    FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = '${database}' AND TABLE_NAME = '${table}'
    ORDER BY ORDINAL_POSITION"

if [ $? -ne 0 ]; then
  echo "get column info failed"
else
  echo "get column info success"
fi
