package validator

import (
	"github.com/stretchr/testify/suite"
	"manala/models"
	"testing"
)

/**************************/
/* Validate Value - Suite */
/**************************/

type ValidateValueTestSuite struct {
	suite.Suite
}

func TestValidateValueTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(ValidateValueTestSuite))
}

/**************************/
/* Validate Value - Tests */
/**************************/

func (s *ValidateValueTestSuite) TestValidateValue() {
	err := ValidateValue(
		"foo.bar",
		map[string]interface{}{"type": "string", "format": "domain"},
	)
	s.NoError(err)
}

func (s *ValidateValueTestSuite) TestValidateValueError() {
	err := ValidateValue(
		"foo",
		map[string]interface{}{"type": "string", "format": "domain"},
	)
	s.Error(err)
	s.Equal("\n- Does not match format 'domain'", err.Error())
}

/****************************/
/* Validate Project - Suite */
/****************************/

type ValidateProjectTestSuite struct {
	suite.Suite
	recipe models.RecipeInterface
}

func TestValidateProjectTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(ValidateProjectTestSuite))
}

func (s *ValidateProjectTestSuite) SetupTest() {
	s.recipe = models.NewRecipe(
		"foo",
		"bar",
		"",
		"baz",
		models.NewRepository(
			"foo",
			"bar",
		),
		nil,
		nil,
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"foo": map[string]interface{}{
					"type": "string",
				},
			},
		},
		nil,
	)
}

/****************************/
/* Validate Project - Tests */
/****************************/

func (s *ValidateProjectTestSuite) TestValidateProject() {
	dir := "foo"

	s.Run("Success", func() {
		prj := models.NewProject(dir, s.recipe, map[string]interface{}{
			"foo": "bar",
		})
		err := ValidateProject(prj)
		s.NoError(err)
	})

	s.Run("Error", func() {
		prj := models.NewProject(dir, s.recipe, map[string]interface{}{
			"foo": 123,
		})
		err := ValidateProject(prj)
		s.Error(err)
		s.Equal("project config errors:\n- foo: Invalid type. Expected: string, given: integer", err.Error())
	})
}

/**************************/
/* Format Checker - Suite */
/**************************/

type FormatCheckerTestSuite struct {
	suite.Suite
}

func TestFormatCheckerTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(FormatCheckerTestSuite))
}

/**************************/
/* Format Checker - Tests */
/**************************/

func (s *FormatCheckerTestSuite) TestGitRepo() {
	checker := GitRepoFormatChecker{}

	s.True(checker.IsFormat("git@github.com:manala/manala.git"))

	s.False(checker.IsFormat("foo"))
}

func (s *FormatCheckerTestSuite) TestFilePath() {
	checker := FilePathFormatChecker{}

	s.True(checker.IsFormat("/"))
	s.True(checker.IsFormat("/foo"))
	s.True(checker.IsFormat("/foo/bar"))
	s.True(checker.IsFormat("/foo-bar"))

	s.False(checker.IsFormat("foo"))
	s.False(checker.IsFormat("/foo/*"))
}

func (s *FormatCheckerTestSuite) TestDomain() {
	checker := DomainFormatChecker{}

	s.True(checker.IsFormat("foo.bar"))
	s.True(checker.IsFormat("foo.bar.baz"))
	s.True(checker.IsFormat("foo-bar.baz"))
	s.True(checker.IsFormat("fo.ba"))

	s.False(checker.IsFormat("foo"))
	s.False(checker.IsFormat("foo.b"))
	s.False(checker.IsFormat("foo_.bar"))
}
