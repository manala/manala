package template

import (
	"reflect"
	"strings"
	textTemplate "text/template"

	"github.com/manala/manala/internal/serrors"

	"github.com/goccy/go-yaml"
)

func FuncMap(template *textTemplate.Template) textTemplate.FuncMap {
	return textTemplate.FuncMap{
		"toYaml":  functionToYaml,
		"include": functionInclude(template),
	}
}

// As seen in helm.
func functionToYaml(value any) string {
	marshalOptions := []yaml.EncodeOption{
		yaml.Indent(4),
		yaml.UseSingleQuote(true),
		yaml.UseLiteralStyleIfMultiline(true),
	}

	// Root sequences should not be indented
	// see: https://github.com/goccy/go-yaml/pull/855
	if reflect.ValueOf(value).Kind() != reflect.Slice {
		marshalOptions = append(marshalOptions, yaml.IndentSequence(true))
	}

	data, err := yaml.MarshalWithOptions(value, marshalOptions...)
	if err != nil {
		// Swallow errors inside a template.
		return ""
	}

	return strings.TrimSuffix(string(data), "\n")
}

// As seen in helm.
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
