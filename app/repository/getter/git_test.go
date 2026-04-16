package getter_test

import (
	"io"
	"net/http"
	"net/http/httptest"
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

type GitSuite struct {
	suite.Suite

	server *httptest.Server
}

func TestGitSuite(t *testing.T) {
	suite.Run(t, new(GitSuite))
}

func (s *GitSuite) SetupSuite() {
	// Server
	mux := http.NewServeMux()
	s.server = httptest.NewServer(mux)

	mux.Handle("/repository.git/", http.StripPrefix("/repository.git/",
		http.FileServer(http.Dir(filepath.FromSlash("testdata/GitSuite/repository.git"))),
	))
}

func (s *GitSuite) TearDownSuite() {
	s.server.Close()
}

func (s *GitSuite) TestLoaderHandler() {
	cacheDir := filepath.FromSlash("testdata/cache")
	cache := caching.NewCache(cacheDir)

	s.Run("Http", func() {
		_ = os.RemoveAll(cacheDir)

		url := s.server.URL + "/repository.git"

		chainMock := &repository.LoaderHandlerChainMock{}

		handler := getter.NewGitLoaderHandler(log.New(io.Discard), cache)
		repository, err := handler.Handle(&repository.LoaderQuery{URL: url}, chainMock)

		s.Require().NoError(err)
		s.NotNil(repository)
		chainMock.AssertExpectations(s.T())

		s.DirExists(repository.Dir())
		heredoc.EqualFile(s.T(), `
			Hello World!
		`, filepath.Join(repository.Dir(), "README"))
	})

	s.Run("HttpForced", func() {
		_ = os.RemoveAll(cacheDir)

		url := "git::" + s.server.URL + "/repository.git"

		chainMock := &repository.LoaderHandlerChainMock{}

		handler := getter.NewGitLoaderHandler(log.New(io.Discard), cache)
		repository, err := handler.Handle(&repository.LoaderQuery{URL: url}, chainMock)

		s.Require().NoError(err)
		s.NotNil(repository)
		chainMock.AssertExpectations(s.T())

		s.DirExists(repository.Dir())
		heredoc.EqualFile(s.T(), `
			Hello World!
		`, filepath.Join(repository.Dir(), "README"))
	})
}
