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
	dir    string
	recipe RecipeInterface
}

func TestProjectTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(ProjectTestSuite))
}

func (s *ProjectTestSuite) SetupTest() {
	s.dir = "foo"
	s.recipe = NewRecipe(
		"foo",
		"bar",
		"baz",
		NewRepository(
			"foo",
			"bar",
		),
	)
	s.recipe.MergeVars(
		&map[string]interface{}{
			"foo": "bar",
			"bar": "baz",
		},
	)
}

/*******************/
/* Project - Tests */
/*******************/

func (s *ProjectTestSuite) TestProject() {
	prj := NewProject(s.dir, s.recipe)
	s.Implements((*ProjectInterface)(nil), prj)
	s.Equal(s.dir, prj.Dir())
	s.Equal(s.recipe, prj.Recipe())
	s.Equal(s.recipe.Vars(), prj.Vars())
}

func (s *ProjectTestSuite) TestProjectMergeVars() {
	prj := NewProject(s.dir, s.recipe)
	prj.MergeVars(
		&map[string]interface{}{
			"bar": "qux",
			"baz": "qux",
		},
	)
	s.Equal(
		map[string]interface{}{
			"foo": "bar",
			"bar": "qux",
			"baz": "qux",
		},
		prj.Vars(),
	)
}
