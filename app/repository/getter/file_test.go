package getter_test

import (
	"log/slog"
	"path/filepath"
	"testing"

	"github.com/manala/manala/app/repository"
	"github.com/manala/manala/app/repository/getter"

	"github.com/stretchr/testify/suite"
)

type FileSuite struct{ suite.Suite }

func TestFileSuite(t *testing.T) {
	suite.Run(t, new(FileSuite))
}

func (s *FileSuite) TestLoaderHandler() {
	s.Run("Relative", func() {
		url := filepath.FromSlash("testdata/FileSuite/repository")

		chainMock := &repository.LoaderHandlerChainMock{}

		handler := getter.NewFileLoaderHandler(slog.New(slog.DiscardHandler))
		repository, err := handler.Handle(&repository.LoaderQuery{URL: url}, chainMock)

		s.Require().NoError(err)
		s.Equal(url, repository.URL())
		chainMock.AssertExpectations(s.T())
	})

	s.Run("Absolute", func() {
		url := filepath.FromSlash("testdata/FileSuite/repository")
		url, _ = filepath.Abs(url)

		chainMock := &repository.LoaderHandlerChainMock{}

		handler := getter.NewFileLoaderHandler(slog.New(slog.DiscardHandler))
		repository, err := handler.Handle(&repository.LoaderQuery{URL: url}, chainMock)

		s.Require().NoError(err)
		s.Equal(url, repository.URL())
		chainMock.AssertExpectations(s.T())
	})
}
