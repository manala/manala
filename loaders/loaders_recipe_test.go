package loaders

import (
	"github.com/stretchr/testify/suite"
	"manala/models"
	"os"
	"testing"
)

/******************/
/* Recipe - Suite */
/******************/

type RecipeTestSuite struct {
	suite.Suite
	repository              models.RepositoryInterface
	repositoryEmpty         models.RepositoryInterface
	repositoryIncorrect     models.RepositoryInterface
	repositoryNoDescription models.RepositoryInterface
}

func TestRecipeTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(RecipeTestSuite))
}

func (s *RecipeTestSuite) SetupTest() {
	s.repository = models.NewRepository(
		"testdata/repository",
		"testdata/repository",
	)
	s.repositoryEmpty = models.NewRepository(
		"testdata/repository_empty",
		"testdata/repository_empty",
	)
	s.repositoryIncorrect = models.NewRepository(
		"testdata/repository_incorrect",
		"testdata/repository_incorrect",
	)
	s.repositoryNoDescription = models.NewRepository(
		"testdata/repository_no_description",
		"testdata/repository_no_description",
	)
}

/******************/
/* Recipe - Tests */
/******************/

func (s *RecipeTestSuite) TestRecipe() {
	ld := NewRecipeLoader()
	s.Implements((*RecipeLoaderInterface)(nil), ld)
}

func (s *RecipeTestSuite) TestRecipeConfigFile() {
	ld := NewRecipeLoader()
	file, err := ld.ConfigFile("testdata/recipe")
	s.NoError(err)
	s.IsType((*os.File)(nil), file)
	s.Equal("testdata/recipe/.manala.yaml", file.Name())
}

func (s *RecipeTestSuite) TestRecipeConfigFileNotFound() {
	ld := NewRecipeLoader()
	file, err := ld.ConfigFile("testdata/recipe_no_config_file")
	s.Error(err)
	s.Equal("open testdata/recipe_no_config_file/.manala.yaml: no such file or directory", err.Error())
	s.Nil(file)
}

func (s *RecipeTestSuite) TestRecipeConfigFileDirectory() {
	ld := NewRecipeLoader()
	file, err := ld.ConfigFile("testdata/recipe_config_directory")
	s.Error(err)
	s.Equal("open testdata/recipe_config_directory/.manala.yaml: is a directory", err.Error())
	s.Nil(file)
}

func (s *RecipeTestSuite) TestRecipeLoad() {
	ld := NewRecipeLoader()
	rec, err := ld.Load("foo", s.repository)
	s.NoError(err)
	s.Implements((*models.RecipeInterface)(nil), rec)
	s.Equal("foo", rec.Name())
	s.Equal("Foo bar", rec.Description())
	s.Equal("testdata/repository/foo", rec.Dir())
	s.Equal(s.repository, rec.Repository())
}

func (s *RecipeTestSuite) TestRecipeLoadNotExist() {
	ld := NewRecipeLoader()
	rec, err := ld.Load("not_exist", s.repository)
	s.Error(err)
	s.Equal("recipe not found", err.Error())
	s.Nil(rec)
}

func (s *RecipeTestSuite) TestRecipeLoadEmpty() {
	ld := NewRecipeLoader()
	rec, err := ld.Load("foo", s.repositoryEmpty)
	s.Error(err)
	s.Equal("empty recipe config \"testdata/repository_empty/foo/.manala.yaml\"", err.Error())
	s.Nil(rec)
}

func (s *RecipeTestSuite) TestRecipeLoadIncorrect() {
	ld := NewRecipeLoader()
	rec, err := ld.Load("foo", s.repositoryIncorrect)
	s.Error(err)
	s.Equal("invalid recipe config \"testdata/repository_incorrect/foo/.manala.yaml\" (yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `foo` into map[string]interface {})", err.Error())
	s.Nil(rec)
}

func (s *RecipeTestSuite) TestRecipeLoadNoDescription() {
	ld := NewRecipeLoader()
	rec, err := ld.Load("foo", s.repositoryNoDescription)
	s.Error(err)
	s.Equal("Key: 'recipeConfig.Description' Error:Field validation for 'Description' failed on the 'required' tag", err.Error())
	s.Nil(rec)
}

func (s *RecipeTestSuite) TestRecipeLoadVars() {
	ld := NewRecipeLoader()
	rec, _ := ld.Load("foo", s.repository)
	s.Equal(
		map[string]interface{}{
			"foo": map[string]interface{}{"foo": "bar", "bar": "baz", "baz": []interface{}{}},
			"bar": map[string]interface{}{"bar": "baz"},
			"baz": map[string]interface{}{"bar": "baz", "baz": "qux"},
		},
		rec.Vars(),
	)
}

func (s *RecipeTestSuite) TestRecipeLoadSyncUnits() {
	ld := NewRecipeLoader()
	rec, _ := ld.Load("foo", s.repository)
	s.Equal(
		[]models.RecipeSyncUnit{
			{Source: "foo", Destination: "foo"},
			{Source: "foo", Destination: "bar"},
		},
		rec.SyncUnits(),
	)
}

func (s *RecipeTestSuite) TestRecipeWalk() {
	ld := NewRecipeLoader()
	results := make(map[string]string)
	err := ld.Walk(s.repository, func(rec models.RecipeInterface) {
		results[rec.Name()] = rec.Description()
	})
	s.NoError(err)
	s.Len(results, 3)
	s.Equal("Foo bar", results["foo"])
	s.Equal("Bar bar", results["bar"])
	s.Equal("Baz bar", results["baz"])
}
