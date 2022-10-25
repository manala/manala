package recipe

import (
	"errors"
	"fmt"
	"github.com/caarlos0/log"
	"manala/core"
	internalLog "manala/internal/log"
	internalOs "manala/internal/os"
	internalReport "manala/internal/report"
	"os"
	"path/filepath"
)

/**************/
/* Repository */
/**************/

func NewRepositoryManager(log *internalLog.Logger, repo core.Repository) *RepositoryManager {
	return &RepositoryManager{
		log:        log,
		repository: repo,
	}
}

type RepositoryManager struct {
	log        *internalLog.Logger
	repository core.Repository
}

func (manager *RepositoryManager) LoadRecipeManifest(name string) (*Manifest, error) {
	// Log
	manager.log.WithFields(log.Fields{
		"name": name,
	}).Debug("try load")

	path := filepath.Join(manager.repository.Dir(), name)

	manifest := NewManifest(path)

	// Stat file
	fileInfo, err := os.Stat(manifest.Path())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, internalReport.NewError(
				core.NewNotFoundRecipeManifestError("recipe manifest not found"),
			).WithField("path", path)
		}
		return nil, internalReport.NewError(internalOs.NewError(err)).
			WithMessage("unable to stat recipe manifest").
			WithField("path", manifest.Path())
	}
	if fileInfo.IsDir() {
		return nil, internalReport.NewError(fmt.Errorf("recipe manifest is a directory")).
			WithField("path", manifest.Path())
	}

	// Open file
	file, err := os.Open(manifest.Path())
	if err != nil {
		return nil, internalReport.NewError(internalOs.NewError(err)).
			WithMessage("unable to open recipe manifest").
			WithField("path", manifest.Path())
	}

	// Read from file
	if err = manifest.ReadFrom(file); err != nil {
		return nil, internalReport.NewError(err).
			WithField("path", manifest.Path())
	}

	// Log
	manager.log.WithFields(log.Fields{
		"description": manifest.Description(),
		"template":    manifest.Template(),
	}).Debug("manifest")

	return manifest, nil
}

func (manager *RepositoryManager) LoadRecipe(name string) (core.Recipe, error) {
	// Load manifest
	manifest, err := manager.LoadRecipeManifest(name)
	if err != nil {
		return nil, err
	}

	rec := NewRecipe(
		name,
		manifest,
		manager.repository,
	)

	return rec, nil
}
