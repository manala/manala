package getter_test

import (
	"manala/app/repository"
	"manala/app/repository/getter"
	"manala/internal/caching"
	"manala/internal/log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type S3Suite struct{ suite.Suite }

func TestS3Suite(t *testing.T) {
	suite.Run(t, new(S3Suite))
}

func (s *S3Suite) TestLoaderHandler() {
	cacheDir := filepath.FromSlash("testdata/cache")
	cache := caching.NewCache(cacheDir)

	_ = os.RemoveAll(cacheDir)

	// See: https://github.com/hairyhenderson/gomplate/blob/29bd9ed54a89918e6225d382b4306e1f68c74f3f/internal/tests/integration/datasources_blob_test.go#L84
	url := "noaa-bathymetry-pds.s3.amazonaws.com/test"

	chainMock := &repository.LoaderHandlerChainMock{}

	handler := getter.NewS3LoaderHandler(log.Discard, cache)
	repository, err := handler.Handle(&repository.LoaderQuery{URL: url}, chainMock)

	s.Require().NoError(err)
	s.NotNil(repository)
	chainMock.AssertExpectations(s.T())

	s.DirExists(filepath.Join(cacheDir, "repositories", "aefad621d26200de9fb8b906a3167feb12d3be18c76eeed63280c97c"))
	s.DirExists(filepath.Join(cacheDir, "repositories", "aefad621d26200de9fb8b906a3167feb12d3be18c76eeed63280c97c", "2020"))
}
