package getter

import (
	"github.com/stretchr/testify/suite"
	"manala/app/repository"
	"manala/internal/log"
	"path/filepath"
	"testing"
)

type FileSuite struct{ suite.Suite }

func TestFileSuite(t *testing.T) {
	suite.Run(t, new(FileSuite))
}

func (s *FileSuite) TestLoaderHandler() {
	s.Run("Relative", func() {
		url := filepath.FromSlash("testdata/FileSuite/TestLoaderHandler/Relative/repository")

		chainMock := &repository.LoaderHandlerChainMock{}

		handler := NewFileLoaderHandler(log.Discard)
		repository, err := handler.Handle(&repository.LoaderQuery{Url: url}, chainMock)

		s.Equal(url, repository.Url())
		s.NoError(err)
		chainMock.AssertExpectations(s.T())
	})

	s.Run("Absolute", func() {
		url := filepath.FromSlash("testdata/FileSuite/TestLoaderHandler/Absolute/repository")
		url, _ = filepath.Abs(url)

		chainMock := &repository.LoaderHandlerChainMock{}

		handler := NewFileLoaderHandler(log.Discard)
		repository, err := handler.Handle(&repository.LoaderQuery{Url: url}, chainMock)

		s.Equal(url, repository.Url())
		s.NoError(err)
		chainMock.AssertExpectations(s.T())
	})
}
