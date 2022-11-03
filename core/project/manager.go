package project

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"github.com/caarlos0/log"
	"github.com/xeipuuv/gojsonschema"
	"manala/core"
	internalLog "manala/internal/log"
	internalOs "manala/internal/os"
	internalReport "manala/internal/report"
	internalValidation "manala/internal/validation"
	internalWatcher "manala/internal/watcher"
	"os"
	"path/filepath"
)

const manifestFilename = ".manala.yaml"

//go:embed resources/manifest.yaml.tmpl
var manifestTemplate string

func NewManager(
	log *internalLog.Logger,
	repositoryManager core.RepositoryManager,
	recipeManager core.RecipeManager,
) *Manager {
	return &Manager{
		log:               log,
		repositoryManager: repositoryManager,
		recipeManager:     recipeManager,
	}
}

type Manager struct {
	log               *internalLog.Logger
	repositoryManager core.RepositoryManager
	recipeManager     core.RecipeManager
}

func (manager *Manager) IsProject(dir string) bool {
	manFile := filepath.Join(dir, manifestFilename)

	if _, err := os.Stat(manFile); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

func (manager *Manager) loadManifest(file string) (core.ProjectManifest, error) {
	// Log
	manager.log.WithFields(log.Fields{
		"file": file,
	}).Debug("load project manifest")

	// Stat file
	if fileInfo, err := os.Stat(file); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, internalReport.NewError(
				core.NewNotFoundProjectManifestError("project manifest not found"),
			).WithField("file", file)
		}
		return nil, internalReport.NewError(internalOs.NewError(err)).
			WithMessage("unable to stat project manifest").
			WithField("file", file)
	} else {
		if fileInfo.IsDir() {
			return nil, internalReport.NewError(fmt.Errorf("project manifest is a directory")).
				WithField("dir", file)
		}
	}

	man := NewManifest()

	// Open file
	if reader, err := os.Open(file); err != nil {
		return nil, internalReport.NewError(internalOs.NewError(err)).
			WithMessage("unable to open project manifest").
			WithField("file", file)
	} else {
		// Read from file
		if err = man.ReadFrom(reader); err != nil {
			return nil, internalReport.NewError(err).
				WithField("file", file)
		}
	}

	// Log
	manager.log.WithFields(log.Fields{
		"repository": man.Repository(),
		"recipe":     man.Recipe(),
	}).Debug("manifest")

	return man, nil
}

func (manager *Manager) LoadProject(dir string) (core.Project, error) {
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
	validation, err := gojsonschema.Validate(
		gojsonschema.NewGoLoader(proj.Recipe().Schema()),
		gojsonschema.NewGoLoader(proj.Vars()),
	)
	if err != nil {
		return nil, internalReport.NewError(err).
			WithMessage("unable to validate project manifest").
			WithField("file", manFile)
	}

	if !validation.Valid() {
		return nil, internalReport.NewError(
			internalValidation.NewError(
				"invalid project manifest vars",
				validation,
			).
				WithReporter(man),
		).WithField("file", manFile)
	}

	return proj, nil
}

func (manager *Manager) CreateProject(dir string, rec core.Recipe, vars map[string]interface{}) (core.Project, error) {
	template := rec.ProjectManifestTemplate().
		WithData(&core.ProjectView{
			Vars:   vars,
			Recipe: core.NewRecipeView(rec),
		}).
		WithDefaultContent(manifestTemplate)

	// Get final manifest content
	buffer := &bytes.Buffer{}
	if err := template.WriteTo(buffer); err != nil {
		return nil, err
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
				return nil, internalReport.NewError(internalOs.NewError(err)).
					WithMessage("unable to create project directory").
					WithField("dir", _dir)
			}
		} else {
			return nil, internalReport.NewError(internalOs.NewError(err)).
				WithMessage("unable to stat project directory").
				WithField("dir", _dir)
		}
	} else if !dirStat.IsDir() {
		return nil, internalReport.NewError(fmt.Errorf("project is not a directory")).
			WithField("file", _dir)
	}

	if writer, err := os.Create(manFile); err != nil {
		return nil, internalReport.NewError(internalOs.NewError(err)).
			WithMessage("unable to create project manifest file").
			WithField("file", manFile)
	} else {
		if _, err := writer.ReadFrom(bytes.NewReader(buffer.Bytes())); err != nil {
			return nil, internalReport.NewError(err).
				WithMessage("unable to save project manifest file").
				WithField("file", manFile)
		}

		if err := writer.Sync(); err != nil {
			return nil, internalReport.NewError(err).
				WithMessage("unable to sync project manifest file").
				WithField("file", manFile)
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

func (manager *Manager) WatchProject(proj core.Project, watcher *internalWatcher.Watcher) error {
	manFile := filepath.Join(proj.Dir(), manifestFilename)

	return watcher.AddGroup("project", manFile)
}
