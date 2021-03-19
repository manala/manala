package template

import (
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"gopkg.in/yaml.v3"
	"io"
	"manala/fs"
	"strings"
	textTemplate "text/template"
)

/***********/
/* Manager */
/***********/

// Create a template manager
func NewManager() *manager {
	return &manager{}
}

type ManagerInterface interface {
	NewFsTemplate(fs fs.ReadInterface) *Template
}

type manager struct {
}

// Create a fs template
func (manager *manager) NewFsTemplate(fs fs.ReadInterface) *Template {
	// Template
	tmpl := &Template{
		txtTmpl: textTemplate.New(""),
		fs:      fs,
	}

	// Execution stops immediately with an error.
	tmpl.txtTmpl.Option("missingkey=error")

	// Functions
	tmpl.txtTmpl.Funcs(sprig.TxtFuncMap())
	tmpl.txtTmpl.Funcs(textTemplate.FuncMap{
		"toYaml":  tmpl.funcToYaml(),
		"include": tmpl.funcInclude(),
	})

	return tmpl
}

/************/
/* Template */
/************/

type Interface interface {
	Parse(text string) error
	ParseFile(filename string) error
	ParseFiles(filenames ...string) error
	Execute(writer io.Writer, vars map[string]interface{}) error
}

type Template struct {
	txtTmpl *textTemplate.Template
	fs      fs.ReadInterface
}

func (tmpl *Template) Parse(text string) error {
	_, err := tmpl.txtTmpl.Parse(
		text,
	)

	return err
}

func (tmpl *Template) ParseFile(filename string) error {
	text, err := tmpl.fs.ReadFile(filename)
	if err != nil {
		return err
	}

	return tmpl.Parse(string(text))
}

func (tmpl *Template) ParseFiles(filenames ...string) error {
	_, err := tmpl.txtTmpl.ParseFS(tmpl.fs, filenames...)

	return err
}

func (tmpl *Template) Execute(writer io.Writer, vars map[string]interface{}) error {
	return tmpl.txtTmpl.Execute(
		writer,
		vars,
	)
}

// As seen in helm
func (tmpl *Template) funcToYaml() func(value interface{}) string {
	return func(value interface{}) string {
		data, err := yaml.Marshal(value)
		if err != nil {
			// Swallow errors inside of a template.
			return ""
		}
		return strings.TrimSuffix(string(data), "\n")
	}
}

// As seen in helm
func (tmpl *Template) funcInclude() func(name string, data interface{}) (string, error) {
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
		err := tmpl.txtTmpl.ExecuteTemplate(&buf, name, data)
		return buf.String(), err
	}
}
