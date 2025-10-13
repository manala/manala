package app_test

import (
	"testing"

	"manala/app"

	"github.com/stretchr/testify/suite"
)

type ErrorsSuite struct{ suite.Suite }

func TestErrorsSuite(t *testing.T) {
	suite.Run(t, new(ErrorsSuite))
}

func (s *ErrorsSuite) TestProject() {
	s.Run("AlreadyExistingProjectError", func() {
		err := &app.AlreadyExistingProjectError{Dir: "dir"}

		s.Require().EqualError(err, "already existing project")
		s.Equal([]any{"dir", "dir"}, err.ErrorArguments())
	})

	s.Run("NotFoundProjectError", func() {
		err := &app.NotFoundProjectError{Dir: "dir"}

		s.Require().EqualError(err, "project not found")
		s.Equal([]any{"dir", "dir"}, err.ErrorArguments())
	})
}

func (s *ErrorsSuite) TestRecipe() {
	s.Run("NotFoundRecipeError", func() {
		repositoryMock := &app.RepositoryMock{}
		repositoryMock.
			On("URL").Return("url")

		err := &app.NotFoundRecipeError{Repository: repositoryMock, Name: "name"}

		s.Require().EqualError(err, "recipe not found")
		s.Equal([]any{"repository", "url", "name", "name"}, err.ErrorArguments())
	})
}

func (s *ErrorsSuite) TestRepository() {
	s.Run("NotFoundRepositoryError", func() {
		err := &app.NotFoundRepositoryError{URL: "url"}

		s.Require().EqualError(err, "repository not found")
		s.Equal([]any{"url", "url"}, err.ErrorArguments())
	})
	s.Run("UnsupportedRepositoryError", func() {
		err := &app.UnsupportedRepositoryError{URL: "url"}

		s.Require().EqualError(err, "unsupported repository url")
		s.Equal([]any{"url", "url"}, err.ErrorArguments())
	})
	s.Run("EmptyRepositoryError", func() {
		repositoryMock := &app.RepositoryMock{}
		repositoryMock.
			On("URL").Return("url")

		err := &app.EmptyRepositoryError{Repository: repositoryMock}

		s.Require().EqualError(err, "empty repository")
		s.Equal([]any{"url", "url"}, err.ErrorArguments())
	})
}
