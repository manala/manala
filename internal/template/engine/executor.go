package engine

import (
	"io"
	"os"
	"text/template"

	"github.com/manala/manala/internal/errors/serror"
	"github.com/manala/manala/internal/errors/source"
	"github.com/manala/manala/internal/errors/std"
	templateerrors "github.com/manala/manala/internal/template/errors"
)

type Executor struct {
	template *template.Template
	data     any
}

func (e *Executor) Execute(writer io.Writer, content string) error {
	if err := e.execute(writer, content); err != nil {
		return serror.New("unable to parse template").
			WithErr(source.From(templateerrors.From(err, content), source.Origin{
				Source:   content,
				Language: "go-template",
			}))
	}
	return nil
}

func (e *Executor) ExecuteTemplate(writer io.Writer, file string) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return serror.New("unable to read template file").
			With("file", file).
			WithErr(std.From(err))
	}

	if err := e.execute(writer, string(content)); err != nil {
		return serror.New("unable to parse template file").
			WithErr(source.From(templateerrors.From(err, string(content)), source.Origin{
				File:     file,
				Source:   string(content),
				Language: "go-template",
			}))
	}
	return nil
}

func (e *Executor) execute(writer io.Writer, content string) error {
	clone, _ := e.template.Clone()

	// Custom funcs are registered on the clone, because
	// they need a reference to the executing template.
	clone.Funcs(Funcs(clone))

	if _, err := clone.Parse(content); err != nil {
		return err
	}

	return clone.Execute(writer, e.data)
}
