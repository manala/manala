package template

import (
	"github.com/Masterminds/sprig/v3"
	"io"
	internalReport "manala/internal/report"
	"path/filepath"
	textTemplate "text/template"
)

type ProviderInterface interface {
	Template() *Template
}

type Provider struct{}

func (provider *Provider) Template() *Template {
	return &Template{}
}

type Template struct {
	defaultContent string
	defaultFiles   []string
	file           string
	data           any
}

func (template *Template) WriteTo(writer io.Writer) error {
	_template := textTemplate.New("")

	// Execution stops immediately with an error.
	_template.Option("missingkey=error")

	// Functions
	_template.Funcs(sprig.TxtFuncMap())
	_template.Funcs(FuncMap(_template))

	// Default files
	if len(template.defaultFiles) > 0 {
		for _, file := range template.defaultFiles {
			if _, err := _template.ParseFiles(file); err != nil {
				return internalReport.NewError(NewParsingError(err)).
					WithMessage("template error").
					WithField("path", file)
			}
		}
	}

	// File
	if template.file != "" {
		if _, err := _template.ParseFiles(template.file); err != nil {
			return internalReport.NewError(NewParsingError(err)).
				WithMessage("template error").
				WithField("path", template.file)
		}

		if err := _template.ExecuteTemplate(writer, filepath.Base(template.file), template.data); err != nil {
			return internalReport.NewError(NewExecutionError(err)).
				WithMessage("template error").
				WithField("path", template.file)
		}

		return nil
	}

	if _, err := _template.Parse(template.defaultContent); err != nil {
		return internalReport.NewError(NewParsingError(err)).
			WithMessage("template error")
	}

	if err := _template.Execute(writer, template.data); err != nil {
		return internalReport.NewError(NewExecutionError(err)).
			WithMessage("template error")
	}

	return nil
}

func (template *Template) WithData(data any) *Template {
	template.data = data

	return template
}

func (template *Template) WithDefaultFile(path string) *Template {
	template.defaultFiles = append(template.defaultFiles, path)

	return template
}

func (template *Template) WithFile(path string) *Template {
	template.file = path

	return template
}

func (template *Template) WithDefaultContent(content string) *Template {
	template.defaultContent = content

	return template
}
