package manifest

import (
	"bytes"
	"errors"
	"manala/app"
	"manala/internal/serrors"
	"os"
	"path/filepath"
)

func NewCreator() *Creator {
	return &Creator{}
}

type Creator struct{}

func (creator *Creator) Create(dir string, recipe app.Recipe, vars map[string]any) (app.Project, error) {
	template := recipe.ProjectManifestTemplate().
		WithData(&app.ProjectView{
			Vars:   vars,
			Recipe: app.NewRecipeView(recipe),
		}).
		WithDefaultContent(_template)

	// Get final manifest content
	buffer := &bytes.Buffer{}
	if err := template.WriteTo(buffer); err != nil {
		return nil, serrors.New("recipe template error").
			WithErrors(err)
	}

	manifestFile := filepath.Join(dir, filename)

	manifest := New()
	if _, err := manifest.ReadFrom(bytes.NewReader(buffer.Bytes())); err != nil {
		return nil, err
	}

	// Ensure directory exists
	_dir := filepath.Dir(manifestFile)
	if dirStat, err := os.Stat(_dir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(_dir, 0755); err != nil {
				return nil, serrors.New("unable to create project directory").
					WithArguments("dir", _dir).
					WithErrors(serrors.NewOs(err))
			}
		} else {
			return nil, serrors.New("unable to stat project directory").
				WithArguments("dir", _dir).
				WithErrors(serrors.NewOs(err))
		}
	} else if !dirStat.IsDir() {
		return nil, serrors.New("project is not a directory").
			WithArguments("path", _dir)
	}

	writer, err := os.Create(manifestFile)
	if err != nil {
		return nil, serrors.New("unable to create project manifest file").
			WithArguments("file", manifestFile).
			WithErrors(serrors.NewOs(err))
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
