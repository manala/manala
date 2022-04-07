package internal

import (
	"github.com/stretchr/testify/suite"
	"io"
	internalLog "manala/internal/log"
	"path/filepath"
	"testing"
)

type RepositorySuite struct{ suite.Suite }

func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}

var repositoryTestPath = filepath.Join("testdata", "repository")

func (s *RepositorySuite) Test() {
	log := internalLog.New(io.Discard)

	repository := &Repository{
		log:  log,
		path: "path",
		dir:  repositoryTestPath,
	}
	repository.recipeLoader = &RecipeRepositoryDirLoader{
		Log:        log,
		Repository: repository,
	}

	s.Equal("path", repository.Path())
	s.Equal("path", repository.Source())
	s.Equal(repositoryTestPath, repository.Dir())

	s.Run("LoadRecipe", func() {
		path := filepath.Join(repositoryTestPath, "load_recipe")
		repository.dir = path

		recipe, err := repository.LoadRecipe("recipe")

		s.NoError(err)
		s.Equal(filepath.Join(path, "recipe"), recipe.Path())
		s.Equal("recipe", recipe.Name())
		s.Equal("description", recipe.Description())
		s.Equal(map[string]interface{}{"foo": "bar"}, recipe.Vars())
	})

	s.Run("WalkRecipes Empty", func() {
		path := filepath.Join(repositoryTestPath, "walk_recipes_empty")
		repository.dir = path

		err := repository.WalkRecipes(func(recipe *Recipe) {})

		s.ErrorAs(err, &internalError)
		s.Equal("empty repository", internalError.Message)
	})

	s.Run("WalkRecipes", func() {
		path := filepath.Join(repositoryTestPath, "walk_recipes")
		repository.dir = path

		err := repository.WalkRecipes(func(recipe *Recipe) {
			s.Equal(filepath.Join(path, "recipe"), recipe.Path())
			s.Equal("recipe", recipe.Name())
			s.Equal("description", recipe.Description())
			s.Equal(map[string]interface{}{"foo": "bar"}, recipe.Vars())
		})

		s.NoError(err)
	})
}
