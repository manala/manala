package getter_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/manala/manala/app/repository"
	"github.com/manala/manala/app/repository/getter"
	"github.com/manala/manala/internal/caching"
	"github.com/manala/manala/internal/log"
	"github.com/manala/manala/internal/testing/heredoc"

	"github.com/stretchr/testify/suite"
)

type GitSuite struct{ suite.Suite }

func TestGitSuite(t *testing.T) {
	suite.Run(t, new(GitSuite))
}

func (s *GitSuite) TestLoaderHandler() {
	cacheDir := filepath.FromSlash("testdata/cache")
	cache := caching.NewCache(cacheDir)

	s.Run("Http", func() {
		_ = os.RemoveAll(cacheDir)

		url := "https://github.com/octocat/Hello-World.git"

		chainMock := &repository.LoaderHandlerChainMock{}

		handler := getter.NewGitLoaderHandler(log.Discard, cache)
		repository, err := handler.Handle(&repository.LoaderQuery{URL: url}, chainMock)

		s.Require().NoError(err)
		s.NotNil(repository)
		chainMock.AssertExpectations(s.T())

		s.DirExists(filepath.Join(cacheDir, "repositories", "0abe222eb5c9f101220b454d055bd5cbbb419722eddec6ad296f343f"))
		heredoc.EqualFile(s.T(), `
			Hello World!
		`, filepath.Join(cacheDir, "repositories", "0abe222eb5c9f101220b454d055bd5cbbb419722eddec6ad296f343f", "README"))
	})

	s.Run("HttpForced", func() {
		_ = os.RemoveAll(cacheDir)

		url := "git::https://github.com/octocat/Hello-World.git"

		chainMock := &repository.LoaderHandlerChainMock{}

		handler := getter.NewGitLoaderHandler(log.Discard, cache)
		repository, err := handler.Handle(&repository.LoaderQuery{URL: url}, chainMock)

		s.Require().NoError(err)
		s.NotNil(repository)
		chainMock.AssertExpectations(s.T())

		s.DirExists(filepath.Join(cacheDir, "repositories", "fb19e43b4b91cd3024866d556a7d04f40cc66e18c5f7098b9a4af8d0"))
		heredoc.EqualFile(s.T(), `
			Hello World!
		`, filepath.Join(cacheDir, "repositories", "fb19e43b4b91cd3024866d556a7d04f40cc66e18c5f7098b9a4af8d0", "README"))
	})
}
