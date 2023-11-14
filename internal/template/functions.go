package template

import (
	goYaml "github.com/goccy/go-yaml"
	"manala/internal/serrors"
	"reflect"
	"strings"
	textTemplate "text/template"
)

func FuncMap(template *textTemplate.Template) textTemplate.FuncMap {
	return textTemplate.FuncMap{
		"toYaml":  functionToYaml,
		"include": functionInclude(template),
	}
}

// As seen in helm
func functionToYaml(value any) string {
	marshalOptions := []goYaml.EncodeOption{
		goYaml.Indent(4),
		goYaml.UseSingleQuote(true),
		goYaml.UseLiteralStyleIfMultiline(true),
	}

	// Root sequences should not be indented
	// see: https://github.com/goccy/go-yaml/issues/287
	if reflect.ValueOf(value).Kind() != reflect.Slice {
		marshalOptions = append(marshalOptions, goYaml.IndentSequence(true))
	}

	data, err := goYaml.MarshalWithOptions(value, marshalOptions...)
	if err != nil {
		// Swallow errors inside a template.
		return ""
	}
	return strings.TrimSuffix(string(data), "\n")
}

// As seen in helm
func functionInclude(template *textTemplate.Template) func(name string, data any) (string, error) {
	includedNames := make(map[string]int)
	return func(name string, data any) (string, error) {
		var buf strings.Builder
		if v, ok := includedNames[name]; ok {
			if v > 1000 {
				return "", serrors.New("rendering template has a nested reference").
					WithArguments("reference", name)
			}
			includedNames[name]++
		} else {
			includedNames[name] = 1
		}
		err := template.ExecuteTemplate(&buf, name, data)
		return buf.String(), err
	}
}
