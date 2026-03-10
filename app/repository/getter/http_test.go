package getter_test

import (
	"github.com/manala/manala/app/repository"
	"github.com/manala/manala/app/repository/getter"
	"github.com/manala/manala/internal/caching"
	"github.com/manala/manala/internal/log"
	"github.com/manala/manala/internal/testing/heredoc"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type HTTPSuite struct{ suite.Suite }

func TestHttpSuite(t *testing.T) {
	suite.Run(t, new(HTTPSuite))
}

func (s *HTTPSuite) TestLoaderHandler() {
	cacheDir := filepath.FromSlash("testdata/cache")
	cache := caching.NewCache(cacheDir)

	s.Run("Zip", func() {
		_ = os.RemoveAll(cacheDir)

		url := "https://github.com/octocat/Hello-World/archive/refs/heads/master.zip"

		chainMock := &repository.LoaderHandlerChainMock{}

		handler := getter.NewHTTPLoaderHandler(log.Discard, cache)
		repository, err := handler.Handle(&repository.LoaderQuery{URL: url}, chainMock)

		s.Require().NoError(err)
		s.NotNil(repository)
		chainMock.AssertExpectations(s.T())

		s.DirExists(filepath.Join(cacheDir, "repositories", "5b3a694d5b3f61795c95c7540461f94ffb950b81deeb3cea08e20e97"))
		heredoc.EqualFile(s.T(), `
			Hello World!
		`, filepath.Join(cacheDir, "repositories", "5b3a694d5b3f61795c95c7540461f94ffb950b81deeb3cea08e20e97", "Hello-World-master", "README"))
	})

	s.Run("ZipSubdirectory", func() {
		_ = os.RemoveAll(cacheDir)

		url := "https://github.com/octocat/Hello-World/archive/refs/heads/master.zip//Hello-World-master"

		chainMock := &repository.LoaderHandlerChainMock{}

		handler := getter.NewHTTPLoaderHandler(log.Discard, cache)
		repository, err := handler.Handle(&repository.LoaderQuery{URL: url}, chainMock)

		s.Require().NoError(err)
		s.NotNil(repository)
		chainMock.AssertExpectations(s.T())

		s.DirExists(filepath.Join(cacheDir, "repositories", "e8b20dc6a839cb37051c82b0acee5e6b63fd96894a9ce12dbf548ba8"))
		heredoc.EqualFile(s.T(), `
			Hello World!
		`, filepath.Join(cacheDir, "repositories", "e8b20dc6a839cb37051c82b0acee5e6b63fd96894a9ce12dbf548ba8", "README"))
	})
}
