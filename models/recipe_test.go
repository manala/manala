package models

import (
	"github.com/stretchr/testify/suite"
	"path/filepath"
	"testing"
)

/******************/
/* Recipe - Suite */
/******************/

type RecipeTestSuite struct {
	suite.Suite
	repository RepositoryInterface
}

func TestRecipeTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(RecipeTestSuite))
}

func (s *RecipeTestSuite) SetupTest() {
	s.repository = NewRepository(
		"foo",
		"bar",
	)
}

/******************/
/* Recipe - Tests */
/******************/

func (s *RecipeTestSuite) TestRecipe() {
	name := "foo"
	description := "bar"
	template := "baz"
	dir := "qux"
	vars := map[string]interface{}{
		"foo": "foo",
		"bar": "bar",
	}
	sync := []RecipeSyncUnit{
		{Source: "foo", Destination: "bar"},
		{Source: "bar", Destination: "baz"},
	}
	schema := map[string]interface{}{
		"foo": "bar",
		"bar": "bas",
	}
	options := []RecipeOption{
		{Label: "foo", Path: "bar", Schema: map[string]interface{}{"baz": "qux"}},
		{Label: "bar", Path: "baz", Schema: map[string]interface{}{"qux": "quobuz"}},
	}

	s.Run("New", func() {
		rec := NewRecipe(name, description, template, dir, s.repository, vars, sync, schema, options)
		s.Implements((*RecipeInterface)(nil), rec)
		s.Equal(name, rec.Name())
		s.Equal(description, rec.Description())
		s.Equal(template, rec.Template())
		s.Equal(filepath.Join(s.repository.getDir(), dir), rec.getDir())
		s.Equal(s.repository, rec.Repository())
		s.Equal(vars, rec.Vars())
		s.Equal(sync, rec.Sync())
		s.Equal(schema, rec.Schema())
		s.Equal(options, rec.Options())
	})
}
