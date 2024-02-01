package api

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/gookit/color"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"

	"github.com/sliveryou/goctl/api/protogen"

	af "github.com/sliveryou/grom/pkg/api/format"
	"github.com/sliveryou/grom/util"
)

const (
	outTplName        = "out"
	serverAPITplName  = "serverAPI"
	convertAPITplName = "convertAPI"
	convertRPCTplName = "convertRPC"
	updateMapTplName  = "updateMap"

	convertAPIOut = "convert-api.txt"
	convertRPCOut = "convert-rpc.txt"
	updateMapOut  = "update-map.txt"
	serverAPIOut  = "server"

	writeFilePerm           = 0o666
	unsignedPrefix          = "u"
	commentPrefix           = "// "
	autoTimeSuffix          = "_at"
	apiFileSuffix           = ".api"
	boolTypeEnums           = "0 1"
	defaultCurrentTimestamp = "CURRENT_TIMESTAMP"
	defaultIdComment        = "ID"
)

var (
	generator *template.Template

	//go:embed tpl/out.tpl
	outTpl string
	//go:embed tpl/server-api.tpl
	serverAPITpl string
	//go:embed tpl/convert-api.tpl
	convertAPITpl string
	//go:embed tpl/convert-rpc.tpl
	convertRPCTpl string
	//go:embed tpl/update-map.tpl
	updateMapTpl string
)

func init() {
	var err error
	generator, err = template.New(outTplName).Parse(outTpl)
	if err != nil {
		log.Fatalln(color.Red.Render("parse out.tpl err:", err))
	}
	generator, err = generator.New(serverAPITplName).Parse(serverAPITpl)
	if err != nil {
		log.Fatalln(color.Red.Render("parse server-api.tpl err:", err))
	}
	generator, err = generator.New(convertAPITplName).Parse(convertAPITpl)
	if err != nil {
		log.Fatalln(color.Red.Render("parse convert-api.tpl err:", err))
	}
	generator, err = generator.New(convertRPCTplName).Parse(convertRPCTpl)
	if err != nil {
		log.Fatalln(color.Red.Render("parse convert-rpc.tpl err:", err))
	}
	generator, err = generator.New(updateMapTplName).Parse(updateMapTpl)
	if err != nil {
		log.Fatalln(color.Red.Render("parse update-map.tpl err:", err))
	}
}

// GenerateProject generates the output project by project config.
func GenerateProject(pc *ProjectConfig) error {
	var cab, crb, umb strings.Builder
	if err := mkdirIfNotExist(pc.Dir); err != nil {
		return errors.WithMessage(err, "mkdirIfNotExist err")
	}
	defer util.CloseDB()

	apiImports := make([]string, 0, len(pc.Tables))
	for _, table := range pc.Tables {
		var apiName string
		c := pc.Config
		c.Table = table
		if pc.NeedTrimTablePrefix {
			c.StructName = strcase.ToCamel(strings.TrimPrefix(table, pc.TablePrefix))
			apiName = strings.ToLower(c.StructName) + apiFileSuffix
		} else {
			c.StructName = strcase.ToCamel(table)
			c.SnakeStructName = strcase.ToSnake(strings.TrimPrefix(table, pc.TablePrefix))
			apiName = strings.ToLower(strcase.ToCamel(c.SnakeStructName)) + apiFileSuffix
		}
		apiImports = append(apiImports, apiName)

		cc := c.GetCmdConfig()
		fields, err := util.GetFields(cc)
		if err != nil {
			return errors.WithMessage(err, "util.GetFields err")
		}

		c.UpdateBy(cc)
		api, err := GenerateAPI(&c, fields)
		if err != nil {
			return errors.WithMessage(err, "GenerateAPI err")
		}

		err = os.WriteFile(path.Join(pc.Dir, apiName), []byte(api), writeFilePerm)
		if err != nil {
			return errors.WithMessage(err, "os.WriteFile err")
		}

		ca, err := GenerateConvertAPI(&c, fields)
		if err != nil {
			return errors.WithMessage(err, "GenerateConvertAPI err")
		}
		cab.WriteString(ca + "\n\n")

		cr, err := GenerateConvertRPC(&c, fields)
		if err != nil {
			return errors.WithMessage(err, "GenerateConvertRPC err")
		}
		crb.WriteString(cr + "\n\n")

		um, err := GenerateUpdateMap(&c, fields)
		if err != nil {
			return errors.WithMessage(err, "GenerateUpdateMap err")
		}
		umb.WriteString(um + "\n\n")
	}

	if len(apiImports) > 0 {
		c := pc.Config
		out, err := GenerateServerAPI(&c, apiImports)
		if err != nil {
			return errors.WithMessage(err, "GenerateServerAPI err")
		}
		fileName := strings.ToLower(strings.Trim(pc.TablePrefix, `_`))
		if fileName == "" {
			fileName = serverAPIOut
		}
		fileName = path.Join(pc.Dir, fileName+apiFileSuffix)
		if err := os.WriteFile(fileName, []byte(out), writeFilePerm); err != nil {
			return errors.WithMessage(err, "os.WriteFile err")
		}
		if err := protogen.DoGenProto(fileName, pc.Dir); err != nil {
			return errors.WithMessage(err, "protogen.DoGenProto err")
		}
	}
	if ca := cab.String(); ca != "" {
		if err := os.WriteFile(path.Join(pc.Dir, convertAPIOut), []byte(ca[:len(ca)-1]), writeFilePerm); err != nil {
			return errors.WithMessage(err, "os.WriteFile err")
		}
	}
	if cr := crb.String(); cr != "" {
		if err := os.WriteFile(path.Join(pc.Dir, convertRPCOut), []byte(cr[:len(cr)-1]), writeFilePerm); err != nil {
			return errors.WithMessage(err, "os.WriteFile err")
		}
	}
	if um := umb.String(); um != "" {
		if err := os.WriteFile(path.Join(pc.Dir, updateMapOut), []byte(um[:len(um)-1]), writeFilePerm); err != nil {
			return errors.WithMessage(err, "os.WriteFile err")
		}
	}

	return nil
}

// GenerateAPI generates the output api by api config and structure fields.
func GenerateAPI(c *Config, fs []*util.StructField) (string, error) {
	gc := getGenerateConfig(c, fs)
	buffer := &bytes.Buffer{}
	err := generator.ExecuteTemplate(buffer, outTplName, struct {
		TableComment          string
		StructName            string // camel
		SnakeStructName       string // snake
		GroupName             string // lower
		Title                 string
		Desc                  string
		Author                string
		Email                 string
		Version               string
		ServiceName           string
		ServerPrefix          string
		GroupPrefix           string
		IdName                string
		IdType                string
		IdComment             string
		IdRawName             string
		IdLabel               string
		StructInfo            string
		StructGetInfo         string
		StructCreateInfo      string
		StructUpdateInfo      string
		StructFilterInfo      string
		StructBatchUpdateInfo string
	}{
		TableComment:          c.TableComment,
		StructName:            c.StructName,
		SnakeStructName:       gc.SnakeStructName,
		GroupName:             gc.GroupName,
		Title:                 c.Title,
		Desc:                  c.Desc,
		Author:                c.Author,
		Email:                 c.Email,
		Version:               c.Version,
		ServiceName:           c.ServiceName,
		ServerPrefix:          strings.Trim(c.ServerPrefix, `/`),
		GroupPrefix:           strings.Trim(c.GroupPrefix, `/`),
		IdName:                gc.IdName,
		IdType:                gc.IdType,
		IdComment:             gc.IdComment,
		IdRawName:             gc.IdRawName,
		IdLabel:               convertComment(gc.IdComment, true),
		StructInfo:            buildStructInfo(gc.StructFields),
		StructGetInfo:         buildStructGetInfo(gc.StructFields),
		StructCreateInfo:      buildStructCreateInfo(gc.StructFields),
		StructUpdateInfo:      buildStructUpdateInfo(gc.StructFields),
		StructFilterInfo:      strings.ReplaceAll(buildStructGetInfo(gc.StructFields), "`form:", "`json:"),
		StructBatchUpdateInfo: buildStructUpdateInfo(gc.StructFields, true),
	})
	if err != nil {
		return "", errors.WithMessage(err, "generator.ExecuteTemplate err")
	}

	api, err := af.APIFormat(buffer.String())
	if err != nil {
		return "", errors.WithMessage(err, "format.APIFormat err")
	}

	return api, nil
}

// buildStructInfo builds struct info.
func buildStructInfo(fs []StructField) string {
	b := &strings.Builder{}

	for _, f := range fs {
		field := fmt.Sprintf("\t%s %s `json:%q`", f.Name, f.Type, f.RawName)
		if f.Comment != "" {
			field += commentPrefix + f.Comment
		}
		b.WriteString(field + "\n")
	}

	return strings.TrimSuffix(b.String(), "\n")
}

// buildStructGetInfo builds struct create info.
func buildStructGetInfo(fs []StructField) string {
	b := &strings.Builder{}

	for _, f := range fs {
		if f.IsPrimaryKey {
			continue
		}
		tag := fmt.Sprintf("form:\"%s,optional\"", f.RawName)
		if f.IsNullable {
			f.Type = toPointer(f.Type)
		}
		if contains([]string{util.GoInt, util.GoInt32}, f.Type) && f.Enums != "" {
			f.Type = toPointer(f.Type)
			tag += fmt.Sprintf(" validate:\"omitempty,oneof=%s\" label:%q",
				f.Enums, convertComment(f.Comment, true))
		}
		if contains([]string{util.GoInt32, util.GoBool}, f.Type) {
			f.Type = toPointer(f.Type)
		}
		field := fmt.Sprintf("\t%s %s `%s`", f.Name, f.Type, tag)
		if f.Comment != "" {
			field += commentPrefix + f.Comment
		}
		b.WriteString(field + "\n")
	}

	return strings.TrimSuffix(b.String(), "\n")
}

// buildStructCreateInfo builds struct create info.
func buildStructCreateInfo(fs []StructField) string {
	b := &strings.Builder{}

	for _, f := range fs {
		if f.IsPrimaryKey || IsAutoTimeField(f) {
			continue
		}
		needLabel := false
		tag := fmt.Sprintf("json:\"%s,optional\"", f.RawName)
		if !f.IsNullable && f.Default == "" {
			validate := " validate:\"required\""
			tag = fmt.Sprintf("json:%q", f.RawName)
			if contains([]string{util.GoInt, util.GoInt32}, f.Type) && f.Enums != "" {
				f.Type = toPointer(f.Type)
				validate = fmt.Sprintf(" validate:\"required,oneof=%s\"", f.Enums)
			}
			if contains([]string{util.GoInt32, util.GoBool}, f.Type) {
				f.Type = toPointer(f.Type)
			}
			tag += validate
			needLabel = true
		} else if contains([]string{util.GoInt, util.GoInt32}, f.Type) && f.Enums != "" {
			f.Type = toPointer(f.Type)
			tag += fmt.Sprintf(" validate:\"omitempty,oneof=%s\"", f.Enums)
			needLabel = true
		}
		if f.Default != "" {
			f.Type = toPointer(f.Type)
		}
		if needLabel && f.Comment != "" {
			tag += fmt.Sprintf(" label:%q", convertComment(f.Comment, true))
		}
		field := fmt.Sprintf("\t%s %s `%s`", f.Name, f.Type, tag)
		if f.Comment != "" {
			field += commentPrefix + f.Comment
		}
		b.WriteString(field + "\n")
	}

	return strings.TrimSuffix(b.String(), "\n")
}

// buildStructUpdateInfo builds struct update info.
func buildStructUpdateInfo(fs []StructField, ignorePrimaryKey ...bool) string {
	b := &strings.Builder{}
	ipk := false
	if len(ignorePrimaryKey) > 0 {
		ipk = ignorePrimaryKey[0]
	}

	for _, f := range fs {
		if IsAutoTimeField(f) {
			continue
		}
		prefix := "json"
		if f.IsPrimaryKey {
			if ipk {
				continue
			}
			prefix = "path"
		}
		needLabel := false
		tag := fmt.Sprintf("%s:\"%s,optional\"", prefix, f.RawName)
		if !f.IsNullable && f.Default == "" {
			validate := " validate:\"required\""
			tag = fmt.Sprintf("%s:%q", prefix, f.RawName)
			if contains([]string{util.GoInt, util.GoInt32}, f.Type) && f.Enums != "" {
				f.Type = toPointer(f.Type)
				validate = fmt.Sprintf(" validate:\"required,oneof=%s\"", f.Enums)
			}
			if contains([]string{util.GoInt32, util.GoBool}, f.Type) {
				f.Type = toPointer(f.Type)
			}
			tag += validate
			needLabel = true
		} else {
			f.Type = toPointer(f.Type)
			if contains([]string{util.GoInt, util.GoInt32}, f.Type) && f.Enums != "" {
				tag += fmt.Sprintf(" validate:\"omitempty,oneof=%s\"", f.Enums)
				needLabel = true
			}
		}
		if f.Default != "" {
			f.Type = toPointer(f.Type)
		}
		if needLabel && f.Comment != "" {
			tag += fmt.Sprintf(" label:%q", convertComment(f.Comment, true))
		}
		if !f.IsNullable && f.IsPrimaryKey {
			tag += " swaggerignore:\"true\""
		}
		field := fmt.Sprintf("\t%s %s `%s`", f.Name, f.Type, tag)
		if f.Comment != "" {
			field += commentPrefix + f.Comment
		}
		b.WriteString(field + "\n")
	}

	return strings.TrimSuffix(b.String(), "\n")
}

// GenerateServerAPI generates the output server api by api config and import apis.
func GenerateServerAPI(c *Config, imports []string) (string, error) {
	buffer := &bytes.Buffer{}
	err := generator.ExecuteTemplate(buffer, serverAPITplName, struct {
		Title   string
		Desc    string
		Author  string
		Email   string
		Version string
		Imports []string
	}{
		Title:   c.Title,
		Desc:    c.Desc,
		Author:  c.Author,
		Email:   c.Email,
		Version: c.Version,
		Imports: imports,
	})
	if err != nil {
		return "", errors.WithMessage(err, "generator.ExecuteTemplate err")
	}

	return buffer.String(), nil
}

// GenerateConvertAPI generates the output api convert functions by api config and structure fields.
func GenerateConvertAPI(c *Config, fs []*util.StructField) (string, error) {
	gc := getGenerateConfig(c, fs)
	buffer := &bytes.Buffer{}

	err := generator.ExecuteTemplate(buffer, convertAPITplName, struct {
		TableComment string
		StructName   string
		ConvertInfo  string
	}{
		TableComment: c.TableComment,
		StructName:   c.StructName,
		ConvertInfo:  buildConvertAPIInfo(gc.StructFields),
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

// buildConvertAPIInfo builds convert api info.
func buildConvertAPIInfo(fs []StructField) string {
	b := &strings.Builder{}

	for _, f := range fs {
		field := fmt.Sprintf("%s: src.%s,\n", f.Name, f.Name)
		b.WriteString(field)
	}

	return strings.TrimSuffix(b.String(), "\n")
}

// GenerateConvertRPC generates the output rpc convert functions by api config and structure fields.
func GenerateConvertRPC(c *Config, fs []*util.StructField) (string, error) {
	gc := getGenerateConfig(c, fs)
	buffer := &bytes.Buffer{}

	convertInfo, ifInfo := buildConvertRPCInfo(gc.StructFields)
	err := generator.ExecuteTemplate(buffer, convertRPCTplName, struct {
		TableComment string
		StructName   string
		ModelName    string
		ConvertInfo  string
		IfInfo       string
	}{
		TableComment: c.TableComment,
		StructName:   c.StructName,
		ModelName:    gc.ModelName,
		ConvertInfo:  convertInfo,
		IfInfo:       ifInfo,
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

// buildConvertRPCInfo builds convert rpc info.
func buildConvertRPCInfo(fs []StructField) (convertInfo, ifInfo string) {
	var b, ib strings.Builder

	for _, f := range fs {
		srcName := initialismsReplacer.Replace(f.Name)
		if IsAutoTimeField(f) || IsTimeField(f) {
			b.WriteString(fmt.Sprintf("%s: %s,\n", f.Name, "0"))
			ib.WriteString(fmt.Sprintf("if src.%s != nil {\n\tdst.%s = src.%s.UnixMilli()\n}\n", srcName, f.Name, srcName))
		} else if !f.IsNullable && f.Default != "" {
			b.WriteString(fmt.Sprintf("%s: %s,\n", f.Name, getTypeEmptyString(f.Type)))
			ib.WriteString(fmt.Sprintf("if src.%s != nil {\n\tdst.%s = *src.%s\n}\n", srcName, f.Name, srcName))
		} else {
			b.WriteString(fmt.Sprintf("%s: src.%s,\n", f.Name, srcName))
		}
	}

	return strings.TrimSuffix(b.String(), "\n"), strings.TrimSuffix(ib.String(), "\n")
}

// GenerateUpdateMap generates the output updateMap by api config and structure fields.
func GenerateUpdateMap(c *Config, fs []*util.StructField) (string, error) {
	gc := getGenerateConfig(c, fs)
	buffer := &bytes.Buffer{}
	symbol := strings.Repeat("-", 20)
	buffer.WriteString(fmt.Sprintf("// %s %s %s %s //\n"+
		"// 构建更新map\nupdateMap := make(map[string]interface{})\n",
		symbol, c.StructName, c.TableComment, symbol))

	for _, field := range gc.StructFields {
		if field.IsPrimaryKey || IsAutoTimeField(field) {
			continue
		}
		err := generator.ExecuteTemplate(buffer, updateMapTplName, struct {
			MemberName           string
			MemberRawName        string
			MemberLowerCamelName string
			ObjectName           string // lower camel
			ObjectMemberName     string
			HasDefault           bool
			IsNullable           bool
			IsTimeField          bool
			IsPointer            bool
		}{
			MemberName:           field.Name,
			MemberRawName:        field.RawName,
			MemberLowerCamelName: strcase.ToLowerCamel(field.Name),
			ObjectName:           strcase.ToLowerCamel(gc.SnakeStructName),
			ObjectMemberName:     initialismsReplacer.Replace(field.Name),
			HasDefault:           field.Default != "",
			IsNullable:           field.IsNullable,
			IsTimeField:          IsTimeField(field),
			IsPointer:            isPointerWhenUpdated(field),
		})
		if err != nil {
			return "", errors.WithMessage(err, "generator.ExecuteTemplate err")
		}
	}

	code, err := format.Source(buffer.Bytes())
	if err != nil {
		return "", errors.WithMessage(err, "format.Source err")
	}

	return string(code[:len(code)-1]), nil
}
