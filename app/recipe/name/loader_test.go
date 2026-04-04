package name_test

import (
	"log/slog"
	"testing"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/recipe"
	"github.com/manala/manala/app/recipe/name"

	"github.com/stretchr/testify/suite"
)

type LoaderSuite struct{ suite.Suite }

func TestLoaderSuite(t *testing.T) {
	suite.Run(t, new(LoaderSuite))
}

func (s *LoaderSuite) TestProcessorHandler() {
	processor := name.NewProcessor(slog.New(slog.DiscardHandler))
	processor.Add("name", 10)

	handler := name.NewProcessorLoaderHandler(slog.New(slog.DiscardHandler), processor)

	repositoryMock := &app.RepositoryMock{}
	recipeMock := &app.RecipeMock{}

	chainMock := &recipe.LoaderHandlerChainMock{}
	chainMock.
		On("Next", &recipe.LoaderQuery{Repository: repositoryMock, Name: "name"}).Return(recipeMock, nil)

	recipe, err := handler.Handle(&recipe.LoaderQuery{Repository: repositoryMock, Name: ""}, chainMock)

	s.Require().NoError(err)
	s.Equal(recipeMock, recipe)
	chainMock.AssertExpectations(s.T())
}
