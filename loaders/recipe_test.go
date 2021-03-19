package loaders

import (
	"github.com/stretchr/testify/suite"
	"manala/fs"
	"manala/logger"
	"manala/models"
	"testing"
)

/******************/
/* Recipe - Suite */
/******************/

type RecipeTestSuite struct {
	suite.Suite
	ld                      RecipeLoaderInterface
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
	s.repository = models.NewRepository("testdata/recipe/_repository", "testdata/recipe/_repository", false)
	s.repositoryEmpty = models.NewRepository("testdata/recipe/_repository_empty", "testdata/recipe/_repository_empty", false)
	s.repositoryInvalid = models.NewRepository("testdata/recipe/_repository_invalid", "testdata/recipe/_repository_invalid", false)
	s.repositoryIncorrect = models.NewRepository("testdata/recipe/_repository_incorrect", "testdata/recipe/_repository_incorrect", false)
	s.repositoryNoDescription = models.NewRepository("testdata/recipe/_repository_no_description", "testdata/recipe/_repository_no_description", false)
	s.repositorySchemaInvalid = models.NewRepository("testdata/recipe/_repository_schema_invalid", "testdata/recipe/_repository_schema_invalid", false)

	log := logger.New(logger.WithDiscardment())

	fsManager := fs.NewManager()
	modelFsManager := models.NewFsManager(fsManager)

	s.ld = NewRecipeLoader(log, modelFsManager)
}

/******************/
/* Recipe - Tests */
/******************/

func (s *RecipeTestSuite) TestRecipe() {
	s.Implements((*RecipeLoaderInterface)(nil), s.ld)
}

func (s *RecipeTestSuite) TestRecipeLoad() {
	rec, err := s.ld.Load("load", s.repository)
	s.NoError(err)
	s.Implements((*models.RecipeInterface)(nil), rec)
	s.Equal("load", rec.Name())
	s.Equal("Load", rec.Description())
	s.Equal(s.repository, rec.Repository())
}

func (s *RecipeTestSuite) TestRecipeLoadNotFound() {
	rec, err := s.ld.Load("not_found", s.repository)
	s.Error(err)
	s.Equal("recipe not found", err.Error())
	s.Nil(rec)
}

func (s *RecipeTestSuite) TestRecipeLoadEmpty() {
	rec, err := s.ld.Load("load", s.repositoryEmpty)
	s.Error(err)
	s.Equal("empty recipe manifest \"load\"", err.Error())
	s.Nil(rec)
}

func (s *RecipeTestSuite) TestRecipeLoadInvalid() {
	rec, err := s.ld.Load("load", s.repositoryInvalid)
	s.Error(err)
	s.Equal("invalid recipe manifest \"load\" (yaml: mapping values are not allowed in this context)", err.Error())
	s.Nil(rec)
}

func (s *RecipeTestSuite) TestRecipeLoadIncorrect() {
	rec, err := s.ld.Load("load", s.repositoryIncorrect)
	s.Error(err)
	s.Equal("incorrect recipe manifest \"load\" (yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `foo` into map[string]interface {})", err.Error())
	s.Nil(rec)
}

func (s *RecipeTestSuite) TestRecipeLoadNoDescription() {
	rec, err := s.ld.Load("load", s.repositoryNoDescription)
	s.Error(err)
	s.Equal("Key: 'recipeConfig.Description' Error:Field validation for 'Description' failed on the 'required' tag", err.Error())
	s.Nil(rec)
}

func (s *RecipeTestSuite) TestRecipeLoadVars() {
	rec, err := s.ld.Load("load_vars", s.repository)
	s.NoError(err)
	s.Equal(
		map[string]interface{}{
			"foo": map[string]interface{}{"foo": "bar", "bar": "baz", "baz": []interface{}{}},
			"bar": map[string]interface{}{"bar": "baz"},
			"baz": map[string]interface{}{"bar": "baz", "baz": "qux"},
		},
		rec.Vars(),
	)
}

func (s *RecipeTestSuite) TestRecipeLoadSync() {
	rec, err := s.ld.Load("load_sync", s.repository)
	s.NoError(err)
	s.Equal(
		[]models.RecipeSyncUnit{
			{Source: "foo", Destination: "foo"},
			{Source: "foo", Destination: "bar"},
		},
		rec.Sync(),
	)
}

func (s *RecipeTestSuite) TestRecipeLoadSchema() {
	rec, err := s.ld.Load("load_schema", s.repository)
	s.NoError(err)
	s.Equal(
		map[string]interface{}{
			"type":                 "object",
			"additionalProperties": false,
			"properties": map[string]interface{}{
				"foo": map[string]interface{}{
					"type":                 "object",
					"additionalProperties": false,
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
					"type":                 "object",
					"additionalProperties": false,
					"properties": map[string]interface{}{
						"bar": map[string]interface{}{},
					},
				},
				"additionalProperties": map[string]interface{}{
					"type":                 "object",
					"additionalProperties": false,
					"properties": map[string]interface{}{
						"object": map[string]interface{}{
							"type":                 "object",
							"additionalProperties": false,
							"properties": map[string]interface{}{
								"foo": map[string]interface{}{},
								"bar": map[string]interface{}{},
							},
						},
						"object_overriden": map[string]interface{}{
							"type":                 "object",
							"additionalProperties": true,
							"properties": map[string]interface{}{
								"foo": map[string]interface{}{},
								"bar": map[string]interface{}{},
							},
						},
						"empty_object": map[string]interface{}{
							"type":                 "object",
							"additionalProperties": true,
							"properties":           map[string]interface{}{},
						},
						"empty_object_overriden": map[string]interface{}{
							"type":                 "object",
							"additionalProperties": false,
							"properties":           map[string]interface{}{},
						},
					},
				},
			},
		},
		rec.Schema(),
	)
}

func (s *RecipeTestSuite) TestRecipeLoadSchemaInvalid() {
	rec, err := s.ld.Load("load", s.repositorySchemaInvalid)
	s.Error(err)
	s.Equal("invalid recipe schema tag at \"/foo\": unexpected end of JSON input", err.Error())
	s.Nil(rec)
}

func (s *RecipeTestSuite) TestRecipeLoadOptions() {
	rec, err := s.ld.Load("load_options", s.repository)
	s.NoError(err)
	s.Equal(
		[]models.RecipeOption{
			{Label: "Foo bar", Path: "/foo/bar", Schema: map[string]interface{}{"enum": []interface{}{nil, "baz"}}},
		},
		rec.Options(),
	)
}

func (s *RecipeTestSuite) TestRecipeWalk() {
	repo := models.NewRepository("testdata/recipe/walk", "testdata/recipe/walk", false)
	results := make(map[string]string)
	err := s.ld.Walk(repo, func(rec models.RecipeInterface) {
		results[rec.Name()] = rec.Description()
	})
	s.NoError(err)
	s.Len(results, 1)
	s.Equal("Default", results["default"])
}
