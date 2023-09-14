package project

import (
	"bytes"
	_ "embed"
	"errors"
	"log/slog"
	"manala/app"
	"manala/internal/serrors"
	"manala/internal/validator"
	"manala/internal/watcher"
	"os"
	"path/filepath"
)

const manifestFilename = ".manala.yaml"

//go:embed resources/manifest.yaml.tmpl
var manifestTemplate string

func NewManager(
	log *slog.Logger,
	repositoryManager app.RepositoryManager,
	recipeManager app.RecipeManager,
) *Manager {
	return &Manager{
		log:               log,
		repositoryManager: repositoryManager,
		recipeManager:     recipeManager,
	}
}

type Manager struct {
	log               *slog.Logger
	repositoryManager app.RepositoryManager
	recipeManager     app.RecipeManager
}

func (manager *Manager) IsProject(dir string) bool {
	manifestFile := filepath.Join(dir, manifestFilename)

	if _, err := os.Stat(manifestFile); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

func (manager *Manager) loadManifest(file string) (*Manifest, error) {
	// Log
	manager.log.Debug("try to load project manifest",
		"file", file,
	)

	// Stat file
	if fileInfo, err := os.Stat(file); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, &app.NotFoundProjectManifestError{File: file}
		}
		return nil, serrors.New("unable to stat project manifest").
			WithArguments("file", file).
			WithErrors(serrors.NewOs(err))
	} else {
		if fileInfo.IsDir() {
			return nil, serrors.New("project manifest is a directory").
				WithArguments("dir", file)
		}
	}

	manifest := NewManifest()

	// Open file
	if reader, err := os.Open(file); err != nil {
		return nil, serrors.New("unable to open project manifest").
			WithArguments("file", file).
			WithErrors(serrors.NewOs(err))
	} else {
		// Read from file
		if _, err = manifest.ReadFrom(reader); err != nil {
			return nil, serrors.New("unable to read project manifest").
				WithArguments("file", file).
				WithErrors(err)
		}
	}

	// Log
	manager.log.Debug("project manifest loaded",
		"repository", manifest.Repository(),
		"recipe", manifest.Recipe(),
	)

	return manifest, nil
}

func (manager *Manager) LoadProject(dir string) (app.Project, error) {
	// Log
	manager.log.Debug("load project",
		"dir", dir,
	)

	// Load manifest
	manifestFile := filepath.Join(dir, manifestFilename)
	manifest, err := manager.loadManifest(manifestFile)
	if err != nil {
		return nil, err
	}

	// Load repository
	repository, err := manager.repositoryManager.LoadRepository(
		manifest.Repository(),
	)
	if err != nil {
		return nil, err
	}

	// Load recipe
	recipe, err := manager.recipeManager.LoadRecipe(
		repository,
		manifest.Recipe(),
	)
	if err != nil {
		return nil, err
	}

	project := New(
		dir,
		manifest,
		recipe,
	)

	// Validate project vars against recipe
	if violations, err := validator.New(
		validator.WithValidators(
			project.Recipe().ProjectValidator(),
		),
		validator.WithFormatters(
			manifest.ValidatorFormatter(),
		),
	).Validate(project.Vars()); err != nil {
		return nil, serrors.New("unable to validate project manifest").
			WithArguments("file", manifestFile).
			WithErrors(err)
	} else if len(violations) != 0 {
		return nil, serrors.New("invalid project manifest vars").
			WithArguments("file", manifestFile).
			WithErrors(violations.StructuredErrors()...)
	}

	return project, nil
}

func (manager *Manager) CreateProject(dir string, recipe app.Recipe, vars map[string]any) (app.Project, error) {
	template := recipe.ProjectManifestTemplate().
		WithData(&app.ProjectView{
			Vars:   vars,
			Recipe: app.NewRecipeView(recipe),
		}).
		WithDefaultContent(manifestTemplate)

	// Get final manifest content
	buffer := &bytes.Buffer{}
	if err := template.WriteTo(buffer); err != nil {
		return nil, serrors.New("recipe template error").
			WithErrors(err)
	}

	manifestFile := filepath.Join(dir, manifestFilename)

	manifest := NewManifest()
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

	if writer, err := os.Create(manifestFile); err != nil {
		return nil, serrors.New("unable to create project manifest file").
			WithArguments("file", manifestFile).
			WithErrors(serrors.NewOs(err))
	} else {
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
	}

	// Final project
	return New(
		dir,
		manifest,
		recipe,
	), nil
}

func (manager *Manager) WatchProject(project app.Project, watcher *watcher.Watcher) error {
	manifestFile := filepath.Join(project.Dir(), manifestFilename)

	return watcher.AddGroup("project", manifestFile)
}
