#!/bin/sh
host="localhost"
port="3306"
user="user"
password="password"
database="database"
table="table"

mysql -h${host} -P${port} -u${user} -p${password} -e \
    "SELECT NON_UNIQUE, INDEX_NAME, SEQ_IN_INDEX, COLUMN_NAME, INDEX_COMMENT
    FROM INFORMATION_SCHEMA.STATISTICS
    WHERE TABLE_SCHEMA = '${database}' AND TABLE_NAME = '${table}'
    ORDER BY INDEX_NAME, SEQ_IN_INDEX"

if [ $? -ne 0 ]; then
    echo "get index info failed"
else
    echo "get index info success"
fi
