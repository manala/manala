package url_test

import (
	"manala/app"
	"manala/app/repository"
	"manala/app/repository/url"
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
	processor := url.NewProcessor(log.Discard)

	handler := url.NewProcessorLoaderHandler(log.Discard, processor)

	chainMock := &repository.LoaderHandlerChainMock{}

	repository, err := handler.Handle(&repository.LoaderQuery{URL: "foo?bar;baz"}, chainMock)

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
	processor := url.NewProcessor(log.Discard)
	processor.Add("url", 10)

	handler := url.NewProcessorLoaderHandler(log.Discard, processor)

	repositoryMock := &app.RepositoryMock{}

	chainMock := &repository.LoaderHandlerChainMock{}
	chainMock.
		On("Next", &repository.LoaderQuery{URL: "url"}).Return(repositoryMock, nil)

	repository, err := handler.Handle(&repository.LoaderQuery{URL: ""}, chainMock)

	s.Require().NoError(err)
	s.Equal(repositoryMock, repository)
	chainMock.AssertExpectations(s.T())
}
