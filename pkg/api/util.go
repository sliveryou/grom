package api

import (
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"

	"github.com/sliveryou/grom/util"
)

var (
	re = regexp.MustCompile(`\d+`)
	// https://github.com/golang/lint/blob/master/lint.go#L770
	initialisms         = []string{"API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS", "QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SSH", "TLS", "TTL", "UID", "UI", "UUID", "URI", "URL", "UTF8", "VM", "XML", "XSRF", "XSS"}
	initialismsReplacer *strings.Replacer
)

func init() {
	initialismsForReplacer := make([]string, 0, len(initialisms))
	for _, initialism := range initialisms {
		initialismsForReplacer = append(initialismsForReplacer, strings.Title(strings.ToLower(initialism)), initialism)
	}
	initialismsReplacer = strings.NewReplacer(initialismsForReplacer...)
}

// convertComment converts comment.
func convertComment(d string, flag bool) string {
	left, right := d, ""
	start := strings.Index(d, "（")
	end := strings.LastIndex(d, "）")
	if start != -1 && end != -1 {
		left, right = d[:start], d[start+3:end]
	}

	if flag {
		return strings.TrimSpace(left)
	}
	return strings.TrimSpace(right)
}

// getEnums get enums contained in the comment.
func getEnums(s string) string {
	right := convertComment(s, false)
	if right != "" {
		if matches := re.FindAllString(right, -1); len(matches) > 0 {
			return strings.Join(matches, " ")
		}
	}

	return ""
}

// contains reports whether v is present in s.
func contains(s []string, v string) bool {
	for _, n := range s {
		if n == v {
			return true
		}
	}

	return false
}

// mkdirIfNotExist makes directories if the input path is not exists.
func mkdirIfNotExist(dir string) error {
	if dir == "" {
		return nil
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return errors.WithMessage(os.MkdirAll(dir, os.ModePerm), "os.MkdirAll err")
	}

	return nil
}

// writeFile writes the content to the named file.
func writeFile(name, content string) error {
	if err := mkdirIfNotExist(path.Dir(name)); err != nil {
		return errors.WithMessage(err, "mkdirIfNotExist err")
	}

	if err := os.WriteFile(name, []byte(content), writeFilePerm); err != nil {
		return errors.WithMessage(err, "os.WriteFile err")
	}

	return nil
}

// getTypeEmptyString gets the type empty value string.
func getTypeEmptyString(t string) string {
	switch t {
	case "":
		return ""
	case util.GoString:
		return "\"\""
	case util.GoInt, util.GoInt32, util.GoInt64, util.GoFloat32, util.GoFloat64:
		return "0"
	case util.GoUint, util.GoUint32, util.GoUint64:
		return "0"
	case util.GoBool:
		return "false"
	case util.GoBytes, util.GoPointerTime:
		return "nil"
	case util.GoTime:
		return "time.Now()"
	}

	return "nil"
}

// isDefaultEmpty reports whether v is default empty value.
func isDefaultEmpty(v, t string) bool {
	if v == "" {
		return true
	}

	switch t = strings.TrimLeft(t, "*"); t {
	case "":
		return true
	case util.GoString:
		return v == ""
	case util.GoInt, util.GoInt32, util.GoInt64, util.GoFloat32, util.GoFloat64:
		return v == "0"
	case util.GoUint, util.GoUint32, util.GoUint64:
		return v == "0"
	case util.GoBool:
		return v == "false" || v == "0"
	case util.GoTime, util.GoPointerTime:
		return v == "CURRENT_TIMESTAMP"
	}

	return true
}

// toPointer makes the type t to pointer type.
func toPointer(t string) string {
	if t == "" {
		return ""
	}

	if isReferenceType(t) {
		return t
	}

	return "*" + strings.TrimPrefix(t, "*")
}

// isReferenceType reports whether t is reference type.
func isReferenceType(t string) bool {
	t = strings.TrimLeft(t, "*")

	return strings.HasPrefix(t, "map") ||
		strings.HasPrefix(t, "[]")
}

// isPointerWhenUpdated reports whether f is pointer type when updated.
func isPointerWhenUpdated(f StructField) bool {
	if f.IsNullable || !isDefaultEmpty(f.Default, f.Type) ||
		f.Type == util.GoInt32 || f.Type == util.GoBool ||
		(f.Type == util.GoInt && f.Enums != "") {
		return true
	}

	return false
}

// toAbbr converts the string s to abbreviation.
func toAbbr(s string) string {
	s = strcase.ToCamel(strings.TrimSpace(s))
	n := strings.Builder{}

	for _, v := range []byte(s) {
		vIsCap := v >= 'A' && v <= 'Z'
		if vIsCap {
			v += 'a'
			v -= 'A'
			n.WriteByte(v)
		}
	}

	return n.String()
}
