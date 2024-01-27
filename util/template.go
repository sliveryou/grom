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

	outTplName    = "out"
	gormTplName   = "grom"
	xormTplName   = "xorm"
	beegoTplName  = "beego"
	gormV2TplName = "gormV2"

	//go:embed tpl/out.tpl
	outTpl string
	//go:embed tpl/gorm.tpl
	gormTpl string
	//go:embed tpl/xorm.tpl
	xormTpl string
	//go:embed tpl/beego.tpl
	beegoTpl string
	//go:embed tpl/gormv2.tpl
	gormV2Tpl string
)

func init() {
	var err error
	generator, err = template.New(outTplName).Parse(outTpl)
	if err != nil {
		log.Fatalln(color.Red.Render("parse out.tpl err:", err))
	}
	generator, err = generator.New(gormTplName).Parse(gormTpl)
	if err != nil {
		log.Fatalln(color.Red.Render("parse gorm.tpl err:", err))
	}
	generator, err = generator.New(xormTplName).Parse(xormTpl)
	if err != nil {
		log.Fatalln(color.Red.Render("parse xorm.tpl err:", err))
	}
	generator, err = generator.New(beegoTplName).Funcs(
		template.FuncMap{"getBeegoType": getBeegoType}).Parse(beegoTpl)
	if err != nil {
		log.Fatalln(color.Red.Render("parse beego.tpl err:", err))
	}
	generator, err = generator.New(gormV2TplName).Parse(gormV2Tpl)
	if err != nil {
		log.Fatalln(color.Red.Render("parse gormv2.tpl err:", err))
	}
}

// generateCode generates the output code by command config and structure fields.
func generateCode(cc *CmdConfig, fields []*StructField) (string, error) {
	buffer := &bytes.Buffer{}
	err := generator.ExecuteTemplate(buffer, outTplName, struct {
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
		EnableSQLNull      bool
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
		NeedImport:         cc.EnableGoTime || cc.EnableSQLNull || cc.EnableGureguNull,
		EnableGoTime:       cc.EnableGoTime,
		EnableSQLNull:      cc.EnableSQLNull,
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
