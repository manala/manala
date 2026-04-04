package url_test

import (
	"log/slog"
	"testing"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/repository"
	"github.com/manala/manala/app/repository/url"
	"github.com/manala/manala/internal/serrors"

	"github.com/stretchr/testify/suite"
)

type LoaderSuite struct{ suite.Suite }

func TestLoaderSuite(t *testing.T) {
	suite.Run(t, new(LoaderSuite))
}

func (s *LoaderSuite) TestProcessorHandlerErrors() {
	processor := url.NewProcessor(slog.New(slog.DiscardHandler))

	handler := url.NewProcessorLoaderHandler(slog.New(slog.DiscardHandler), processor)

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
	processor := url.NewProcessor(slog.New(slog.DiscardHandler))
	processor.Add("url", 10)

	handler := url.NewProcessorLoaderHandler(slog.New(slog.DiscardHandler), processor)

	repositoryMock := &app.RepositoryMock{}

	chainMock := &repository.LoaderHandlerChainMock{}
	chainMock.
		On("Next", &repository.LoaderQuery{URL: "url"}).Return(repositoryMock, nil)

	repository, err := handler.Handle(&repository.LoaderQuery{URL: ""}, chainMock)

	s.Require().NoError(err)
	s.Equal(repositoryMock, repository)
	chainMock.AssertExpectations(s.T())
}
