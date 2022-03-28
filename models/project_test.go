package models

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

/*******************/
/* Project - Suite */
/*******************/

type ProjectTestSuite struct {
	suite.Suite
	recipeVars map[string]interface{}
	recipe     RecipeInterface
}

func TestProjectTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(ProjectTestSuite))
}

func (s *ProjectTestSuite) SetupTest() {
	s.recipeVars = map[string]interface{}{
		"foo": "foo",
		"bar": "foo",
	}
	s.recipe = NewRecipe(
		"foo",
		"bar",
		"",
		"baz",
		NewRepository(
			"foo",
			"bar",
		),
		s.recipeVars,
		nil,
		nil,
		nil,
	)
}

/*******************/
/* Project - Tests */
/*******************/

func (s *ProjectTestSuite) TestProject() {
	dir := "foo"

	s.Run("New", func() {
		prj := NewProject(dir, s.recipe, map[string]interface{}{
			"bar": "bar",
			"baz": "baz",
		})

		s.Implements((*ProjectInterface)(nil), prj)
		s.Equal(dir, prj.getDir())
		s.Equal(s.recipe, prj.Recipe())
		s.Equal(map[string]interface{}{
			"foo": "foo",
			"bar": "bar",
			"baz": "baz",
		}, prj.Vars())

		// Recipe vars should stay untouched
		s.Equal(s.recipeVars, s.recipe.Vars())
	})
}
