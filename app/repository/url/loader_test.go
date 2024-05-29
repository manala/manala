package url

import (
	"manala/app"
	"manala/app/repository"
	"manala/internal/log"
	"manala/internal/serrors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type LoaderSuite struct{ suite.Suite }

func TestLoaderSuite(t *testing.T) {
	suite.Run(t, new(LoaderSuite))
}

func (s *LoaderSuite) TestProcessorHandlerErrors() {
	processor := NewProcessor(log.Discard)

	handler := NewProcessorLoaderHandler(log.Discard, processor)

	chainMock := &repository.LoaderHandlerChainMock{}

	repository, err := handler.Handle(&repository.LoaderQuery{Url: "foo?bar;baz"}, chainMock)

	s.Nil(repository)
	serrors.Equal(s.T(), &serrors.Assertion{
		Message: "unable to process repository query",
		Arguments: []any{
			"query", "bar;baz",
		},
		Errors: []*serrors.Assertion{
			{
				Message: "invalid semicolon separator in query",
			},
		},
	}, err)
	chainMock.AssertExpectations(s.T())
}

func (s *LoaderSuite) TestProcessorHandler() {
	processor := NewProcessor(log.Discard)
	processor.Add("url", 10)

	handler := NewProcessorLoaderHandler(log.Discard, processor)

	repositoryMock := &app.RepositoryMock{}

	chainMock := &repository.LoaderHandlerChainMock{}
	chainMock.
		On("Next", &repository.LoaderQuery{Url: "url"}).Return(repositoryMock, nil)

	repository, err := handler.Handle(&repository.LoaderQuery{Url: ""}, chainMock)

	s.Equal(repositoryMock, repository)
	s.NoError(err)
	chainMock.AssertExpectations(s.T())
}
