# gorm

*[English](README.md) ∙ [简体中文](README_zh-CN.md)*

[![Github License](https://img.shields.io/github/license/sliveryou/grom.svg?style=flat)](https://github.com/sliveryou/grom/blob/master/LICENSE)
[![Go Doc](https://godoc.org/github.com/sliveryou/grom?status.svg)](https://pkg.go.dev/github.com/sliveryou/grom)
[![Go Report](https://goreportcard.com/badge/github.com/sliveryou/grom)](https://goreportcard.com/report/github.com/sliveryou/grom)
[![Github Latest Release](https://img.shields.io/github/release/sliveryou/grom.svg?style=flat)](https://github.com/sliveryou/grom/releases/latest)
[![Github Latest Tag](https://img.shields.io/github/tag/sliveryou/grom.svg?style=flat)](https://github.com/sliveryou/grom/tags)
[![Github Stars](https://img.shields.io/github/stars/sliveryou/grom.svg?style=flat)](https://github.com/sliveryou/grom/stargazers)

Grom 是一个可以将 mysql 的表字段转换为 golang 的模型结构的命令行工具。它的全称是 golang relational object mapping（GROM，Golang 关系对象映射）。

## 安装

使用如下命令下载并安装包：

```shell script
$ go get -u github.com/sliveryou/grom
```

如果要从源码开始构建的话，需要有 [Go](https://golang.org/dl/) 环境（1.14 及以上版本），并使用如下命令：

```shell script
$ git clone https://github.com/sliveryou/grom.git
$ cd grom
$ sh scripts/install.sh
```

或者从 github 的 [release](https://github.com/sliveryou/grom/releases) 页面下载预编译好的二进制文件。

## Grom 命令行接口

```shell script
$ grom -h
利用 mysql 的信息模式获取 golang 的模型结构

用法:
  grom [command]

例子:
  grom generate -n ./grom.json
  grom convert -n ./grom.json

可用命令:
  convert     将 mysql 的表字段转换为 golang 的模型结构
  generate    生成 grom 的配置文件
  help        获取有关任何命令的帮助
  version     显示 grom 的版本信息

标记:
  -h, --help   获取有关 grom 命令的帮助

使用 "grom [command] --help" 获取有关命令的详细信息。

$ grom generate -h
将会生成如下的 grom 配置文件：
{
    "host": "localhost",            // 将要连接的 mysql 主机
    "port": 3306,                   // 将要连接的 mysql 端口
    "user": "user",                 // 将要连接的 mysql 用户
    "password": "password",         // 将要连接的 mysql 密码 
    "database": "database",         // 将要连接的 mysql 数据库
    "table": "table",               // 将要连接的 mysql 数据表
    "package_name": "package_name", // 转换后的模型结构的包名称
    "struct_name": "struct_name",   // 转换后的模型结构的结构体名称
    "enable_initialism": true,      // 是否开启常用缩写词映射
    "enable_field_comment": true,   // 是否启用字段注释
    "enable_sql_null": false,       // 是否启用 sql.Null 类型
    "enable_guregu_null": false,    // 是否启用 null.Null 类型
    "enable_json_tag": true,        // 是否启用 json 标签
    "enable_xml_tag": false,        // 是否启用 xml 标签
    "enable_gorm_tag": true,        // 是否启用 gorm 标签
    "enable_xorm_tag": false,       // 是否启用 xorm 标签
    "enable_beego_tag": false,      // 是否启用 beego orm 标签
    "enable_gorose_tag": false      // 是否启用 gorose 标签
}

用法:
  grom generate [flags]

例子:
  grom generate -n ./grom.json

标记:
  -h, --help          获取有关 generate 命令的帮助
  -n, --name string   生成的 grom 配置文件的名称（默认为 "grom.json"）

$ grom convert -h
通过 information_schema.columns 表和 information_schema.statistics 表，将 mysql 的表字段转换为 golang 的模型结构

用法:
  grom convert [flags]

例子:
  grom convert -n ./grom.json
  grom convert -H localhost -P 3306 -u user -p password -d database -t table -e INITIALISM,FIELD_COMMENT,JSON_TAG,GORM_TAG --package PACKAGE_NAME --struct STRUCT_NAME

标记:
  -d, --database string   将要连接的 mysql 数据库
  -e, --enable strings    启用的服务（必须包含在 [INITIALISM,FIELD_COMMENT,SQL_NULL,GUREGU_NULL,JSON_TAG,XML_TAG,GORM_TAG,XORM_TAG,BEEGO_TAG,GOROSE_TAG] 之中）
  -h, --help              获取有关 convert 命令的帮助
  -H, --host string       将要连接的 mysql 主机
  -n, --name string       指定的 grom 配置文件的名称
      --package string    转换后的模型结构的包名称
  -p, --password string   将要连接的 mysql 密码
  -P, --port int          将要连接的 mysql 端口
      --struct string     转换后的模型结构的结构体名称
  -t, --table string      将要连接的 mysql 数据表
  -u, --user string       将要连接的 mysql 用户

$ grom version -h
显示 grom 版本信息，如项目名称、项目版本、go 版本、git 提交 id 和构建时间等

用法:
  grom version [flags]

标记:
  -h, --help   获取有关 version 命令的帮助
```

## 目前支持生成的类型和标签

类型：

- [x] [sql](https://godoc.org/database/sql#NullBool)
- [x] [null](https://godoc.org/github.com/guregu/null#Bool)

标签：

- [x] json
- [x] xml
- [x] [gorm](https://gorm.io/zh_CN/docs/models.html)
- [x] [xorm](https://gobook.io/read/gitea.com/xorm/manual-zh-CN/chapter-02/4.columns.html)
- [x] [beego orm](https://beego.me/docs/mvc/model/models.md)
- [x] [gorose](https://www.kancloud.cn/fizz/gorose-2/1135839)

## 支持的标签生成规则

|   标签    | 主键 | 自增 | 列名 | 类型 | 是否为 null | normal 索引 | unique 索引 | 默认值 | 注释 | 外键 |
|-----------|------|------|------|------|-------------|-------------|-------------|--------|------|------|
|   json    |   ×  |  ×   |  √   |  ×   |      ×      |      ×      |      ×      |    ×   |  ×   |  ×   |
|   xml     |   ×  |  ×   |  √   |  ×   |      ×      |      ×      |      ×      |    ×   |  ×   |  ×   |
|   gorm    |   √  |  √   |  √   |  √   |      √      |      √      |      √      |    √   |  √   |  ×   |
|   xorm    |   √  |  √   |  √   |  √   |      √      |      √      |      √      |    √   |  √   |  ×   |
| beego orm |   √  |  √   |  √   |  √   |      √      |      ×      |      ×      |    √   |  √   |  ×   |
|  gorose   |   ×  |  ×   |  √   |  ×   |      ×      |      ×      |      ×      |    ×   |  ×   |  ×   |

## 支持的函数生成规则

|   标签    | 表名函数（TableName） | 表 normal 索引函数（TableIndex） | 表 unique 索引函数（TableUnique） |
|-----------|-----------------------|----------------------------------|-----------------------------------|
|   json    |            ×          |                ×                 |                 ×                 |
|   xml     |            ×          |                ×                 |                 ×                 |
|   gorm    |            √          |                ×                 |                 ×                 |
|   xorm    |            √          |                ×                 |                 ×                 |
| beego orm |            √          |                √                 |                 √                 |
|  gorose   |            √          |                ×                 |                 ×                 |

## 用法举例

1. 通过以下 sql 语句创建名为 api 的表：

```mysql
CREATE TABLE `api` (
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

2. 生成并编辑 grom 的配置文件：

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
    "enable_initialism": true,
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

你也可以在命令行中填写参数，而不生成配置文件：

```shell script
$ grom convert -H localhost -P 3306 -u user -p password -d database -t api -e INITIALISM,FIELD_COMMENT,JSON_TAG,GORM_TAG
```

然后你将会得到生成的代码：

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

3. 尽情享受吧。
