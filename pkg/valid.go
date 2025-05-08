package pkg

import (
	"regexp"
	"strings"
	"unicode"
)

func NormalizeName(s string) string {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
}

func IsValidName(name string) bool {
	match, _ := regexp.MatchString(`^[\p{L}]+$`, name) // Только буквы (любой алфавит)
	if !match {
		return false
	}
	runes := []rune(name)
	return unicode.IsUpper(runes[0])
}
