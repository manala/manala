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
	"os"
	"path/filepath"
)

//go:embed resources/manifest.yaml.tmpl
var manifestTemplate string

func NewManager(
	log *internalLog.Logger,
	repositoryManager core.RepositoryManager,
) *Manager {
	return &Manager{
		log:               log,
		repositoryManager: repositoryManager,
	}
}

type Manager struct {
	log               *internalLog.Logger
	repositoryManager core.RepositoryManager
}

func (manager *Manager) LoadProjectManifest(path string) (core.ProjectManifest, error) {
	// Log
	manager.log.WithFields(log.Fields{
		"path": path,
	}).Debug("load project manifest")

	manifest := NewManifest(path)

	// Stat file
	fileInfo, err := os.Stat(manifest.Path())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, internalReport.NewError(
				core.NewNotFoundProjectManifestError("project manifest not found"),
			).WithField("path", path)
		}
		return nil, internalReport.NewError(internalOs.NewError(err)).
			WithMessage("unable to stat project manifest").
			WithField("path", manifest.Path())
	}
	if fileInfo.IsDir() {
		return nil, internalReport.NewError(fmt.Errorf("project manifest is a directory")).
			WithField("path", manifest.Path())
	}

	// Open file
	file, err := os.Open(manifest.Path())
	if err != nil {
		return nil, internalReport.NewError(internalOs.NewError(err)).
			WithMessage("unable to open project manifest").
			WithField("path", manifest.Path())
	}

	// Read from file
	if err = manifest.ReadFrom(file); err != nil {
		return nil, internalReport.NewError(err).
			WithField("path", manifest.Path())
	}

	// Log
	manager.log.WithFields(log.Fields{
		"repository": manifest.Repository(),
		"recipe":     manifest.Recipe(),
	}).Debug("manifest")

	return manifest, nil
}

func (manager *Manager) LoadProject(path string, repoPath string, recName string) (core.Project, error) {
	// Load manifest
	manifest, err := manager.LoadProjectManifest(path)
	if err != nil {
		return nil, err
	}

	// Log
	manager.log.WithFields(log.Fields{
		"repository": manifest.Repository(),
		"recipe":     manifest.Recipe(),
	}).Debug("manifest")

	// Log
	manager.log.WithFields(log.Fields{
		"path":          repoPath,
		"manifest.path": manifest.Repository(),
	}).Debug("load repository")
	manager.log.IncreasePadding()

	// Load repository
	repo, err := manager.repositoryManager.LoadRepository([]string{
		repoPath,
		manifest.Repository(),
	})
	if err != nil {
		return nil, err
	}

	// Log
	manager.log.DecreasePadding()
	manager.log.WithFields(log.Fields{
		"name":          recName,
		"manifest.name": manifest.Recipe(),
	}).Debug("load recipe")
	manager.log.IncreasePadding()

	// Load recipe
	if recName == "" {
		recName = manifest.Recipe()
	}
	rec, err := repo.LoadRecipe(recName)
	if err != nil {
		return nil, err
	}

	proj := NewProject(
		manifest,
		rec,
	)

	// Log
	manager.log.DecreasePadding()

	// Validate vars against recipe
	validation, err := gojsonschema.Validate(
		gojsonschema.NewGoLoader(proj.Recipe().Schema()),
		gojsonschema.NewGoLoader(proj.Vars()),
	)
	if err != nil {
		return nil, internalReport.NewError(err).
			WithMessage("unable to validate project manifest").
			WithField("path", manifest.Path())
	}

	if !validation.Valid() {
		return nil, internalReport.NewError(
			internalValidation.NewError(
				"invalid project manifest",
				validation,
				internalValidation.WithReporter(manifest),
			),
		).WithField("path", manifest.Path())
	}

	return proj, nil
}

func (manager *Manager) CreateProject(path string, rec core.Recipe, vars map[string]interface{}) (core.Project, error) {
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

	manifest := NewManifest(path)
	if err := manifest.ReadFrom(bytes.NewReader(buffer.Bytes())); err != nil {
		return nil, err
	}

	// Ensure directory exists
	dir := filepath.Dir(manifest.Path())
	if dirStat, err := os.Stat(dir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, internalReport.NewError(internalOs.NewError(err)).
					WithMessage("unable to create project directory").
					WithField("path", dir)
			}
		} else {
			return nil, internalReport.NewError(internalOs.NewError(err)).
				WithMessage("unable to stat project directory").
				WithField("path", dir)
		}
	} else if !dirStat.IsDir() {
		return nil, internalReport.NewError(fmt.Errorf("project is not a directory")).
			WithField("path", dir)
	}

	manifestFile, err := os.Create(manifest.Path())
	if err != nil {
		return nil, internalReport.NewError(internalOs.NewError(err)).
			WithMessage("unable to create project manifest file").
			WithField("path", manifest.Path())
	}

	if _, err := manifestFile.ReadFrom(bytes.NewReader(buffer.Bytes())); err != nil {
		return nil, internalReport.NewError(err).
			WithMessage("unable to save project manifest file").
			WithField("path", manifest.Path())
	}

	if err := manifestFile.Sync(); err != nil {
		return nil, internalReport.NewError(err).
			WithMessage("unable to sync project manifest file").
			WithField("path", manifest.Path())
	}

	// Final project
	proj := NewProject(
		manifest,
		rec,
	)

	return proj, nil
}
