package manifest

import (
	"bytes"
	_ "embed"
	"errors"
	"os"
	"path/filepath"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/template"
	"github.com/manala/manala/internal/errors/serror"
	"github.com/manala/manala/internal/errors/std"
)

//go:embed template.yaml.tmpl
var _template string

type Creator struct {
	templateEngine *template.Engine
}

func NewCreator(templateEngine *template.Engine) *Creator {
	return &Creator{
		templateEngine: templateEngine,
	}
}

func (creator *Creator) Create(dir string, recipe app.Recipe, vars map[string]any) (app.Project, error) {
	templateExecutor, err := creator.templateEngine.Executor(
		vars,
		recipe,
		dir,
	)
	if err != nil {
		return nil, err
	}

	// Get final manifest content
	buffer := &bytes.Buffer{}
	if template := recipe.Template(); template != "" {
		if err := templateExecutor.ExecuteTemplate(buffer, template); err != nil {
			return nil, serror.New("recipe template error").
				WithErr(err)
		}
	} else {
		if err := templateExecutor.Execute(buffer, _template); err != nil {
			return nil, serror.New("recipe template error").
				WithErr(err)
		}
	}

	manifestFile := filepath.Join(dir, filename)

	// Ensure directory exists
	_dir := filepath.Dir(manifestFile)
	if dirStat, err := os.Stat(_dir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(_dir, 0o755); err != nil {
				return nil, serror.New("unable to create project directory").
					With("dir", _dir).
					WithErr(std.From(err))
			}
		} else {
			return nil, serror.New("unable to stat project directory").
				With("dir", _dir).
				WithErr(std.From(err))
		}
	} else if !dirStat.IsDir() {
		return nil, serror.New("project is not a directory").
			With("path", _dir)
	}

	writer, err := os.Create(manifestFile)
	if err != nil {
		return nil, serror.New("unable to create project manifest file").
			With("file", manifestFile).
			WithErr(std.From(err))
	}

	if _, err := writer.ReadFrom(bytes.NewReader(buffer.Bytes())); err != nil {
		return nil, serror.New("unable to save project manifest file").
			With("file", manifestFile).
			WithErr(err)
	}

	if err := writer.Sync(); err != nil {
		return nil, serror.New("unable to sync project manifest file").
			With("file", manifestFile).
			WithErr(err)
	}

	// Final project
	return &Project{
		dir:    dir,
		recipe: recipe,
		vars:   vars,
	}, nil
}
