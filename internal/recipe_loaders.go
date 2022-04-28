package internal

import (
	"errors"
	"github.com/apex/log"
	"io"
	internalLog "manala/internal/log"
	internalOs "manala/internal/os"
	"os"
	"path/filepath"
)

type RecipeLoaderInterface interface {
	LoadRecipe(name string) (*Recipe, error)
}

type RecipeWalkerInterface interface {
	WalkRecipes(walker func(recipe *Recipe)) error
}

/******************/
/* Repository Dir */
/******************/

type RecipeRepositoryDirLoader struct {
	Log        *internalLog.Logger
	Repository *Repository
}

func (loader *RecipeRepositoryDirLoader) LoadRecipeManifest(name string) (*RecipeManifest, error) {
	// Log
	loader.Log.WithFields(log.Fields{
		"name": name,
	}).Debug("try load")

	manifest := NewRecipeManifest(filepath.Join(loader.Repository.Dir(), name))

	// Stat file
	fileInfo, err := os.Stat(manifest.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, NotFoundRecipeManifestPathError(manifest.path)
		}
		return nil, internalOs.FileSystemError(err)
	}
	if fileInfo.IsDir() {
		return nil, WrongRecipeManifestPathError(manifest.path)
	}

	// Open file
	file, err := os.Open(manifest.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, NotFoundRecipeManifestPathError(manifest.path)
		}
		return nil, internalOs.FileSystemError(err)
	}

	// Write file content into manifest
	if _, err := io.Copy(manifest, file); err != nil {
		return nil, internalOs.FileSystemError(err)
	}

	if err := manifest.Load(); err != nil {
		return nil, err
	}

	// Log
	loader.Log.WithFields(log.Fields{
		"description": manifest.Description,
		"template":    manifest.Template,
	}).Debug("manifest")

	return manifest, nil
}

func (loader *RecipeRepositoryDirLoader) LoadRecipe(name string) (*Recipe, error) {
	// Load manifest
	manifest, err := loader.LoadRecipeManifest(name)
	if err != nil {
		return nil, err
	}

	recipe := &Recipe{
		name:       name,
		manifest:   manifest,
		repository: loader.Repository,
	}

	return recipe, nil
}
