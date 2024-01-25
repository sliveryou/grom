package api

import (
	"os"
	"regexp"
	"strings"
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
		return os.MkdirAll(dir, os.ModePerm)
	}

	return nil
}
