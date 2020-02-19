package models

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

/******************/
/* Recipe - Suite */
/******************/

type RecipeTestSuite struct {
	suite.Suite
	name        string
	description string
	dir         string
	repository  RepositoryInterface
}

func TestRecipeTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(RecipeTestSuite))
}

func (s *RecipeTestSuite) SetupTest() {
	s.name = "foo"
	s.description = "bar"
	s.dir = "baz"
	s.repository = NewRepository(
		"foo",
		"bar",
	)
}

/******************/
/* Recipe - Tests */
/******************/

func (s *RecipeTestSuite) TestRecipe() {
	rec := NewRecipe(s.name, s.description, s.dir, s.repository)
	s.Implements((*RecipeInterface)(nil), rec)
	s.Equal(s.name, rec.Name())
	s.Equal(s.description, rec.Description())
	s.Equal(s.dir, rec.Dir())
	s.Equal(s.repository, rec.Repository())
	s.Len(rec.Vars(), 0)
	s.Len(rec.SyncUnits(), 0)
	s.Len(rec.Schema(), 0)
}

func (s *RecipeTestSuite) TestRecipeVars() {
	rec := NewRecipe(s.name, s.description, s.dir, s.repository)
	vars := map[string]interface{}{
		"foo": "bar",
		"bar": "bas",
	}
	rec.MergeVars(&vars)
	s.Equal(vars, rec.Vars())
}

func (s *RecipeTestSuite) TestRecipeSyncUnits() {
	rec := NewRecipe(s.name, s.description, s.dir, s.repository)
	syncUnits := []RecipeSyncUnit{
		{Source: "foo", Destination: "bar"},
		{Source: "bar", Destination: "baz"},
	}
	rec.AddSyncUnits(syncUnits)
	s.Equal(syncUnits, rec.SyncUnits())
}

func (s *RecipeTestSuite) TestRecipeSchema() {
	rec := NewRecipe(s.name, s.description, s.dir, s.repository)
	schema := map[string]interface{}{
		"foo": "bar",
		"bar": "bas",
	}
	rec.MergeVars(&schema)
	s.Equal(schema, rec.Vars())
}

func (s *RecipeTestSuite) TestRecipeOptions() {
	rec := NewRecipe(s.name, s.description, s.dir, s.repository)
	s.False(rec.HasOptions())
	options := []RecipeOption{
		{Label: "foo", Path: "bar", Schema: map[string]interface{}{"baz": "qux"}},
		{Label: "bar", Path: "baz", Schema: map[string]interface{}{"qux": "quobuz"}},
	}
	rec.AddOptions(options)
	s.True(rec.HasOptions())
	s.Equal(options, rec.Options())
}
