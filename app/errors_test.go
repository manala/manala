package app

import (
	"testing"

	"github.com/stretchr/testify/suite"
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

	s.Run("NotFoundProjectError", func() {
		err := &NotFoundProjectError{Dir: "dir"}

		s.EqualError(err, "project not found")
		s.Equal([]any{"dir", "dir"}, err.ErrorArguments())
	})
}

func (s *ErrorsSuite) TestRecipe() {
	s.Run("NotFoundRecipeError", func() {
		repositoryMock := &RepositoryMock{}
		repositoryMock.
			On("URL").Return("url")

		err := &NotFoundRecipeError{Repository: repositoryMock, Name: "name"}

		s.EqualError(err, "recipe not found")
		s.Equal([]any{"repository", "url", "name", "name"}, err.ErrorArguments())
	})
}

func (s *ErrorsSuite) TestRepository() {
	s.Run("NotFoundRepositoryError", func() {
		err := &NotFoundRepositoryError{URL: "url"}

		s.EqualError(err, "repository not found")
		s.Equal([]any{"url", "url"}, err.ErrorArguments())
	})
	s.Run("UnsupportedRepositoryError", func() {
		err := &UnsupportedRepositoryError{URL: "url"}

		s.EqualError(err, "unsupported repository url")
		s.Equal([]any{"url", "url"}, err.ErrorArguments())
	})
	s.Run("EmptyRepositoryError", func() {
		repositoryMock := &RepositoryMock{}
		repositoryMock.
			On("URL").Return("url")

		err := &EmptyRepositoryError{Repository: repositoryMock}

		s.EqualError(err, "empty repository")
		s.Equal([]any{"url", "url"}, err.ErrorArguments())
	})
}
