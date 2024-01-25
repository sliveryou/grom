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
	generator, err = template.New("out").Parse(outTpl)
	if err != nil {
		log.Fatalln(color.Red.Render("parse out.tpl err:", err))
	}
	generator, err = generator.New("serverAPI").Parse(serverAPITpl)
	if err != nil {
		log.Fatalln(color.Red.Render("parse server-api.tpl err:", err))
	}
	generator, err = generator.New("convertAPI").Parse(convertAPITpl)
	if err != nil {
		log.Fatalln(color.Red.Render("parse convert-api.tpl err:", err))
	}
	generator, err = generator.New("convertRPC").Parse(convertRPCTpl)
	if err != nil {
		log.Fatalln(color.Red.Render("parse convert-rpc.tpl err:", err))
	}
	generator, err = generator.New("updateMap").Parse(updateMapTpl)
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

	var apiImports []string
	for _, table := range pc.Tables {
		c := pc.Config
		c.Table = table
		apiName := ""
		if pc.NeedTrimTablePrefix {
			c.StructName = strcase.ToCamel(strings.TrimPrefix(table, pc.TablePrefix))
			apiName = strings.ToLower(c.StructName) + ".api"
		} else {
			c.StructName = strcase.ToCamel(table)
			c.SnakeStructName = strcase.ToSnake(strings.TrimPrefix(table, pc.TablePrefix))
			apiName = strings.ToLower(strcase.ToCamel(c.SnakeStructName)) + ".api"
		}
		apiImports = append(apiImports, apiName)

		cc := c.GetCMDConfig()
		fields, err := util.GetFields(cc)
		if err != nil {
			return errors.WithMessage(err, "util.GetFields err")
		}

		c.UpdateBy(cc)
		api, err := GenerateAPI(&c, fields)
		if err != nil {
			return errors.WithMessage(err, "GenerateAPI err")
		}

		err = os.WriteFile(path.Join(pc.Dir, apiName), []byte(api), 0o666)
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
		if out, err := GenerateServerAPI(&c, apiImports); err == nil {
			fileName := strings.ToLower(strings.Trim(pc.TablePrefix, `_`))
			if fileName == "" {
				fileName = "server"
			}
			fileName = path.Join(pc.Dir, fileName+".api")
			if err := os.WriteFile(fileName, []byte(out), 0o666); err != nil {
				return errors.WithMessage(err, "os.WriteFile err")
			}
			if err := protogen.DoGenProto(fileName, pc.Dir); err != nil {
				return errors.WithMessage(err, "protogen.DoGenProto err")
			}
		} else {
			return errors.WithMessage(err, "GenerateServerAPI err")
		}
	}
	if ca := cab.String(); ca != "" {
		if err := os.WriteFile(path.Join(pc.Dir, "convert-api.txt"), []byte(ca[:len(ca)-1]), 0o666); err != nil {
			return errors.WithMessage(err, "os.WriteFile err")
		}
	}
	if cr := crb.String(); cr != "" {
		if err := os.WriteFile(path.Join(pc.Dir, "convert-rpc.txt"), []byte(cr[:len(cr)-1]), 0o666); err != nil {
			return errors.WithMessage(err, "os.WriteFile err")
		}
	}
	if um := umb.String(); um != "" {
		if err := os.WriteFile(path.Join(pc.Dir, "update-map.txt"), []byte(um[:len(um)-1]), 0o666); err != nil {
			return errors.WithMessage(err, "os.WriteFile err")
		}
	}

	return nil
}

// GenerateAPI generates the output api by api config and structure fields.
func GenerateAPI(c *Config, fs []*util.StructField) (string, error) {
	gc := getGenerateConfig(c, fs)
	buffer := &bytes.Buffer{}
	err := generator.ExecuteTemplate(buffer, "out", struct {
		TableComment     string
		StructName       string // camel
		SnakeStructName  string // snake
		GroupName        string // lower
		Title            string
		Desc             string
		Author           string
		Email            string
		Version          string
		ServiceName      string
		ServerPrefix     string
		GroupPrefix      string
		IdComment        string
		IdLabel          string
		StructInfo       string
		StructGetInfo    string
		StructCreateInfo string
		StructUpdateInfo string
	}{
		TableComment:     c.TableComment,
		StructName:       c.StructName,
		SnakeStructName:  gc.SnakeStructName,
		GroupName:        gc.GroupName,
		Title:            c.Title,
		Desc:             c.Desc,
		Author:           c.Author,
		Email:            c.Email,
		Version:          c.Version,
		ServiceName:      c.ServiceName,
		ServerPrefix:     strings.Trim(c.ServerPrefix, `/`),
		GroupPrefix:      strings.Trim(c.GroupPrefix, `/`),
		IdComment:        gc.IdComment,
		IdLabel:          convertComment(gc.IdComment, true),
		StructInfo:       BuildStructInfo(gc.StructFields),
		StructGetInfo:    BuildStructGetInfo(gc.StructFields),
		StructCreateInfo: BuildStructCreateInfo(gc.StructFields),
		StructUpdateInfo: BuildStructUpdateInfo(gc.StructFields),
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

// BuildStructInfo builds struct info.
func BuildStructInfo(fs []util.StructField) string {
	b := &strings.Builder{}

	for _, f := range fs {
		field := fmt.Sprintf("\t%s %s `json:%q`", f.Name, f.Type, f.RawName)
		if f.Comment != "" {
			field += "// " + f.Comment
		}
		b.WriteString(field + "\n")
	}

	return strings.TrimSuffix(b.String(), "\n")
}

// BuildStructGetInfo builds struct create info.
func BuildStructGetInfo(fs []util.StructField) string {
	b := &strings.Builder{}

	for _, f := range fs {
		if f.IsPrimaryKey {
			continue
		}
		tag := fmt.Sprintf("form:\"%s,optional\"", f.RawName)
		enums := getEnums(f.Comment)
		if contains([]string{"int", "int32"}, f.Type) && enums != "" {
			f.Type = "*" + f.Type
			tag += fmt.Sprintf(" validate:\"omitempty,oneof=%s\" label:%q",
				enums, convertComment(f.Comment, true))
		}
		if contains([]string{"int32", "bool"}, f.Type) {
			f.Type = "*" + f.Type
		}
		field := fmt.Sprintf("\t%s %s `%s`", f.Name, f.Type, tag)
		if f.Comment != "" {
			field += "// " + f.Comment
		}
		b.WriteString(field + "\n")
	}

	return strings.TrimSuffix(b.String(), "\n")
}

// BuildStructCreateInfo builds struct create info.
func BuildStructCreateInfo(fs []util.StructField) string {
	b := &strings.Builder{}

	for _, f := range fs {
		if f.IsPrimaryKey {
			continue
		}
		needLabel := false
		enums := getEnums(f.Comment)
		tag := fmt.Sprintf("json:\"%s,optional\"", f.RawName)
		if !f.IsNullable {
			validate := " validate:\"required\""
			tag = fmt.Sprintf("json:%q", f.RawName)
			if contains([]string{"int", "int32"}, f.Type) && enums != "" {
				f.Type = "*" + f.Type
				validate = fmt.Sprintf(" validate:\"required,oneof=%s\"", enums)
			}
			if contains([]string{"int32", "bool"}, f.Type) {
				f.Type = "*" + f.Type
			}
			tag += validate
			needLabel = true

		} else if contains([]string{"int", "int32"}, f.Type) && enums != "" {
			f.Type = "*" + f.Type
			tag += fmt.Sprintf(" validate:\"omitempty,oneof=%s\"", enums)
			needLabel = true
		}
		if needLabel && f.Comment != "" {
			tag += fmt.Sprintf(" label:%q", convertComment(f.Comment, true))
		}
		field := fmt.Sprintf("\t%s %s `%s`", f.Name, f.Type, tag)
		if f.Comment != "" {
			field += "// " + f.Comment
		}
		b.WriteString(field + "\n")
	}

	return strings.TrimSuffix(b.String(), "\n")
}

// BuildStructUpdateInfo builds struct update info.
func BuildStructUpdateInfo(fs []util.StructField) string {
	b := &strings.Builder{}

	for _, f := range fs {
		prefix := "json"
		if f.IsPrimaryKey {
			prefix = "path"
		}
		needLabel := false
		enums := getEnums(f.Comment)
		tag := fmt.Sprintf("%s:\"%s,optional\"", prefix, f.RawName)
		if !f.IsNullable {
			validate := " validate:\"required\""
			tag = fmt.Sprintf("%s:%q", prefix, f.RawName)
			if contains([]string{"int", "int32"}, f.Type) && enums != "" {
				f.Type = "*" + f.Type
				validate = fmt.Sprintf(" validate:\"required,oneof=%s\"", enums)
			}
			if contains([]string{"int32", "bool"}, f.Type) {
				f.Type = "*" + f.Type
			}
			tag += validate
			needLabel = true
		} else {
			f.Type = "*" + f.Type
			if contains([]string{"int", "int32"}, f.Type) && enums != "" {
				tag += fmt.Sprintf(" validate:\"omitempty,oneof=%s\"", enums)
				needLabel = true
			}
		}
		if needLabel && f.Comment != "" {
			tag += fmt.Sprintf(" label:%q", convertComment(f.Comment, true))
		}
		if !f.IsNullable && f.IsPrimaryKey {
			tag += " swaggerignore:\"true\""
		}
		field := fmt.Sprintf("\t%s %s `%s`", f.Name, f.Type, tag)
		if f.Comment != "" {
			field += "// " + f.Comment
		}
		b.WriteString(field + "\n")
	}

	return strings.TrimSuffix(b.String(), "\n")
}

// GenerateServerAPI generates the output server api by api config and import apis.
func GenerateServerAPI(c *Config, imports []string) (string, error) {
	buffer := &bytes.Buffer{}
	err := generator.ExecuteTemplate(buffer, "serverAPI", struct {
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

	err := generator.ExecuteTemplate(buffer, "convertAPI", struct {
		TableComment string
		StructName   string
		ConvertInfo  string
	}{
		TableComment: c.TableComment,
		StructName:   c.StructName,
		ConvertInfo:  BuildConvertAPIInfo(gc.StructFields),
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

// BuildConvertAPIInfo builds convert api info.
func BuildConvertAPIInfo(fs []util.StructField) string {
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

	err := generator.ExecuteTemplate(buffer, "convertRPC", struct {
		TableComment string
		StructName   string
		ModelName    string
		ConvertInfo  string
	}{
		TableComment: c.TableComment,
		StructName:   c.StructName,
		ModelName:    gc.ModelName,
		ConvertInfo:  BuildConvertRPCInfo(gc.StructFields),
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

// BuildConvertRPCInfo builds convert rpc info.
func BuildConvertRPCInfo(fs []util.StructField) string {
	b := &strings.Builder{}

	for _, f := range fs {
		field := fmt.Sprintf("%s: src.%s,\n", f.Name, initialismsReplacer.Replace(f.Name))
		b.WriteString(field)
	}

	return strings.TrimSuffix(b.String(), "\n")
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
		if field.IsPrimaryKey {
			continue
		}
		err := generator.ExecuteTemplate(buffer, "updateMap", struct {
			MemberName       string
			MemberRawName    string
			ObjectName       string // lower camel
			ObjectMemberName string
			IsNullable       bool
		}{
			MemberName:       field.Name,
			MemberRawName:    field.RawName,
			ObjectName:       strcase.ToLowerCamel(gc.SnakeStructName),
			ObjectMemberName: initialismsReplacer.Replace(field.Name),
			IsNullable:       true,
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
