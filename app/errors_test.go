package app

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type ErrorsSuite struct{ suite.Suite }

func TestErrorsSuite(t *testing.T) {
	suite.Run(t, new(ErrorsSuite))
}

func (s *ErrorsSuite) TestProject() {
	s.Run("AlreadyExistingProjectError", func() {
		err := &AlreadyExistingProjectError{Dir: "dir"}

		s.EqualError(err, "already existing project")
		s.Equal([]any{"dir", "dir"}, err.ErrorArguments())
	})

	s.Run("NotFoundProjectManifestError", func() {
		err := &NotFoundProjectManifestError{File: "file"}

		s.EqualError(err, "project manifest not found")
		s.Equal([]any{"file", "file"}, err.ErrorArguments())
	})
}

func (s *ErrorsSuite) TestRecipe() {
	s.Run("NotFoundRecipeManifestError", func() {
		err := &NotFoundRecipeManifestError{File: "file"}

		s.EqualError(err, "recipe manifest not found")
		s.Equal([]any{"file", "file"}, err.ErrorArguments())
	})

	s.Run("UnprocessableRecipeNameError", func() {
		err := &UnprocessableRecipeNameError{}

		s.EqualError(err, "unable to process recipe name")
	})
}

func (s *ErrorsSuite) TestRepository() {
	s.Run("UnsupportedRepositoryError", func() {
		err := &UnsupportedRepositoryError{Url: "url"}

		s.EqualError(err, "unsupported repository url")
		s.Equal([]any{"url", "url"}, err.ErrorArguments())
	})

	s.Run("UnprocessableRepositoryUrlError", func() {
		err := &UnprocessableRepositoryUrlError{}

		s.EqualError(err, "unable to process repository url")
	})
}
