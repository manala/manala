package url_test

import (
	"testing"

	"github.com/manala/manala/app/repository"
	"github.com/manala/manala/app/repository/url"
	"github.com/manala/manala/app/testing/mocks"
	"github.com/manala/manala/internal/log"
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/testing/expect"

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
	expect.Error(s.T(), serrors.Expectation{
		Message: "unable to process repository query",
		Attrs: [][2]any{
			{"query", "bar;baz"},
		},
		Errors: []expect.ErrorExpectation{
			expect.ErrorMessageExpectation("invalid semicolon separator in query"),
		},
	}, err)
	chainMock.AssertExpectations(s.T())
}

func (s *LoaderSuite) TestProcessorHandler() {
	processor := url.NewProcessor(log.Discard)
	processor.Add("url", 10)

	handler := url.NewProcessorLoaderHandler(log.Discard, processor)

	repositoryMock := &mocks.RepositoryMock{}

	chainMock := &repository.LoaderHandlerChainMock{}
	chainMock.
		On("Next", &repository.LoaderQuery{URL: "url"}).Return(repositoryMock, nil)

	repository, err := handler.Handle(&repository.LoaderQuery{URL: ""}, chainMock)

	s.Require().NoError(err)
	s.Equal(repositoryMock, repository)
	chainMock.AssertExpectations(s.T())
}
