package project

import (
	"bytes"
	_ "embed"
	"errors"
	"github.com/xeipuuv/gojsonschema"
	"log/slog"
	"manala/app/interfaces"
	"manala/app/views"
	"manala/core"
	"manala/internal/errors/serrors"
	"manala/internal/validation"
	"manala/internal/watcher"
	"os"
	"path/filepath"
)

const manifestFilename = ".manala.yaml"

//go:embed resources/manifest.yaml.tmpl
var manifestTemplate string

func NewManager(
	log *slog.Logger,
	repositoryManager interfaces.RepositoryManager,
	recipeManager interfaces.RecipeManager,
) *Manager {
	return &Manager{
		log:               log,
		repositoryManager: repositoryManager,
		recipeManager:     recipeManager,
	}
}

type Manager struct {
	log               *slog.Logger
	repositoryManager interfaces.RepositoryManager
	recipeManager     interfaces.RecipeManager
}

func (manager *Manager) IsProject(dir string) bool {
	manFile := filepath.Join(dir, manifestFilename)

	if _, err := os.Stat(manFile); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

func (manager *Manager) loadManifest(file string) (interfaces.ProjectManifest, error) {
	// Log
	manager.log.Debug("try to load project manifest",
		"file", file,
	)

	// Stat file
	if fileInfo, err := os.Stat(file); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, &core.NotFoundProjectManifestError{File: file}
		}
		return nil, serrors.WrapOs("unable to stat project manifest", err).
			WithArguments("file", file)
	} else {
		if fileInfo.IsDir() {
			return nil, serrors.New("project manifest is a directory").
				WithArguments("dir", file)
		}
	}

	man := NewManifest()

	// Open file
	if reader, err := os.Open(file); err != nil {
		return nil, serrors.WrapOs("unable to open project manifest", err).
			WithArguments("file", file)
	} else {
		// Read from file
		if err = man.ReadFrom(reader); err != nil {
			return nil, serrors.Wrap("unable to read project manifest", err).
				WithArguments("file", file)
		}
	}

	// Log
	manager.log.Debug("project manifest loaded",
		"repository", man.Repository(),
		"recipe", man.Recipe(),
	)

	return man, nil
}

func (manager *Manager) LoadProject(dir string) (interfaces.Project, error) {
	// Log
	manager.log.Debug("load project",
		"dir", dir,
	)

	// Load manifest
	manFile := filepath.Join(dir, manifestFilename)
	man, err := manager.loadManifest(manFile)
	if err != nil {
		return nil, err
	}

	// Load repository
	repo, err := manager.repositoryManager.LoadRepository(
		man.Repository(),
	)
	if err != nil {
		return nil, err
	}

	// Load recipe
	rec, err := manager.recipeManager.LoadRecipe(
		repo,
		man.Recipe(),
	)
	if err != nil {
		return nil, err
	}

	proj := NewProject(
		dir,
		man,
		rec,
	)

	// Validate vars against recipe
	val, err := gojsonschema.Validate(
		gojsonschema.NewGoLoader(proj.Recipe().Schema()),
		gojsonschema.NewGoLoader(proj.Vars()),
	)
	if err != nil {
		return nil, serrors.Wrap("unable to validate project manifest", err).
			WithArguments("file", manFile)
	}

	if !val.Valid() {
		return nil, validation.NewError(
			"invalid project manifest vars",
			val,
		).
			WithArguments("file", manFile).
			WithResultErrorDecorator(man.ValidationResultErrorDecorator())
	}

	return proj, nil
}

func (manager *Manager) CreateProject(dir string, rec interfaces.Recipe, vars map[string]interface{}) (interfaces.Project, error) {
	template := rec.ProjectManifestTemplate().
		WithData(&views.ProjectView{
			Vars:   vars,
			Recipe: views.NormalizeRecipe(rec),
		}).
		WithDefaultContent(manifestTemplate)

	// Get final manifest content
	buffer := &bytes.Buffer{}
	if err := template.WriteTo(buffer); err != nil {
		return nil, serrors.Wrap("recipe template error", err)
	}

	manFile := filepath.Join(dir, manifestFilename)

	man := NewManifest()
	if err := man.ReadFrom(bytes.NewReader(buffer.Bytes())); err != nil {
		return nil, err
	}

	// Ensure directory exists
	_dir := filepath.Dir(manFile)
	if dirStat, err := os.Stat(_dir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(_dir, 0755); err != nil {
				return nil, serrors.WrapOs("unable to create project directory", err).
					WithArguments("dir", _dir)
			}
		} else {
			return nil, serrors.WrapOs("unable to stat project directory", err).
				WithArguments("dir", _dir)
		}
	} else if !dirStat.IsDir() {
		return nil, serrors.New("project is not a directory").
			WithArguments("path", _dir)
	}

	if writer, err := os.Create(manFile); err != nil {
		return nil, serrors.WrapOs("unable to create project manifest file", err).
			WithArguments("file", manFile)
	} else {
		if _, err := writer.ReadFrom(bytes.NewReader(buffer.Bytes())); err != nil {
			return nil, serrors.Wrap("unable to save project manifest file", err).
				WithArguments("file", manFile)
		}

		if err := writer.Sync(); err != nil {
			return nil, serrors.Wrap("unable to sync project manifest file", err).
				WithArguments("file", manFile)
		}
	}

	// Final project
	proj := NewProject(
		dir,
		man,
		rec,
	)

	return proj, nil
}

func (manager *Manager) WatchProject(proj interfaces.Project, watcher *watcher.Watcher) error {
	manFile := filepath.Join(proj.Dir(), manifestFilename)

	return watcher.AddGroup("project", manFile)
}
