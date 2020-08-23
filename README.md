# grom

*[English](README.md) ∙ [简体中文](README_zh-CN.md)*

[![Github License](https://img.shields.io/github/license/sliveryou/grom.svg?style=flat)](https://github.com/sliveryou/grom/blob/master/LICENSE)
[![Go Doc](https://godoc.org/github.com/sliveryou/grom?status.svg)](https://pkg.go.dev/github.com/sliveryou/grom)
[![Go Report](https://goreportcard.com/badge/github.com/sliveryou/grom)](https://goreportcard.com/report/github.com/sliveryou/grom)
[![Github Latest Release](https://img.shields.io/github/release/sliveryou/grom.svg?style=flat)](https://github.com/sliveryou/grom/releases/latest)
[![Github Latest Tag](https://img.shields.io/github/tag/sliveryou/grom.svg?style=flat)](https://github.com/sliveryou/grom/tags)
[![Github Stars](https://img.shields.io/github/stars/sliveryou/grom.svg?style=flat)](https://github.com/sliveryou/grom/stargazers)

Grom is a powerful command line tool that can convert mysql table fields to golang model structure. 
Its full name is golang relational object mapping (GROM).

## Installation

Download package by using:

```shell script
$ go get -u github.com/sliveryou/grom
```

To build from source code, you need [Go environment](https://golang.org/dl/) (1.14 or newer) and use the following commands:

```shell script
$ git clone https://github.com/sliveryou/grom.git
$ cd grom
$ sh scripts/build.sh
```

Or download a pre-compiled binary from the [release page](https://github.com/sliveryou/grom/releases).

## Grom CLI

```shell script
$ grom -h
Get golang model structure by mysql information schema

Usage:
  grom [command]

Examples:
  grom generate -n ./grom.json
  grom convert -n ./grom.json

Available Commands:
  convert     Convert mysql table fields to golang model structure
  generate    Generate grom configuration file
  help        Help about any command
  version     Show the grom version information

Flags:
  -h, --help   help for grom

Use "grom [command] --help" for more information about a command.

$ grom generate -h
Generate grom configuration file like this:
{
    "host": "localhost",
    "port": 3306,
    "user": "user",
    "password": "password",
    "database": "database",
    "table": "table",
    "package_name": "package_name",
    "struct_name": "struct_name",
    "enable_field_comment": true,
    "enable_sql_null": false,
    "enable_guregu_null": false,
    "enable_json_tag": true,
    "enable_xml_tag": false,
    "enable_gorm_tag": true,
    "enable_xorm_tag": false,
    "enable_beego_tag": false,
    "enable_gorose_tag": false
}

Usage:
  grom generate [flags]

Examples:
  grom generate -n ./grom.json

Flags:
  -h, --help          help for generate
  -n, --name string   the name of the generated grom configuration file (default "grom.json")

$ grom convert -h
Convert mysql table fields to golang model structure by information_schema.columns and information_schema.statistics

Usage:
  grom convert [flags]

Examples:
  grom convert -n ./grom.json
  grom convert -H localhost -P 3306 -u user -p password -d database -t table -e FIELD_COMMENT,JSON_TAG,GORM_TAG --package PACKAGE_NAME --struct STRUCT_NAME

Flags:
  -d, --database string   the database of mysql
  -e, --enable strings    enable services (must in [FIELD_COMMENT,SQL_NULL,GUREGU_NULL,JSON_TAG,XML_TAG,GORM_TAG,XORM_TAG,BEEGO_TAG,GOROSE_TAG])
  -h, --help              help for convert
  -H, --host string       the host of mysql
  -n, --name string       the name of the grom configuration file
      --package string    the package name of the converted model structure
  -p, --password string   the password of mysql
  -P, --port int          the port of mysql
      --struct string     the struct name of the converted model structure
  -t, --table string      the table of mysql
  -u, --user string       the user of mysql

$ grom version -h
Show the grom version information, such as project name, project version, go version, git commit id, build time, etc

Usage:
  grom version [flags]

Flags:
  -h, --help   help for version
```

## Supported Generated Types And Tags

Types:

- [x] [sql](https://godoc.org/database/sql#NullBool)
- [x] [null](https://godoc.org/github.com/guregu/null#Bool)

Tags:

- [x] json
- [x] xml
- [x] [gorm](https://gorm.io/docs/models.html)
- [x] [xorm](https://gobook.io/read/gitea.com/xorm/manual-en-US/chapter-02/4.columns.html)
- [x] [beego orm](https://beego.me/docs/mvc/model/models.md)
- [x] [gorose](https://www.kancloud.cn/fizz/gorose-2/1135839)

## Supported Tag Generation Rules

|    Tag    | PrimaryKey | AutoIncrement | ColumnName | Type | IsNullable | Indexes | Uniques | Default | Comment | ForeignKey |
|-----------|------------|---------------|------------|------|------------|---------|---------|---------|---------|------------|
|   json    |      ×     |       ×       |      √     |  ×   |      ×     |    ×    |    ×    |    ×    |    ×    |      ×     |
|   xml     |      ×     |       ×       |      √     |  ×   |      ×     |    ×    |    ×    |    ×    |    ×    |      ×     |
|   gorm    |      √     |       √       |      √     |  √   |      √     |    √    |    √    |    √    |    √    |      ×     |
|   xorm    |      √     |       √       |      √     |  √   |      √     |    √    |    √    |    √    |    √    |      ×     |
| beego orm |      √     |       √       |      √     |  √   |      √     |    ×    |    ×    |    √    |    √    |      ×     |
|  gorose   |      ×     |       ×       |      √     |  ×   |      ×     |    ×    |    ×    |    ×    |    ×    |      ×     |

## Supported Function Generation Rules

|    Tag    | TableName | TableIndex | TableUnique |
|-----------|-----------|------------|-------------|
|   json    |     ×     |     ×      |      ×      |
|   xml     |     ×     |     ×      |      ×      |
|   gorm    |     √     |     ×      |      ×      |
|   xorm    |     √     |     ×      |      ×      |
| beego orm |     √     |     √      |      √      |
|  gorose   |     √     |     ×      |      ×      |

## Usage Example

1. Create the table named api by following sql: 

```mysql
CREATE TABLE `api`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '接口id',
  `path` varchar(255) NULL DEFAULT NULL COMMENT '接口路径',
  `description` varchar(255) NULL DEFAULT NULL COMMENT '接口描述',
  `group` varchar(255) NULL DEFAULT NULL COMMENT '接口属组',
  `method` varchar(255) NULL DEFAULT 'POST' COMMENT '接口方法',
  `create_time` bigint(20) NULL DEFAULT NULL COMMENT '创建时间',
  `update_time` bigint(20) NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `path_method`(`path`, `method`),
  INDEX `group`(`group`)
) ENGINE = InnoDB AUTO_INCREMENT = 1;
```

2. Generate the grom configuration file and edit it:

```shell script
$ grom generate -n grom.json 
$ vim grom.json
{
    "host": "localhost",
    "port": 3306,
    "user": "user",
    "password": "password",
    "database": "database",
    "table": "api",
    "package_name": "model",
    "struct_name": "API",
    "enable_field_comment": true,
    "enable_sql_null": false,
    "enable_guregu_null": false,
    "enable_json_tag": true,
    "enable_xml_tag": false,
    "enable_gorm_tag": true,
    "enable_xorm_tag": false,
    "enable_beego_tag": false,
    "enable_gorose_tag": false
}
$ grom convert -n grom.json
```

You can also fill in the parameters on the command line without generating a configuration file:

```shell script
$ grom convert -H localhost -P 3306 -u user -p password -d database -t api -e FIELD_COMMENT,JSON_TAG,GORM_TAG
```

Then you will get the generated code:

```go
package model

type API struct {
        ID          int    `json:"id" gorm:"primary_key;column:id;type:int(11) auto_increment;comment:'接口id'"`                           // 接口id
        Path        string `json:"path" gorm:"column:path;type:varchar(255);unique_index:path_method;comment:'接口路径'"`                    // 接口路径
        Description string `json:"description" gorm:"column:description;type:varchar(255);comment:'接口描述'"`                               // 接口描述
        Group       string `json:"group" gorm:"column:group;type:varchar(255);index:group;comment:'接口属组'"`                               // 接口属组
        Method      string `json:"method" gorm:"column:method;type:varchar(255);unique_index:path_method;default:'POST';comment:'接口方法'"` // 接口方法
        CreateTime  int64  `json:"create_time" gorm:"column:create_time;type:bigint(20);comment:'创建时间'"`                                 // 创建时间
        UpdateTime  int64  `json:"update_time" gorm:"column:update_time;type:bigint(20);comment:'更新时间'"`                                 // 更新时间
}

// TableName returns the table name of the API model
func (a *API) TableName() string {
        return "api"
}
```

3. Enjoy yourself!
