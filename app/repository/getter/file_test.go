package getter_test

import (
	"path/filepath"
	"testing"

	"manala/app/repository"
	"manala/app/repository/getter"
	"manala/internal/log"

	"github.com/stretchr/testify/suite"
)

type FileSuite struct{ suite.Suite }

func TestFileSuite(t *testing.T) {
	suite.Run(t, new(FileSuite))
}

func (s *FileSuite) TestLoaderHandler() {
	s.Run("Relative", func() {
		url := filepath.FromSlash("testdata/FileSuite/TestLoaderHandler/Relative/repository")

		chainMock := &repository.LoaderHandlerChainMock{}

		handler := getter.NewFileLoaderHandler(log.Discard)
		repository, err := handler.Handle(&repository.LoaderQuery{URL: url}, chainMock)

		s.Require().NoError(err)
		s.Equal(url, repository.URL())
		chainMock.AssertExpectations(s.T())
	})

	s.Run("Absolute", func() {
		url := filepath.FromSlash("testdata/FileSuite/TestLoaderHandler/Absolute/repository")
		url, _ = filepath.Abs(url)

		chainMock := &repository.LoaderHandlerChainMock{}

		handler := getter.NewFileLoaderHandler(log.Discard)
		repository, err := handler.Handle(&repository.LoaderQuery{URL: url}, chainMock)

		s.Require().NoError(err)
		s.Equal(url, repository.URL())
		chainMock.AssertExpectations(s.T())
	})
}
