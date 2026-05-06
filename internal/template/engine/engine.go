package engine

import (
	"os"
	"text/template"

	"github.com/manala/manala/internal/errors/serror"
	"github.com/manala/manala/internal/errors/source"
	"github.com/manala/manala/internal/errors/std"
	templateerrors "github.com/manala/manala/internal/template/errors"

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
			return nil, serror.New("unable to read template file").
				With("path", file).
				WithErr(std.From(err))
		}

		if _, err := clone.Parse(string(content)); err != nil {
			return nil, serror.New("unable to parse template file").
				WithErr(source.From(templateerrors.From(err, string(content)), source.Origin{
					File:     file,
					Source:   string(content),
					Language: "go-template",
				}))
		}
	}

	return &Executor{
		template: clone,
		data:     data,
	}, nil
}
