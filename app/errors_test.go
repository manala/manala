package app_test

import (
	"testing"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/testing/errors"
	"github.com/manala/manala/app/testing/mocks"
	"github.com/manala/manala/internal/testing/expect"

	"github.com/stretchr/testify/suite"
)

type ErrorsSuite struct{ suite.Suite }

func TestErrorsSuite(t *testing.T) {
	suite.Run(t, new(ErrorsSuite))
}

func (s *ErrorsSuite) TestProject() {
	s.Run("AlreadyExistingProjectError", func() {
		err := &app.AlreadyExistingProjectError{Dir: "dir"}

		expect.Error(s.T(), errors.Expectation{
			Type:  &app.AlreadyExistingProjectError{},
			Attrs: [][2]any{{"dir", "dir"}},
		}, err)
	})

	s.Run("NotFoundProjectError", func() {
		err := &app.NotFoundProjectError{Dir: "dir"}

		expect.Error(s.T(), errors.Expectation{
			Type:  &app.NotFoundProjectError{},
			Attrs: [][2]any{{"dir", "dir"}},
		}, err)
	})
}

func (s *ErrorsSuite) TestRecipe() {
	s.Run("NotFoundRecipeError", func() {
		repositoryMock := &mocks.RepositoryMock{}
		repositoryMock.
			On("URL").Return("url")

		err := &app.NotFoundRecipeError{Repository: repositoryMock, Name: "name"}

		expect.Error(s.T(), errors.Expectation{
			Type:  &app.NotFoundRecipeError{},
			Attrs: [][2]any{{"repository", "url"}, {"name", "name"}},
		}, err)
	})
}

func (s *ErrorsSuite) TestRepository() {
	s.Run("NotFoundRepositoryError", func() {
		err := &app.NotFoundRepositoryError{URL: "url"}

		expect.Error(s.T(), errors.Expectation{
			Type:  &app.NotFoundRepositoryError{},
			Attrs: [][2]any{{"url", "url"}},
		}, err)
	})

	s.Run("EmptyRepositoryError", func() {
		repositoryMock := &mocks.RepositoryMock{}
		repositoryMock.
			On("URL").Return("url")

		err := &app.EmptyRepositoryError{Repository: repositoryMock}

		expect.Error(s.T(), errors.Expectation{
			Type:  &app.EmptyRepositoryError{},
			Attrs: [][2]any{{"url", "url"}},
		}, err)
	})
}
