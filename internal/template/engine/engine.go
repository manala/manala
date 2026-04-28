package engine

import (
	"os"
	"text/template"

	"github.com/manala/manala/internal/parsing"
	"github.com/manala/manala/internal/serrors"

	"github.com/Masterminds/sprig/v3"
)

type Engine struct {
	template *template.Template
}

func New() *Engine {
	t := template.New("")

	// Execution stops immediately with an error.
	t.Option("missingkey=error")

	// Sprig functions
	t.Funcs(sprig.TxtFuncMap())

	return &Engine{
		template: t,
	}
}

func (engine *Engine) Executor(data any, files ...string) (*Executor, error) {
	// Clone base template to isolate each executor.
	clone, _ := engine.template.Clone()

	// Parse files.
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return nil, serrors.New("unable to read template file").
				With("path", file).
				WithErrors(serrors.FromOs(err))
		}

		if _, err := clone.Parse(string(content)); err != nil {
			return nil, serrors.New("unable to parse template file").
				WithDumper(parsing.ErrorDumper{
					Err:   ErrorFrom(err, string(content)),
					File:  file,
					Src:   string(content),
					Lexer: "go-template",
				})
		}
	}

	return &Executor{
		template: clone,
		data:     data,
	}, nil
}
