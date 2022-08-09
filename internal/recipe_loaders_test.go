package internal

import (
	"github.com/stretchr/testify/suite"
	"io"
	internalLog "manala/internal/log"
	internalTesting "manala/internal/testing"
	"testing"
)

type RecipeLoadersSuite struct{ suite.Suite }

func TestRecipeLoadersSuite(t *testing.T) {
	suite.Run(t, new(RecipeLoadersSuite))
}

func (s *RecipeLoadersSuite) TestRepositoryDir() {
	log := internalLog.New(io.Discard)

	s.Run("LoadRecipeManifest Not Found", func() {
		loader := &RecipeRepositoryDirLoader{
			Log: log,
			Repository: &Repository{
				path: internalTesting.DataPath(s, "repository"),
				dir:  internalTesting.DataPath(s, "repository"),
			},
		}

		recipeManifest, err := loader.LoadRecipeManifest("recipe")

		var _err *NotFoundRecipeManifestError
		s.ErrorAs(err, &_err)

		s.ErrorAs(err, &internalError)
		s.Equal("recipe manifest not found", internalError.Message)
		s.Nil(recipeManifest)
	})

	s.Run("LoadRecipeManifest Wrong", func() {
		loader := &RecipeRepositoryDirLoader{
			Log: log,
			Repository: &Repository{
				path: internalTesting.DataPath(s, "repository"),
				dir:  internalTesting.DataPath(s, "repository"),
			},
		}

		recipeManifest, err := loader.LoadRecipeManifest("recipe")

		s.ErrorAs(err, &internalError)
		s.Equal("wrong recipe manifest", internalError.Message)
		s.Nil(recipeManifest)
	})

	s.Run("LoadRecipeManifest", func() {
		loader := &RecipeRepositoryDirLoader{
			Log: log,
			Repository: &Repository{
				path: internalTesting.DataPath(s, "repository"),
				dir:  internalTesting.DataPath(s, "repository"),
			},
		}

		recipeManifest, err := loader.LoadRecipeManifest("recipe")

		s.NoError(err)
		s.Equal("Recipe", recipeManifest.Description)
		s.Equal(map[string]interface{}{"foo": "bar"}, recipeManifest.Vars)
	})

	s.Run("LoadRecipe", func() {
		loader := &RecipeRepositoryDirLoader{
			Log: log,
			Repository: &Repository{
				path: internalTesting.DataPath(s, "repository"),
				dir:  internalTesting.DataPath(s, "repository"),
			},
		}

		recipe, err := loader.LoadRecipe("recipe")

		s.NoError(err)
		s.Equal(internalTesting.DataPath(s, "repository", "recipe"), recipe.Path())
		s.Equal(internalTesting.DataPath(s, "repository"), recipe.Repository().Path())
	})
}
