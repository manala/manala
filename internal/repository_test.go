package internal

import (
	"github.com/stretchr/testify/suite"
	"io"
	internalLog "manala/internal/log"
	internalTesting "manala/internal/testing"
	"testing"
)

type RepositorySuite struct{ suite.Suite }

func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}

func (s *RepositorySuite) Test() {
	log := internalLog.New(io.Discard)

	s.Run("Path", func() {
		repository := &Repository{
			path: "path",
		}

		s.Equal("path", repository.Path())
	})

	s.Run("Source", func() {
		repository := &Repository{
			path: "path",
		}

		s.Equal("path", repository.Source())
	})

	s.Run("Dir", func() {
		repository := &Repository{
			dir: "dir",
		}

		s.Equal("dir", repository.Dir())
	})

	s.Run("LoadRecipe", func() {
		repository := &Repository{
			log: log,
			dir: internalTesting.DataPath(s, "repository"),
		}
		repository.recipeLoader = &RecipeRepositoryDirLoader{
			Log:        log,
			Repository: repository,
		}

		recipe, err := repository.LoadRecipe("recipe")

		s.NoError(err)
		s.Equal(internalTesting.DataPath(s, "repository", "recipe"), recipe.Path())
		s.Equal("recipe", recipe.Name())
		s.Equal("description", recipe.Description())
		s.Equal(map[string]interface{}{"foo": "bar"}, recipe.Vars())
	})

	s.Run("WalkRecipes Empty", func() {
		repository := &Repository{
			log: log,
			dir: internalTesting.DataPath(s, "repository"),
		}
		repository.recipeLoader = &RecipeRepositoryDirLoader{
			Log:        log,
			Repository: repository,
		}

		err := repository.WalkRecipes(func(recipe *Recipe) {})

		s.ErrorAs(err, &internalError)
		s.Equal("empty repository", internalError.Message)
	})

	s.Run("WalkRecipes", func() {
		repository := &Repository{
			log: log,
			dir: internalTesting.DataPath(s, "repository"),
		}
		repository.recipeLoader = &RecipeRepositoryDirLoader{
			Log:        log,
			Repository: repository,
		}

		count := 1

		err := repository.WalkRecipes(func(recipe *Recipe) {
			switch count {
			case 1:
				s.Equal(internalTesting.DataPath(s, "repository", "bar"), recipe.Path())
				s.Equal("bar", recipe.Name())
				s.Equal("Bar", recipe.Description())
				s.Equal(map[string]interface{}{"bar": "bar"}, recipe.Vars())
			case 2:
				s.Equal(internalTesting.DataPath(s, "repository", "foo"), recipe.Path())
				s.Equal("foo", recipe.Name())
				s.Equal("Foo", recipe.Description())
				s.Equal(map[string]interface{}{"foo": "foo"}, recipe.Vars())
			}

			count++
		})

		s.NoError(err)
	})
}
