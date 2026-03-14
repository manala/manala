package getter_test

import (
	"fmt"
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

type S3Suite struct {
	suite.Suite

	server *httptest.Server
}

func TestS3Suite(t *testing.T) {
	suite.Run(t, new(S3Suite))
}

func (s *S3Suite) SetupSuite() {
	// Server
	mux := http.NewServeMux()
	s.server = httptest.NewServer(mux)

	// ListObjects — GET /bucket?prefix=repository
	mux.HandleFunc("GET /bucket", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		http.ServeFile(w, r, filepath.FromSlash("testdata/S3Suite/bucket/prefix.xml"))
	})
	// GetObject — GET /bucket/repository/…
	mux.HandleFunc("GET /bucket/repository/{object...}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("testdata", "S3Suite", "bucket", "repository", r.PathValue("object")))
	})
}

func (s *S3Suite) TearDownSuite() {
	s.server.Close()
}

func (s *S3Suite) TestLoaderHandler() {
	cacheDir := filepath.FromSlash("testdata/cache")
	cache := caching.NewCache(cacheDir)

	_ = os.RemoveAll(cacheDir)

	url := fmt.Sprintf("s3::%s/bucket/repository?aws_access_key_id=%s&aws_access_key_secret=%s",
		s.server.URL,
		"access_key_id",
		"access_key_secret",
	)

	chainMock := &repository.LoaderHandlerChainMock{}

	handler := getter.NewS3LoaderHandler(log.Discard, cache)
	repository, err := handler.Handle(&repository.LoaderQuery{URL: url}, chainMock)

	s.Require().NoError(err)
	s.NotNil(repository)
	chainMock.AssertExpectations(s.T())

	s.DirExists(repository.Dir())
	heredoc.EqualFile(s.T(), `
		foo
	`, filepath.Join(repository.Dir(), "file"))
}
