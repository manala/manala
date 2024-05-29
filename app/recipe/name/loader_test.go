package name

import (
	"manala/app"
	"manala/app/recipe"
	"manala/internal/log"
	"testing"

	"github.com/stretchr/testify/suite"
)

type LoaderSuite struct{ suite.Suite }

func TestLoaderSuite(t *testing.T) {
	suite.Run(t, new(LoaderSuite))
}

func (s *LoaderSuite) TestProcessorHandler() {
	processor := NewProcessor(log.Discard)
	processor.Add("name", 10)

	handler := NewProcessorLoaderHandler(log.Discard, processor)

	repositoryMock := &app.RepositoryMock{}
	recipeMock := &app.RecipeMock{}

	chainMock := &recipe.LoaderHandlerChainMock{}
	chainMock.
		On("Next", &recipe.LoaderQuery{Repository: repositoryMock, Name: "name"}).Return(recipeMock, nil)

	recipe, err := handler.Handle(&recipe.LoaderQuery{Repository: repositoryMock, Name: ""}, chainMock)

	s.Equal(recipeMock, recipe)
	s.NoError(err)
	chainMock.AssertExpectations(s.T())
}
