package command

import (
	"strings"
	"unicode"
)

func isExported(name string) bool {
	for _, c := range name {
		if unicode.IsUpper(c) {
			return true
		}
		return false
	}
	return false
}

func CreateFlagName(name string) string {
	var buf strings.Builder
	lastUpper := true
	for _, c := range name {
		upper := unicode.IsUpper(c)
		if upper {
			c = unicode.ToLower(c)
		}
		if !lastUpper && upper {
			buf.WriteString("-")
		}
		lastUpper = upper
		buf.WriteString(string(c))
	}
	return buf.String()
}

func Quote(s string) string {
	chars := "\\\""
	var buf strings.Builder
	start := 0
	length := len(s)
	buf.WriteString(`"`)
	for start < length {
		i := strings.IndexAny(s[start:], chars)
		if i >= 0 {
			buf.WriteString(s[start:i])
			buf.WriteString("\\")
			buf.WriteString(s[i : i+1])
			start = i + 1
		} else {
			break
		}
	}
	if start < length {
		buf.WriteString(s[start:])
	}
	buf.WriteString(`"`)
	return buf.String()
}
