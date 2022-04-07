package template

import (
	"fmt"
	"github.com/goccy/go-yaml"
	"strings"
	textTemplate "text/template"
)

// As seen in helm
func funcToYaml(template *textTemplate.Template) func(value interface{}) string {
	return func(value interface{}) string {
		data, err := yaml.MarshalWithOptions(value, yaml.Indent(4), yaml.IndentSequence(true))
		if err != nil {
			// Swallow errors inside a template.
			return ""
		}
		return strings.TrimSuffix(string(data), "\n")
	}
}

// As seen in helm
func funcInclude(template *textTemplate.Template) func(name string, data interface{}) (string, error) {
	includedNames := make(map[string]int)
	return func(name string, data interface{}) (string, error) {
		var buf strings.Builder
		if v, ok := includedNames[name]; ok {
			if v > 1000 {
				return "", fmt.Errorf("rendering template has a nested reference name: %s", name)
			}
			includedNames[name]++
		} else {
			includedNames[name] = 1
		}
		err := template.ExecuteTemplate(&buf, name, data)
		return buf.String(), err
	}
}
