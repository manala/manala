package engine

import (
	"io"
	"os"
	"text/template"

	"github.com/manala/manala/internal/parsing"
	"github.com/manala/manala/internal/serrors"
)

type Executor struct {
	template *template.Template
	data     any
}

func (e *Executor) Execute(writer io.Writer, content string) error {
	if err := e.execute(writer, content); err != nil {
		return parsing.ErrorTo(
			serrors.New("unable to parse template"),
			ErrorFrom(err, content),
			parsing.Options{
				Src:   content,
				Lexer: "go-template",
			},
		)
	}
	return nil
}

func (e *Executor) ExecuteTemplate(writer io.Writer, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return serrors.New("unable to read template file").
			WithArguments("path", path).
			WithErrors(serrors.NewOs(err))
	}

	if err := e.execute(writer, string(content)); err != nil {
		return parsing.ErrorTo(
			serrors.New("unable to parse template file").
				WithArguments("path", path),
			ErrorFrom(err, string(content)),
			parsing.Options{
				Src:   string(content),
				Lexer: "go-template",
			},
		)
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
