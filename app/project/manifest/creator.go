package manifest

import (
	"bytes"
	_ "embed"
	"errors"
	"os"
	"path/filepath"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/template"
	"github.com/manala/manala/internal/serrors"
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
			return nil, serrors.New("recipe template error").
				WithErrors(err)
		}
	} else {
		if err := templateExecutor.Execute(buffer, _template); err != nil {
			return nil, serrors.New("recipe template error").
				WithErrors(err)
		}
	}

	manifestFile := filepath.Join(dir, filename)

	manifest := New()
	if err := manifest.Unmarshal(buffer.Bytes()); err != nil {
		return nil, err
	}

	// Ensure directory exists
	_dir := filepath.Dir(manifestFile)
	if dirStat, err := os.Stat(_dir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(_dir, 0o755); err != nil {
				return nil, serrors.New("unable to create project directory").
					WithArguments("dir", _dir).
					WithErrors(serrors.FromOs(err))
			}
		} else {
			return nil, serrors.New("unable to stat project directory").
				WithArguments("dir", _dir).
				WithErrors(serrors.FromOs(err))
		}
	} else if !dirStat.IsDir() {
		return nil, serrors.New("project is not a directory").
			WithArguments("path", _dir)
	}

	writer, err := os.Create(manifestFile)
	if err != nil {
		return nil, serrors.New("unable to create project manifest file").
			WithArguments("file", manifestFile).
			WithErrors(serrors.FromOs(err))
	}

	if _, err := writer.ReadFrom(bytes.NewReader(buffer.Bytes())); err != nil {
		return nil, serrors.New("unable to save project manifest file").
			WithArguments("file", manifestFile).
			WithErrors(err)
	}

	if err := writer.Sync(); err != nil {
		return nil, serrors.New("unable to sync project manifest file").
			WithArguments("file", manifestFile).
			WithErrors(err)
	}

	// Final project
	return NewProject(dir, manifest, recipe), nil
}
