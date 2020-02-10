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
	repositoryInvalid       models.RepositoryInterface
	repositoryIncorrect     models.RepositoryInterface
	repositoryNoDescription models.RepositoryInterface
	repositorySchemaInvalid models.RepositoryInterface
}

func TestRecipeTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(RecipeTestSuite))
}

func (s *RecipeTestSuite) SetupTest() {
	s.repository = models.NewRepository("testdata/recipe/_repository", "testdata/recipe/_repository")
	s.repositoryEmpty = models.NewRepository("testdata/recipe/_repository_empty", "testdata/recipe/_repository_empty")
	s.repositoryInvalid = models.NewRepository("testdata/recipe/_repository_invalid", "testdata/recipe/_repository_invalid")
	s.repositoryIncorrect = models.NewRepository("testdata/recipe/_repository_incorrect", "testdata/recipe/_repository_incorrect")
	s.repositoryNoDescription = models.NewRepository("testdata/recipe/_repository_no_description", "testdata/recipe/_repository_no_description")
	s.repositorySchemaInvalid = models.NewRepository("testdata/recipe/_repository_schema_invalid", "testdata/recipe/_repository_schema_invalid")
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
	file, err := ld.ConfigFile("testdata/recipe/config_file")
	s.NoError(err)
	s.IsType((*os.File)(nil), file)
	s.Equal("testdata/recipe/config_file/.manala.yaml", file.Name())
}

func (s *RecipeTestSuite) TestRecipeConfigFileNotFound() {
	ld := NewRecipeLoader()
	file, err := ld.ConfigFile("testdata/recipe/config_file_not_found")
	s.Error(err)
	s.Equal("open testdata/recipe/config_file_not_found/.manala.yaml: no such file or directory", err.Error())
	s.Nil(file)
}

func (s *RecipeTestSuite) TestRecipeConfigFileDirectory() {
	ld := NewRecipeLoader()
	file, err := ld.ConfigFile("testdata/recipe/config_file_directory")
	s.Error(err)
	s.Equal("open testdata/recipe/config_file_directory/.manala.yaml: is a directory", err.Error())
	s.Nil(file)
}

func (s *RecipeTestSuite) TestRecipeLoad() {
	ld := NewRecipeLoader()
	rec, err := ld.Load("load", s.repository)
	s.NoError(err)
	s.Implements((*models.RecipeInterface)(nil), rec)
	s.Equal("load", rec.Name())
	s.Equal("Load", rec.Description())
	s.Equal("testdata/recipe/_repository/load", rec.Dir())
	s.Equal(s.repository, rec.Repository())
}

func (s *RecipeTestSuite) TestRecipeLoadNotFound() {
	ld := NewRecipeLoader()
	rec, err := ld.Load("not_found", s.repository)
	s.Error(err)
	s.Equal("recipe not found", err.Error())
	s.Nil(rec)
}

func (s *RecipeTestSuite) TestRecipeLoadEmpty() {
	ld := NewRecipeLoader()
	rec, err := ld.Load("load", s.repositoryEmpty)
	s.Error(err)
	s.Equal("empty recipe config \"testdata/recipe/_repository_empty/load/.manala.yaml\"", err.Error())
	s.Nil(rec)
}

func (s *RecipeTestSuite) TestRecipeLoadInvalid() {
	ld := NewRecipeLoader()
	rec, err := ld.Load("load", s.repositoryInvalid)
	s.Error(err)
	s.Equal("invalid recipe config \"testdata/recipe/_repository_invalid/load/.manala.yaml\" (yaml: mapping values are not allowed in this context)", err.Error())
	s.Nil(rec)
}

func (s *RecipeTestSuite) TestRecipeLoadIncorrect() {
	ld := NewRecipeLoader()
	rec, err := ld.Load("load", s.repositoryIncorrect)
	s.Error(err)
	s.Equal("incorrect recipe config \"testdata/recipe/_repository_incorrect/load/.manala.yaml\" (yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `foo` into map[string]interface {})", err.Error())
	s.Nil(rec)
}

func (s *RecipeTestSuite) TestRecipeLoadNoDescription() {
	ld := NewRecipeLoader()
	rec, err := ld.Load("load", s.repositoryNoDescription)
	s.Error(err)
	s.Equal("Key: 'recipeConfig.Description' Error:Field validation for 'Description' failed on the 'required' tag", err.Error())
	s.Nil(rec)
}

func (s *RecipeTestSuite) TestRecipeLoadVars() {
	ld := NewRecipeLoader()
	rec, _ := ld.Load("load_vars", s.repository)
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
	rec, _ := ld.Load("load_sync_units", s.repository)
	s.Equal(
		[]models.RecipeSyncUnit{
			{Source: "foo", Destination: "foo"},
			{Source: "foo", Destination: "bar"},
		},
		rec.SyncUnits(),
	)
}

func (s *RecipeTestSuite) TestRecipeLoadSchema() {
	ld := NewRecipeLoader()
	rec, _ := ld.Load("load_schema", s.repository)
	s.Equal(
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"foo": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"foo": map[string]interface{}{},
						"bar": map[string]interface{}{
							"enum": []interface{}{
								nil,
								"baz",
							},
						},
						"baz": map[string]interface{}{
							"type": "array",
						},
					},
					"required": []interface{}{
						"foo",
						"bar",
					},
				},
				"bar": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"bar": map[string]interface{}{},
					},
				},
			},
		},
		rec.Schema(),
	)
}

func (s *RecipeTestSuite) TestRecipeLoadSchemaInvalid() {
	ld := NewRecipeLoader()
	rec, err := ld.Load("load", s.repositorySchemaInvalid)
	s.Error(err)
	s.Equal("invalid recipe schema tag at \"/foo\": unexpected end of JSON input", err.Error())
	s.Nil(rec)
}

func (s *RecipeTestSuite) TestRecipeWalk() {
	ld := NewRecipeLoader()
	results := make(map[string]string)
	err := ld.Walk(s.repository, func(rec models.RecipeInterface) {
		results[rec.Name()] = rec.Description()
	})
	s.NoError(err)
	s.Len(results, 4)
	s.Equal("Load", results["load"])
	s.Equal("Load vars", results["load_vars"])
	s.Equal("Load sync units", results["load_sync_units"])
	s.Equal("Load schema", results["load_schema"])
}
