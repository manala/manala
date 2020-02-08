package validator

import (
	"github.com/stretchr/testify/suite"
	"manala/models"
	"testing"
)

/****************************/
/* Validate Project - Suite */
/****************************/

type ValidateProjectTestSuite struct {
	suite.Suite
	project models.ProjectInterface
}

func TestValidateProjectTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(ValidateProjectTestSuite))
}

func (s *ValidateProjectTestSuite) SetupTest() {
	s.project = models.NewProject(
		"foo",
		models.NewRecipe(
			"foo",
			"bar",
			"baz",
			models.NewRepository(
				"foo",
				"bar",
			),
		),
	)
}

/****************************/
/* Validate Project - Tests */
/****************************/

func (s *ValidateProjectTestSuite) TestValidateProject() {
	s.project.Recipe().MergeSchema(
		&map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"foo": map[string]interface{}{
					"type": "string",
				},
			},
		},
	)
	s.project.MergeVars(
		&map[string]interface{}{
			"foo": "bar",
		},
	)
	err := ValidateProject(s.project)
	s.NoError(err)
}

func (s *ValidateProjectTestSuite) TestValidateProjectErrors() {
	s.project.Recipe().MergeSchema(
		&map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"foo": map[string]interface{}{
					"type": "string",
				},
			},
		},
	)
	s.project.MergeVars(
		&map[string]interface{}{
			"foo": 123,
		},
	)
	err := ValidateProject(s.project)
	s.Error(err)
	s.Equal("project config errors:\n- foo: Invalid type. Expected: string, given: integer", err.Error())
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

	s.False(checker.IsFormat("foo"))
	s.False(checker.IsFormat("foo.b"))
	s.False(checker.IsFormat("foo_.bar"))
}
