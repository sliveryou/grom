package api

import (
	"bytes"
	_ "embed"
	stderrors "errors"
	"fmt"
	"go/format"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/gookit/color"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
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
	modelDir      = "model"

	writeFilePerm           = 0o666
	unsignedPrefix          = "u"
	commentPrefix           = "// "
	autoTimeSuffix          = "_at"
	apiFileSuffix           = ".api"
	goFileSuffix            = ".go"
	boolTypeEnums           = "0 1"
	defaultCurrentTimestamp = "CURRENT_TIMESTAMP"
	defaultIdComment        = "ID"
	deleteAt                = "delete_at"
	gormDeleteAt            = "gorm.DeletedAt"
	dataTypeJSON            = "json"
	dataTypeMap             = "map[string]interface{}"
	dataTypesJSON           = "datatypes.JSON"
)

const (
	// RouteStyleSnake snake route style.
	RouteStyleSnake = "snake"
	// RouteStyleKebab kebab route style.
	RouteStyleKebab = "kebab"
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

	errDBConfig         = stderrors.New("invalid db config")
	errEmptyServiceName = stderrors.New("service name can not be empty")
	errEmptyDir         = stderrors.New("dir can not be empty")
	errNoTables         = stderrors.New("there are no tables")
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
func GenerateProject(pc ProjectConfig) error {
	if err := pc.Check(); err != nil {
		return errors.WithMessage(err, "Check err")
	}

	var cab, crb, umb strings.Builder
	if err := mkdirIfNotExist(path.Join(pc.Dir, modelDir)); err != nil {
		return errors.WithMessage(err, "mkdirIfNotExist err")
	}
	defer util.CloseDB()

	apiImports := make([]string, 0, len(pc.Tables))
	for _, table := range pc.Tables {
		if table == "" {
			continue
		}

		var apiName, modelName string
		c := pc.Config
		c.Table = table
		if pc.EnableTrimTablePrefix {
			c.StructName = strcase.ToCamel(strings.TrimPrefix(table, pc.TablePrefix))
			baseName := strings.ToLower(c.StructName)
			apiName = baseName + apiFileSuffix
			modelName = baseName + goFileSuffix
		} else {
			c.StructName = strcase.ToCamel(table)
			c.RouteName = strcase.ToDelimited(strings.TrimPrefix(table, pc.TablePrefix), c.GetDelimiter())
			baseName := strings.ToLower(strcase.ToCamel(c.RouteName))
			apiName = baseName + apiFileSuffix
			modelName = baseName + goFileSuffix
		}

		cc := c.GetCmdConfig()
		fields, err := util.GetFields(cc)
		if err != nil {
			return errors.WithMessage(err, "util.GetFields err")
		}
		if len(fields) == 0 {
			color.Red.Printf("table: %s has no fields, continue to next one\n", table)
			continue
		}

		if c.EnableModel {
			cloneFields := cloneStructFields(cc, fields)
			model, err := util.GenerateCode(cc, cloneFields)
			if err != nil {
				return errors.WithMessage(err, "util.GenerateCode err")
			}
			if err := os.WriteFile(path.Join(pc.Dir, modelDir, modelName), []byte(model), writeFilePerm); err != nil {
				return errors.WithMessage(err, "os.WriteFile err")
			}
		}

		c.UpdateBy(cc)
		api, err := GenerateAPI(c, fields)
		if err != nil {
			return errors.WithMessage(err, "GenerateAPI err")
		}
		if err := os.WriteFile(path.Join(pc.Dir, apiName), []byte(api), writeFilePerm); err != nil {
			return errors.WithMessage(err, "os.WriteFile err")
		}
		apiImports = append(apiImports, apiName)

		ca, err := GenerateConvertAPI(c, fields)
		if err != nil {
			return errors.WithMessage(err, "GenerateConvertAPI err")
		}
		cab.WriteString(ca + "\n\n")

		cr, err := GenerateConvertRPC(c, fields)
		if err != nil {
			return errors.WithMessage(err, "GenerateConvertRPC err")
		}
		crb.WriteString(cr + "\n\n")

		um, err := GenerateUpdateMap(c, fields)
		if err != nil {
			return errors.WithMessage(err, "GenerateUpdateMap err")
		}
		umb.WriteString(um + "\n\n")
	}

	if len(apiImports) > 0 {
		c := pc.Config
		fileName := strings.ToLower(strings.Trim(pc.TablePrefix, `_`))
		if fileName == "" {
			fileName = serverAPIOut
		}
		fileName = path.Join(pc.Dir, fileName+apiFileSuffix)

		out, err := GenerateServerAPI(c, apiImports, fileName)
		if err != nil {
			return errors.WithMessage(err, "GenerateServerAPI err")
		}
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
func GenerateAPI(c Config, fs []*util.StructField) (string, error) {
	gc := getGenerateConfig(c, fs)
	buffer := &bytes.Buffer{}
	routeName := gc.RouteName
	if c.EnablePlural {
		routeName = inflection.Plural(routeName)
	}
	err := generator.ExecuteTemplate(buffer, outTplName, struct {
		TableComment          string
		StructName            string // camel
		RouteName             string // snake or kebab
		GroupName             string // lower
		APIInfo               string
		ServiceName           string
		RoutePrefix           string
		GroupPrefix           string
		IdName                string
		IdNamePlural          string
		IdType                string
		IdComment             string
		IdRawName             string
		IdRawNamePlural       string
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
		RouteName:             routeName,
		GroupName:             gc.GroupName,
		APIInfo:               buildAPIInfo(c),
		ServiceName:           c.ServiceName,
		RoutePrefix:           strings.Trim(c.RoutePrefix, `/`),
		GroupPrefix:           strings.Trim(c.GroupPrefix, `/`),
		IdName:                gc.IdName,
		IdNamePlural:          gc.IdNamePlural,
		IdType:                gc.IdType,
		IdComment:             gc.IdComment,
		IdRawName:             gc.IdRawName,
		IdRawNamePlural:       gc.IdRawNamePlural,
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

// buildAPIInfo builds api info.
func buildAPIInfo(c Config) string {
	b := &strings.Builder{}

	if c.Title != "" {
		b.WriteString(fmt.Sprintf("title: %q\n", c.Title))
	}
	if c.Desc != "" {
		b.WriteString(fmt.Sprintf("desc: %q\n", c.Desc))
	}
	if c.Author != "" {
		b.WriteString(fmt.Sprintf("author: %q\n", c.Author))
	}
	if c.Email != "" {
		b.WriteString(fmt.Sprintf("email: %q\n", c.Email))
	}
	if c.Version != "" {
		b.WriteString(fmt.Sprintf("version: %q\n", c.Version))
	}

	return b.String()
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
		if f.IsPrimaryKey || isReferenceType(f.Type) {
			continue
		}
		tag := fmt.Sprintf("form:\"%s,optional\"", f.RawName)
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
func buildStructUpdateInfo(fs []StructField, isBatchUpdate ...bool) string {
	b := &strings.Builder{}
	isBatch := false
	if len(isBatchUpdate) > 0 {
		isBatch = isBatchUpdate[0]
	}

	for _, f := range fs {
		if IsAutoTimeField(f) {
			continue
		}
		prefix := "json"
		if f.IsPrimaryKey {
			if isBatch {
				continue
			}
			prefix = "path"
		}
		needLabel := false
		tag := fmt.Sprintf("%s:\"%s,optional\"", prefix, f.RawName)
		if !f.IsNullable && f.Default == "" && !isBatch {
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
func GenerateServerAPI(c Config, imports []string, filename ...string) (string, error) {
	buffer := &bytes.Buffer{}
	err := generator.ExecuteTemplate(buffer, serverAPITplName, struct {
		Imports []string
		APIInfo string
	}{
		Imports: imports,
		APIInfo: buildAPIInfo(c),
	})
	if err != nil {
		return "", errors.WithMessage(err, "generator.ExecuteTemplate err")
	}

	api, err := af.APIFormat(buffer.String(), filename...)
	if err != nil {
		return "", errors.WithMessage(err, "format.APIFormat err")
	}

	return api, nil
}

// GenerateConvertAPI generates the output api convert functions by api config and structure fields.
func GenerateConvertAPI(c Config, fs []*util.StructField) (string, error) {
	gc := getGenerateConfig(c, fs)
	buffer := &bytes.Buffer{}

	convertInfo, ifInfo := buildConvertAPIInfo(c.StructName, gc.StructFields)
	err := generator.ExecuteTemplate(buffer, convertAPITplName, struct {
		TableComment string
		StructName   string
		ConvertInfo  string
		IfInfo       string
	}{
		TableComment: c.TableComment,
		StructName:   c.StructName,
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

// buildConvertAPIInfo builds convert api info.
func buildConvertAPIInfo(structName string, fs []StructField) (convertInfo, ifInfo string) {
	var b, ib strings.Builder

	for _, f := range fs {
		if f.Type == dataTypeMap {
			b.WriteString(f.Name + ": make(map[string]interface{}),\n")
			ib.WriteString(fmt.Sprintf("if src.%s != nil {\nif err := json.Unmarshal(src.%s, &dst.%s); err != nil {\nreturn %s{}, errors.WithMessage(err, \"json.Unmarshal %s err\")\n}\n}\n",
				f.Name, f.Name, f.Name, structName, f.Name))
		} else {
			b.WriteString(fmt.Sprintf("%s: src.%s,\n", f.Name, f.Name))
		}
	}

	return strings.TrimSuffix(b.String(), "\n"), strings.TrimSuffix(ib.String(), "\n")
}

// GenerateConvertRPC generates the output rpc convert functions by api config and structure fields.
func GenerateConvertRPC(c Config, fs []*util.StructField) (string, error) {
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
			b.WriteString(f.Name + ": 0,\n")
			ib.WriteString(fmt.Sprintf("if src.%s != nil {\ndst.%s = src.%s.UnixMilli()\n}\n", srcName, f.Name, srcName))
		} else if !f.IsNullable && f.Default != "" {
			b.WriteString(fmt.Sprintf("%s: %s,\n", f.Name, getTypeEmptyString(f.Type)))
			ib.WriteString(fmt.Sprintf("if src.%s != nil {\ndst.%s = *src.%s\n}\n", srcName, f.Name, srcName))
		} else {
			b.WriteString(fmt.Sprintf("%s: src.%s,\n", f.Name, srcName))
		}
	}

	return strings.TrimSuffix(b.String(), "\n"), strings.TrimSuffix(ib.String(), "\n")
}

// GenerateUpdateMap generates the output updateMap by api config and structure fields.
func GenerateUpdateMap(c Config, fs []*util.StructField) (string, error) {
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
			IsDataTypeJSON       bool
		}{
			MemberName:           field.Name,
			MemberRawName:        field.RawName,
			MemberLowerCamelName: strcase.ToLowerCamel(field.Name),
			ObjectName:           strcase.ToLowerCamel(gc.RouteName),
			ObjectMemberName:     initialismsReplacer.Replace(field.Name),
			HasDefault:           field.Default != "",
			IsNullable:           field.IsNullable,
			IsTimeField:          IsTimeField(field),
			IsPointer:            isPointerWhenUpdated(field),
			IsDataTypeJSON:       field.DataType == dataTypeJSON,
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
