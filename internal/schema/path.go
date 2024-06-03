package schema

import (
	"manala/internal/path"
	"regexp"

	"github.com/xeipuuv/gojsonschema"
)

var fieldPathRegex = regexp.MustCompile(`\.(\d+)`)

func FieldPath(field string) path.Path {
	if field == gojsonschema.STRING_CONTEXT_ROOT {
		field = ""
	}

	// Index: foo.0 -> foo[0]
	field = fieldPathRegex.ReplaceAllString(field, "[${1}]")

	return path.Path(field)
}
