package getter

import (
	"github.com/stretchr/testify/suite"
	"manala/app/repository"
	"manala/internal/cache"
	"manala/internal/log"
	"manala/internal/testing/heredoc"
	"os"
	"path/filepath"
	"testing"
)

type GitSuite struct{ suite.Suite }

func TestGitSuite(t *testing.T) {
	suite.Run(t, new(GitSuite))
}

func (s *GitSuite) TestLoaderHandler() {
	cacheDir := filepath.FromSlash("testdata/cache")
	cache := cache.New(cacheDir)

	s.Run("Http", func() {
		_ = os.RemoveAll(cacheDir)

		url := "https://github.com/octocat/Hello-World.git"

		chainMock := &repository.LoaderHandlerChainMock{}

		handler := NewGitLoaderHandler(log.Discard, cache)
		repository, err := handler.Handle(&repository.LoaderQuery{Url: url}, chainMock)

		s.NotNil(repository)
		s.NoError(err)
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

		handler := NewGitLoaderHandler(log.Discard, cache)
		repository, err := handler.Handle(&repository.LoaderQuery{Url: url}, chainMock)

		s.NotNil(repository)
		s.NoError(err)
		chainMock.AssertExpectations(s.T())

		s.DirExists(filepath.Join(cacheDir, "repositories", "fb19e43b4b91cd3024866d556a7d04f40cc66e18c5f7098b9a4af8d0"))
		heredoc.EqualFile(s.T(), `
			Hello World!
		`, filepath.Join(cacheDir, "repositories", "fb19e43b4b91cd3024866d556a7d04f40cc66e18c5f7098b9a4af8d0", "README"))
	})
}
