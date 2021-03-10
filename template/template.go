package template

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"gopkg.in/yaml.v3"
	"io"
	"strings"
	textTemplate "text/template"
)

type Template struct {
	tmpl *textTemplate.Template
}

func New() *Template {
	// Text template
	tmpl := &Template{
		tmpl: textTemplate.New(""),
	}

	// Execution stops immediately with an error.
	tmpl.tmpl.Option("missingkey=error")

	tmpl.tmpl.Funcs(sprig.TxtFuncMap())
	tmpl.tmpl.Funcs(textTemplate.FuncMap{
		"toYaml":  tmpl.funcToYaml(),
		"include": tmpl.funcInclude(),
	})

	return tmpl
}

func (tmpl *Template) ParseFiles(filenames ...string) error {
	_, err := tmpl.tmpl.ParseFiles(filenames...)

	return err
}

func (tmpl *Template) Parse(text string) error {
	_, err := tmpl.tmpl.Parse(text)

	return err
}

func (tmpl *Template) Execute(wr io.Writer, data interface{}) error {
	return tmpl.tmpl.Execute(wr, data)
}

// As seen in helm
func (tmpl *Template) funcToYaml() func(value interface{}) string {
	return func(value interface{}) string {
		var buf bytes.Buffer

		enc := yaml.NewEncoder(&buf)

		if err := enc.Encode(value); err != nil {
			// Swallow errors inside of a template.
			return ""
		}

		return strings.TrimSuffix(buf.String(), "\n")
	}
}

// As seen in helm
func (tmpl *Template) funcInclude() func(name string, data interface{}) (string, error) {
	includedNames := make([]string, 0)
	return func(name string, data interface{}) (string, error) {
		var buf strings.Builder
		includedCount := 0
		for _, n := range includedNames {
			if n == name {
				includedCount += 1
			}
		}
		if includedCount >= 16 {
			return "", fmt.Errorf("rendering template has reached the maximum nested reference name level: %s", name)
		}
		includedNames = append(includedNames, name)
		err := tmpl.tmpl.ExecuteTemplate(&buf, name, data)
		includedNames = includedNames[:len(includedNames)-1]
		return buf.String(), err
	}
}
