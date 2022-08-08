package internal

import (
	"errors"
	"github.com/caarlos0/log"
	"io"
	internalLog "manala/internal/log"
	internalOs "manala/internal/os"
	internalValidator "manala/internal/validator"
	"os"
)

type ProjectLoaderInterface interface {
	LoadProjectManifest(path string) (*ProjectManifest, error)
	LoadProject(path string, repositoryPath string, recipeName string) (*Project, error)
}

/*******/
/* Dir */
/*******/

type ProjectDirLoader struct {
	Log              *internalLog.Logger
	RepositoryLoader RepositoryLoaderInterface
}

func (loader *ProjectDirLoader) LoadProjectManifest(path string) (*ProjectManifest, error) {
	// Log
	loader.Log.WithFields(log.Fields{
		"path": path,
	}).Debug("load project manifest")

	manifest := NewProjectManifest(path)

	// Stat file
	fileInfo, err := os.Stat(manifest.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, NotFoundProjectManifestPathError(manifest.path)
		}
		return nil, internalOs.FileSystemError(err)
	}
	if fileInfo.IsDir() {
		return nil, WrongProjectManifestPathError(manifest.path)
	}

	// Open file
	file, err := os.Open(manifest.path)
	if err != nil {
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
		"repository": manifest.Repository,
		"recipe":     manifest.Recipe,
	}).Debug("manifest")

	return manifest, nil
}

func (loader *ProjectDirLoader) LoadProject(path string, repositoryPath string, recipeName string) (*Project, error) {
	// Load manifest
	manifest, err := loader.LoadProjectManifest(path)
	if err != nil {
		return nil, err
	}

	project := &Project{
		manifest: manifest,
	}

	// Log
	loader.Log.WithFields(log.Fields{
		"repository": project.manifest.Repository,
		"recipe":     project.manifest.Recipe,
	}).Debug("manifest")

	// Log
	loader.Log.WithFields(log.Fields{
		"path":          repositoryPath,
		"manifest.path": project.manifest.Repository,
	}).Debug("load repository")
	loader.Log.IncreasePadding()

	// Load repository
	repository, err := loader.RepositoryLoader.LoadRepository([]string{
		repositoryPath,
		project.manifest.Repository,
	})
	if err != nil {
		return nil, err
	}

	// Log
	loader.Log.DecreasePadding()
	loader.Log.WithFields(log.Fields{
		"name":          recipeName,
		"manifest.name": project.manifest.Recipe,
	}).Debug("load recipe")
	loader.Log.IncreasePadding()

	// Load recipe
	if recipeName == "" {
		recipeName = project.manifest.Recipe
	}
	project.recipe, err = repository.LoadRecipe(recipeName)
	if err != nil {
		return nil, err
	}

	// Log
	loader.Log.DecreasePadding()

	// Validate vars against recipe
	if err, errs, ok := internalValidator.Validate(project.recipe.Schema(), project.Vars(), internalValidator.WithYamlContent(manifest.content)); !ok {
		return nil, ValidationProjectManifestPathError(manifest.path, err, errs)
	}

	return project, nil
}
