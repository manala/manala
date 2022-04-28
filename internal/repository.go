package internal

import (
	"errors"
	"github.com/apex/log"
	internalFilepath "manala/internal/filepath"
	internalLog "manala/internal/log"
	internalOs "manala/internal/os"
	"os"
)

type Repository struct {
	log          *internalLog.Logger
	path         string
	dir          string
	recipeLoader RecipeLoaderInterface
}

func (repository *Repository) Path() string {
	return repository.path
}

// Source keep backward compatibility, when "source" was used instead of "path"
// to define repository origin
func (repository *Repository) Source() string {
	return repository.Path()
}

func (repository *Repository) Dir() string {
	return repository.dir
}

func (repository *Repository) LoadRecipe(name string) (*Recipe, error) {
	return repository.recipeLoader.LoadRecipe(name)
}

func (repository *Repository) WalkRecipes(walker func(recipe *Recipe)) error {
	dir, err := os.Open(repository.dir)
	if err != nil {
		return internalOs.FileSystemError(err)
	}
	defer dir.Close()

	files, err := dir.Readdir(0) // 0 to read all files and folders
	if err != nil {
		return internalOs.FileSystemError(err)
	}

	empty := true

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		// Exclusions
		if internalFilepath.Exclude(file.Name()) {
			repository.log.WithFields(log.Fields{
				"name": file.Name(),
			}).Debug("exclude recipe")
			continue
		}

		// Log
		repository.log.WithFields(log.Fields{
			"name": file.Name(),
		}).Debug("load recipe")
		repository.log.PaddingUp()

		recipe, err := repository.LoadRecipe(file.Name())

		// Log
		repository.log.PaddingDown()

		if err != nil {
			var _err *NotFoundRecipeManifestError
			if errors.As(err, &_err) {
				continue
			}
			return err
		}

		empty = false

		walker(recipe)
	}

	if empty {
		return EmptyRepositoryDirError(repository.dir)
	}

	return nil
}
