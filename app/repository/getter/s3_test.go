package getter

import (
	"fmt"
	"github.com/johannesboyne/gofakes3"
	"github.com/johannesboyne/gofakes3/backend/s3mem"
	"github.com/stretchr/testify/suite"
	"manala/app/repository"
	"manala/internal/cache"
	"manala/internal/log"
	"manala/internal/testing/heredoc"
	"net"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type S3Suite struct{ suite.Suite }

func TestS3Suite(t *testing.T) {
	suite.Run(t, new(S3Suite))
}

func (s *S3Suite) TestLoaderHandler() {
	cacheDir := filepath.FromSlash("testdata/cache")
	cache := cache.New(cacheDir)

	_ = os.RemoveAll(cacheDir)

	// Fake S3
	s3Server, s3Backend := s.fakeS3("127.0.0.1:1234")
	defer s3Server.Close()

	s3File := "foo\n"
	_ = s3Backend.CreateBucket("bucket")
	_, _ = s3Backend.PutObject("bucket", "repository/file", map[string]string{}, strings.NewReader(s3File), int64(len(s3File)))

	url := fmt.Sprintf("s3::%s/bucket/repository?aws_access_key_id=%s&aws_access_key_secret=%s",
		s3Server.URL,
		"access_key_id",
		"access_key_secret",
	)

	chainMock := &repository.LoaderHandlerChainMock{}

	handler := NewS3LoaderHandler(log.Discard, cache)
	repository, err := handler.Handle(&repository.LoaderQuery{Url: url}, chainMock)

	s.NotNil(repository)
	s.NoError(err)
	chainMock.AssertExpectations(s.T())

	s.DirExists(filepath.Join(cacheDir, "repositories", "0b05624a43aa6dfd14fa0dd68105f49f20466339e403bdb1ad7ae55b"))
	heredoc.EqualFile(s.T(), `
		foo
	`, filepath.Join(cacheDir, "repositories", "0b05624a43aa6dfd14fa0dd68105f49f20466339e403bdb1ad7ae55b", "file"))
}

func (s *S3Suite) fakeS3(address string) (*httptest.Server, gofakes3.Backend) {
	// Create a listener with the desired address
	listener, _ := net.Listen("tcp", address)

	backend := s3mem.New()
	fake := gofakes3.New(backend)

	server := httptest.NewUnstartedServer(fake.Server())

	// Close the automatically created listener and replace with the one we created
	_ = server.Listener.Close()
	server.Listener = listener

	server.Start()

	return server, backend
}
