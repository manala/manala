package sync_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/manala/manala/internal/log"
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/sync"
	"github.com/manala/manala/internal/template"

	"github.com/stretchr/testify/suite"
)

type SyncerSuite struct {
	suite.Suite

	syncer           *sync.Syncer
	templateProvider template.ProviderInterface
}

func TestSyncerSuite(t *testing.T) {
	suite.Run(t, new(SyncerSuite))
}

func (s *SyncerSuite) SetupTest() {
	s.syncer = sync.NewSyncer(log.Discard)
	s.templateProvider = &template.Provider{}
}

func (s *SyncerSuite) TestSync() {
	sourcePath := filepath.FromSlash("testdata/SyncerSuite/TestSync/source")
	destinationPath := filepath.FromSlash("testdata/SyncerSuite/TestSync/destination")

	_ = os.RemoveAll(destinationPath)
	_ = os.Mkdir(destinationPath, 0o755)
	_ = os.WriteFile(filepath.Join(destinationPath, "file_foo"), []byte("foo"), 0o666)
	_ = os.WriteFile(filepath.Join(destinationPath, "file_bar"), []byte("bar"), 0o666)
	_ = os.Mkdir(filepath.Join(destinationPath, "dir_empty"), 0o755)
	_ = os.Mkdir(filepath.Join(destinationPath, "dir"), 0o755)
	f, _ := os.Create(filepath.Join(destinationPath, "dir", "foo"))
	_ = f.Close()
	_ = os.WriteFile(filepath.Join(destinationPath, "dir", "foo"), []byte("bar"), 0o666)
	_ = os.Mkdir(filepath.Join(destinationPath, "dir", "bar"), 0o755)
	f, _ = os.Create(filepath.Join(destinationPath+"dir", "bar", "foo"))
	_ = f.Close()

	s.Run("SourceNotExists", func() {
		err := s.syncer.Sync(sourcePath, "baz", destinationPath, "baz", nil)

		serrors.Equal(s.T(), &serrors.Assertion{
			Message: "no source file or directory",
			Arguments: []any{
				"path", filepath.Join(sourcePath, "baz"),
			},
		}, err)
	})

	s.Run("DestinationFileNotExists", func() {
		err := s.syncer.Sync(sourcePath, "foo", destinationPath, "foo", nil)
		s.Require().NoError(err)
		s.FileExists(filepath.Join(destinationPath, "foo"))

		content, _ := os.ReadFile(filepath.Join(destinationPath, "foo"))
		s.Equal("bar", string(content))
	})

	s.Run("DestinationFileExistsAndSame", func() {
		err := s.syncer.Sync(sourcePath, "foo", destinationPath, "file_bar", nil)
		s.Require().NoError(err)
		s.FileExists(filepath.Join(destinationPath, "file_bar"))

		content, _ := os.ReadFile(filepath.Join(destinationPath, "file_bar"))
		s.Equal("bar", string(content))
	})

	s.Run("DestinationFileExistsAndDifferent", func() {
		err := s.syncer.Sync(sourcePath, "foo", destinationPath, "file_foo", nil)
		s.Require().NoError(err)
		s.FileExists(filepath.Join(destinationPath, "file_foo"))

		content, _ := os.ReadFile(filepath.Join(destinationPath, "file_foo"))
		s.Equal("bar", string(content))
	})

	s.Run("SourceFileOverDestinationDirectoryEmpty", func() {
		err := s.syncer.Sync(sourcePath, "foo", destinationPath, "dir_empty", nil)
		s.Require().NoError(err)
		s.FileExists(filepath.Join(destinationPath, "dir_empty"))

		content, _ := os.ReadFile(filepath.Join(destinationPath, "dir_empty"))
		s.Equal("bar", string(content))
	})

	s.Run("SourceFileOverDestinationDirectory", func() {
		err := s.syncer.Sync(sourcePath, "foo", destinationPath, "dir", nil)
		s.Require().NoError(err)
		s.FileExists(filepath.Join(destinationPath, "dir"))

		content, _ := os.ReadFile(filepath.Join(destinationPath, "dir"))
		s.Equal("bar", string(content))
	})

	s.Run("DestinationDirectoryNotExists", func() {
		err := s.syncer.Sync(sourcePath, "bar", destinationPath, "bar", nil)
		s.Require().NoError(err)
		s.FileExists(filepath.Join(destinationPath, "bar", "foo"))

		content, _ := os.ReadFile(filepath.Join(destinationPath, "bar", "foo"))
		s.Equal("baz", string(content))
	})

	s.Run("DestinationDirectoryExists", func() {
		err := s.syncer.Sync(sourcePath, "bar", destinationPath, "dir", nil)
		s.Require().NoError(err)
		s.FileExists(filepath.Join(destinationPath, "dir", "foo"))

		content, _ := os.ReadFile(filepath.Join(destinationPath, "dir", "foo"))
		s.Equal("baz", string(content))
	})

	s.Run("DestinationFileDirectoryNotexists", func() {
		err := s.syncer.Sync(sourcePath, "foo", destinationPath, filepath.Join("baz", "foo"), nil)
		s.Require().NoError(err)
		s.FileExists(filepath.Join(destinationPath, "baz", "foo"))

		content, _ := os.ReadFile(filepath.Join(destinationPath, "baz", "foo"))
		s.Equal("bar", string(content))
	})
}

func (s *SyncerSuite) TestSyncExecutable() {
	// Irrelevant on Windows
	//goland:noinspection GoBoolExpressions
	if runtime.GOOS == "windows" {
		s.T().Skip()
	}

	sourcePath := filepath.FromSlash("testdata/SyncerSuite/TestSyncExecutable/source")
	destinationPath := filepath.FromSlash("testdata/SyncerSuite/TestSyncExecutable/destination")

	_ = os.RemoveAll(destinationPath)
	_ = os.Mkdir(destinationPath, 0o755)
	_ = os.WriteFile(filepath.Join(destinationPath, "executable_true"), []byte(""), 0o777)
	_ = os.WriteFile(filepath.Join(destinationPath, "executable_false"), []byte(""), 0o666)

	s.Run("SourceTrue", func() {
		err := s.syncer.Sync(sourcePath, "executable_true", destinationPath, "executable", nil)
		s.Require().NoError(err)

		stat, _ := os.Stat(filepath.Join(destinationPath, "executable"))
		s.NotEqual(0, int(stat.Mode()&0o100))
	})

	s.Run("SourceFalse", func() {
		err := s.syncer.Sync(sourcePath, "executable_false", destinationPath, "executable", nil)
		s.Require().NoError(err)

		stat, _ := os.Stat(filepath.Join(destinationPath, "executable"))
		s.Equal(0, int(stat.Mode()&0o100))
	})

	s.Run("SourceFalseDestinationFalse", func() {
		err := s.syncer.Sync(sourcePath, "executable_false", destinationPath, "executable_false", nil)
		s.Require().NoError(err)

		stat, _ := os.Stat(filepath.Join(destinationPath, "executable_false"))
		s.Equal(0, int(stat.Mode()&0o100))
	})

	s.Run("SourceTrueDestinationFalse", func() {
		err := s.syncer.Sync(sourcePath, "executable_true", destinationPath, "executable_false", nil)
		s.Require().NoError(err)

		stat, _ := os.Stat(filepath.Join(destinationPath, "executable_false"))
		s.NotEqual(0, int(stat.Mode()&0o100))
	})

	s.Run("SourceFalseDestinationTrue", func() {
		err := s.syncer.Sync(sourcePath, "executable_false", destinationPath, "executable_true", nil)
		s.Require().NoError(err)

		stat, _ := os.Stat(filepath.Join(destinationPath, "executable_true"))
		s.Equal(0, int(stat.Mode()&0o100))
	})

	s.Run("SourceTrueDestinationTrue", func() {
		err := s.syncer.Sync(sourcePath, "executable_true", destinationPath, "executable_true", nil)
		s.Require().NoError(err)

		stat, _ := os.Stat(filepath.Join(destinationPath, "executable_true"))
		s.NotEqual(0, int(stat.Mode()&0o100))
	})
}

func (s *SyncerSuite) TestSyncTemplate() {
	sourcePath := filepath.FromSlash("testdata/SyncerSuite/TestSyncTemplate/source")
	destinationPath := filepath.FromSlash("testdata/SyncerSuite/TestSyncTemplate/destination")

	_ = os.RemoveAll(destinationPath)
	_ = os.Mkdir(destinationPath, 0o755)
	_ = os.WriteFile(filepath.Join(destinationPath, "file_foo"), []byte("foo"), 0o666)
	_ = os.WriteFile(filepath.Join(destinationPath, "file_bar"), []byte("bar"), 0o666)

	s.Run("SourceNotExists", func() {
		err := s.syncer.Sync(sourcePath, "baz.tmpl", destinationPath, "baz", s.templateProvider)

		serrors.Equal(s.T(), &serrors.Assertion{
			Message: "no source file or directory",
			Arguments: []any{
				"path", filepath.Join(sourcePath, "baz.tmpl"),
			},
		}, err)
	})

	s.Run("DestinationFileNotExists", func() {
		err := s.syncer.Sync(sourcePath, "foo.tmpl", destinationPath, "foo", s.templateProvider)
		s.Require().NoError(err)
		s.FileExists(filepath.Join(destinationPath, "foo"))

		content, _ := os.ReadFile(filepath.Join(destinationPath, "foo"))
		s.Equal("bar", string(content))
	})

	s.Run("DestinationFileExistsAndSame", func() {
		err := s.syncer.Sync(sourcePath, "foo.tmpl", destinationPath, "file_bar", s.templateProvider)
		s.Require().NoError(err)
		s.FileExists(filepath.Join(destinationPath, "file_bar"))

		content, _ := os.ReadFile(filepath.Join(destinationPath, "file_bar"))
		s.Equal("bar", string(content))
	})

	s.Run("DestinationFileExistsAndDifferent", func() {
		err := s.syncer.Sync(sourcePath, "foo.tmpl", destinationPath, "file_foo", s.templateProvider)
		s.Require().NoError(err)
		s.FileExists(filepath.Join(destinationPath, "file_foo"))

		content, _ := os.ReadFile(filepath.Join(destinationPath, "file_foo"))
		s.Equal("bar", string(content))
	})

	s.Run("Invalid", func() {
		err := s.syncer.Sync(sourcePath, "invalid.tmpl", destinationPath, "invalid", s.templateProvider)

		serrors.Equal(s.T(), &serrors.Assertion{
			Message: "template error",
			Errors: []*serrors.Assertion{
				{
					Message: "nil data; no entry for key \"foo\"",
					Arguments: []any{
						"context", ".foo",
						"template", "invalid.tmpl",
						"line", 1,
						"column", 3,
						"file", filepath.Join(sourcePath, "invalid.tmpl"),
					},
				},
			},
		}, err)
	})
}
