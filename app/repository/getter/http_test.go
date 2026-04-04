package getter_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/manala/manala/app/repository"
	"github.com/manala/manala/app/repository/getter"
	"github.com/manala/manala/internal/caching"
	"github.com/manala/manala/internal/testing/heredoc"

	"github.com/stretchr/testify/suite"
)

type HTTPSuite struct {
	suite.Suite

	server *httptest.Server
}

func TestHttpSuite(t *testing.T) {
	suite.Run(t, new(HTTPSuite))
}

func (s *HTTPSuite) SetupSuite() {
	// Server
	mux := http.NewServeMux()
	s.server = httptest.NewServer(mux)

	mux.HandleFunc("GET /{file...}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("testdata", "HTTPSuite", r.PathValue("file")))
	})
}

func (s *HTTPSuite) TearDownSuite() {
	s.server.Close()
}

func (s *HTTPSuite) TestLoaderHandler() {
	cacheDir := filepath.FromSlash("testdata/cache")
	cache := caching.NewCache(cacheDir)

	s.Run("Zip", func() {
		_ = os.RemoveAll(cacheDir)

		url := s.server.URL + "/archive.zip"

		chainMock := &repository.LoaderHandlerChainMock{}

		handler := getter.NewHTTPLoaderHandler(slog.New(slog.DiscardHandler), cache)
		repository, err := handler.Handle(&repository.LoaderQuery{URL: url}, chainMock)

		s.Require().NoError(err)
		s.NotNil(repository)
		chainMock.AssertExpectations(s.T())

		s.DirExists(repository.Dir())
		heredoc.EqualFile(s.T(), `
			Hello World!
		`, filepath.Join(repository.Dir(), "Hello-World-master", "README"))
	})

	s.Run("ZipSubdirectory", func() {
		_ = os.RemoveAll(cacheDir)

		url := s.server.URL + "/archive.zip//Hello-World-master"

		chainMock := &repository.LoaderHandlerChainMock{}

		handler := getter.NewHTTPLoaderHandler(slog.New(slog.DiscardHandler), cache)
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
