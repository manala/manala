package project

import (
	"manala/app"
	"manala/internal/log"
	"manala/internal/serrors"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type LoaderSuite struct{ suite.Suite }

func TestLoaderSuite(t *testing.T) {
	suite.Run(t, new(LoaderSuite))
}

func (s *LoaderSuite) TestLoadErrors() {
	loader := NewLoader(log.Discard)

	s.Run("NotFound", func() {
		project, err := loader.Load("dir")

		s.Nil(project)
		serrors.Equal(s.T(), &serrors.Assertion{
			Type:    &app.NotFoundProjectError{},
			Message: "project not found",
			Arguments: []any{
				"dir", "dir",
			},
		}, err)
	})
}

func (s *LoaderSuite) TestLoad() {
	projectMock := &app.ProjectMock{}

	handlerMock := &LoaderHandlerMock{}
	handlerMock.
		On("Handle", &LoaderQuery{Dir: "dir"}, mock.Anything).Return(projectMock, nil)

	loader := NewLoader(log.Discard, WithLoaderHandlers(handlerMock))

	project, err := loader.Load("dir")

	s.Require().NoError(err)
	s.Equal(projectMock, project)
	handlerMock.AssertExpectations(s.T())
}

func (s *LoaderSuite) TestLoadRecursiveErrors() {
	s.Run("NotFound", func() {
		handlerMock := &LoaderHandlerMock{}

		loader := NewLoader(log.Discard, WithLoaderHandlers(handlerMock))

		err := loader.LoadRecursive("dir", func(project app.Project) error {
			return nil
		})

		serrors.Equal(s.T(), &serrors.Assertion{
			Message: "dir not found",
			Arguments: []any{
				"dir", "dir",
			},
		}, err)
		handlerMock.AssertExpectations(s.T())
	})
}

func (s *LoaderSuite) TestLoadRecursive() {
	projectsDir := filepath.FromSlash("testdata/LoaderSuite/TestLoadRecursive/projects")
	projectMock := &app.ProjectMock{}

	handlerMock := &LoaderHandlerMock{}
	handlerMock.
		On("Handle", &LoaderQuery{Dir: projectsDir}, mock.Anything).Return(projectMock, nil).
		On("Handle", &LoaderQuery{Dir: filepath.Join(projectsDir, "bar")}, mock.Anything).Return(projectMock, nil).
		On("Handle", &LoaderQuery{Dir: filepath.Join(projectsDir, "bar", "baz")}, mock.Anything).Return(projectMock, nil).
		On("Handle", &LoaderQuery{Dir: filepath.Join(projectsDir, "foo")}, mock.Anything).Return(projectMock, nil)

	loader := NewLoader(log.Discard, WithLoaderHandlers(handlerMock))

	err := loader.LoadRecursive(projectsDir, func(project app.Project) error {
		s.Equal(projectMock, project)

		return nil
	})

	s.Require().NoError(err)
	handlerMock.AssertExpectations(s.T())
}
