package path

import (
	"strconv"
	"strings"
)

var (
	fromJSONPointerReplacer = strings.NewReplacer("~1", "/", "~0", "~")
	toJSONPointerReplacer   = strings.NewReplacer("~", "~0", "/", "~1")
)

func FromJSONPointer(pointer string) string {
	var b strings.Builder
	b.WriteByte('$')
	if pointer != "" && pointer != "/" {
		for token := range strings.SplitSeq(strings.TrimPrefix(pointer, "/"), "/") {
			token = fromJSONPointerReplacer.Replace(token)
			if _, err := strconv.Atoi(token); err == nil {
				b.WriteByte('[')
				b.WriteString(token)
				b.WriteByte(']')
			} else {
				b.WriteByte('.')
				b.WriteString(token)
			}
		}
	}
	return b.String()
}

func ToJSONPointer(path string) string {
	if path == "" || path == "$" {
		return ""
	}
	normalized := strings.ReplaceAll(strings.ReplaceAll(path[1:], "[", "."), "]", "")
	var b strings.Builder
	for token := range strings.SplitSeq(normalized, ".") {
		if token == "" {
			continue
		}
		b.WriteByte('/')
		b.WriteString(toJSONPointerReplacer.Replace(token))
	}
	return b.String()
}
