package repository

import (
	"errors"
	"fmt"
	"github.com/caarlos0/log"
	"manala/core"
	"manala/core/recipe"
	internalFilepath "manala/internal/filepath"
	internalLog "manala/internal/log"
	internalOs "manala/internal/os"
	internalReport "manala/internal/report"
	"os"
	"sort"
)

func NewRepository(log *internalLog.Logger, path string, dir string) *Repository {
	repo := &Repository{
		log:  log,
		path: path,
		dir:  dir,
	}

	repo.recipeManager = recipe.NewRepositoryManager(
		log,
		repo,
	)

	return repo
}

type Repository struct {
	log           *internalLog.Logger
	path          string
	dir           string
	recipeManager core.RecipeManager
}

func (repo *Repository) Path() string {
	return repo.path
}

// Source keep backward compatibility, when "source" was used instead of "path"
// to define repository origin
func (repo *Repository) Source() string {
	return repo.Path()
}

func (repo *Repository) Dir() string {
	return repo.dir
}

func (repo *Repository) LoadRecipe(name string) (core.Recipe, error) {
	return repo.recipeManager.LoadRecipe(name)
}

func (repo *Repository) WalkRecipes(walker func(rec core.Recipe)) error {
	dir, err := os.Open(repo.dir)
	if err != nil {
		return internalReport.NewError(internalOs.NewError(err)).
			WithMessage("file system error")
	}

	//goland:noinspection GoUnhandledErrorResult
	defer dir.Close()

	files, err := dir.ReadDir(0) // 0 to read all files and folders
	if err != nil {
		return internalReport.NewError(internalOs.NewError(err)).
			WithMessage("file system error")
	}

	// Sort alphabetically
	sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })

	empty := true

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		// Exclusions
		if internalFilepath.Exclude(file.Name()) {
			repo.log.WithFields(log.Fields{
				"name": file.Name(),
			}).Debug("exclude recipe")
			continue
		}

		// Log
		repo.log.WithFields(log.Fields{
			"name": file.Name(),
		}).Debug("load recipe")
		repo.log.IncreasePadding()

		rec, err := repo.LoadRecipe(file.Name())

		// Log
		repo.log.DecreasePadding()

		if err != nil {
			var _notFoundRecipeManifestError *core.NotFoundRecipeManifestError
			if errors.As(err, &_notFoundRecipeManifestError) {
				continue
			}
			return err
		}

		empty = false

		walker(rec)
	}

	if empty {
		return internalReport.NewError(fmt.Errorf("empty repository")).
			WithField("dir", repo.dir)
	}

	return nil
}
