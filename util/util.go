package util

import (
	"strings"
	"unicode"
)

// convertName converts the name to camel case name.
func convertName(name string) string {
	if name == "" {
		return ""
	}

	var cn string
	s := strings.Split(name, "_")

	for _, v := range s {
		upperV := strings.ToUpper(v)
		if _, ok := abbreviation[upperV]; ok {
			cn += upperV
		} else {
			if runesV := []rune(v); len(runesV) > 0 {
				for i, r := range runesV {
					if i == 0 {
						runesV[i] = unicode.ToUpper(r)
					} else {
						runesV[i] = unicode.ToLower(r)
					}
				}
				cn += string(runesV)
			}
		}
	}

	return cn
}
