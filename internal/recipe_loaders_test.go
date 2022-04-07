package internal

import (
	"github.com/stretchr/testify/suite"
	"io"
	internalLog "manala/internal/log"
	"path/filepath"
	"testing"
)

type RecipeLoadersSuite struct{ suite.Suite }

func TestRecipeLoadersSuite(t *testing.T) {
	suite.Run(t, new(RecipeLoadersSuite))
}

var recipeLoadersTestPath = filepath.Join("testdata", "recipe_loaders")

func (s *RecipeLoadersSuite) TestRepositoryDir() {
	log := internalLog.New(io.Discard)
	loader := &RecipeRepositoryDirLoader{
		Log:        log,
		Repository: &Repository{path: recipeLoadersTestPath, dir: recipeLoadersTestPath},
	}

	s.Run("LoadRecipeManifest Not Found", func() {
		recipeManifest, err := loader.LoadRecipeManifest("recipe_not_found")

		var _err *NotFoundRecipeManifestError
		s.ErrorAs(err, &_err)

		s.ErrorAs(err, &internalError)
		s.Equal("recipe manifest not found", internalError.Message)
		s.Nil(recipeManifest)
	})

	s.Run("LoadRecipeManifest Wrong", func() {
		recipeManifest, err := loader.LoadRecipeManifest("recipe_wrong")

		s.ErrorAs(err, &internalError)
		s.Equal("wrong recipe manifest", internalError.Message)
		s.Nil(recipeManifest)
	})

	s.Run("LoadRecipeManifest", func() {
		recipeManifest, err := loader.LoadRecipeManifest("recipe")

		s.NoError(err)
		s.Equal("description", recipeManifest.Description)
		s.Equal(map[string]interface{}{"foo": "bar"}, recipeManifest.Vars)
	})

	s.Run("LoadRecipe", func() {
		recipe, err := loader.LoadRecipe("recipe")

		s.NoError(err)
		s.Equal(filepath.Join(recipeLoadersTestPath, "recipe"), recipe.Path())
		s.Equal(recipeLoadersTestPath, recipe.Repository().Path())
	})
}
