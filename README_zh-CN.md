# gorm

*[English](README.md) ∙ [简体中文](README_zh-CN.md)*

[![Github License](https://img.shields.io/github/license/sliveryou/grom.svg?style=flat)](https://github.com/sliveryou/grom/blob/master/LICENSE)
[![Go Doc](https://godoc.org/github.com/sliveryou/grom?status.svg)](https://pkg.go.dev/github.com/sliveryou/grom)
[![Go Report](https://goreportcard.com/badge/github.com/sliveryou/grom)](https://goreportcard.com/report/github.com/sliveryou/grom)
[![Github Latest Release](https://img.shields.io/github/release/sliveryou/grom.svg?style=flat)](https://github.com/sliveryou/grom/releases/latest)
[![Github Latest Tag](https://img.shields.io/github/tag/sliveryou/grom.svg?style=flat)](https://github.com/sliveryou/grom/tags)
[![Github Stars](https://img.shields.io/github/stars/sliveryou/grom.svg?style=flat)](https://github.com/sliveryou/grom/stargazers)

Grom 是一个可以将 mysql 的表字段转换为 golang 的模型结构的命令行工具。
它的全称是 golang relational object mapping（GROM，Golang 关系对象映射）。

## 安装

使用如下命令下载并安装包：

```sh
$ go get -u github.com/sliveryou/grom
```

如果要从源码开始构建的话，需要有 [Go](https://golang.org/dl/) 环境（1.14 及以上版本），并使用如下命令：

```shell
$ git clone https://github.com/sliveryou/grom.git
$ cd grom
$ sh scripts/build.sh
```

或者从 github 的 [release](https://github.com/sliveryou/grom/releases) 页面下载预编译好的二进制文件。

## grom cli

```sh
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
  grom convert -H localhost -P 3306 -u user -p password -d database -t table -e FIELD_COMMENT,JSON_TAG,GORM_TAG --package PACKAGE_NAME --struct STRUCT_NAME

标记:
  -d, --database string   将要连接的 mysql 数据库
  -e, --enable strings    启用的服务（必须包含在 [FIELD_COMMENT,SQL_NULL,GUREGU_NULL,JSON_TAG,XML_TAG,GORM_TAG,XORM_TAG,BEEGO_TAG,GOROSE_TAG] 之中）
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

|   标签    | 表名 | 表 normal 索引 | 表 unique索引 |
|-----------|------|----------------|---------------|
|   json    |   ×  |       ×        |       ×       |
|   xml     |   ×  |       ×        |       ×       |
|   gorm    |   √  |       ×        |       ×       |
|   xorm    |   √  |       ×        |       ×       |
| beego orm |   √  |       √        |       √       |
|  gorose   |   √  |       ×        |       ×       |
