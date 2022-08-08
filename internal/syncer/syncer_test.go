package syncer

import (
	"bytes"
	"github.com/stretchr/testify/suite"
	internalLog "manala/internal/log"
	internalTemplate "manala/internal/template"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

type SyncerSuite struct {
	suite.Suite
	stderr           *bytes.Buffer
	syncer           *Syncer
	templateProvider internalTemplate.ProviderInterface
}

func TestSyncerSuite(t *testing.T) {
	suite.Run(t, new(SyncerSuite))
}

func (s *SyncerSuite) SetupTest() {
	s.stderr = &bytes.Buffer{}
	s.syncer = &Syncer{
		Log: internalLog.New(s.stderr),
	}
	s.templateProvider = &internalTemplate.Provider{}
}

var syncerTestPath = filepath.Join("testdata", "syncer")

func (s *SyncerSuite) TestSync() {
	sourcePath := filepath.Join(syncerTestPath, "sync", "source")
	destinationPath := filepath.Join(syncerTestPath, "sync", "destination")

	_ = os.RemoveAll(destinationPath)
	_ = os.Mkdir(destinationPath, 0755)
	_ = os.WriteFile(filepath.Join(destinationPath, "file_foo"), []byte("foo"), 0666)
	_ = os.WriteFile(filepath.Join(destinationPath, "file_bar"), []byte("bar"), 0666)
	_ = os.Mkdir(filepath.Join(destinationPath, "dir_empty"), 0755)
	_ = os.Mkdir(filepath.Join(destinationPath, "dir"), 0755)
	f, _ := os.Create(filepath.Join(destinationPath, "dir", "foo"))
	_ = f.Close()
	_ = os.WriteFile(filepath.Join(destinationPath, "dir", "foo"), []byte("bar"), 0666)
	_ = os.Mkdir(filepath.Join(destinationPath, "dir", "bar"), 0755)
	f, _ = os.Create(filepath.Join(destinationPath+"dir", "bar", "foo"))
	_ = f.Close()

	s.Run("Source not exists", func() {
		err := s.syncer.Sync(sourcePath, "baz", destinationPath, "baz", nil)
		s.Error(err, "no source file or directory")
	})

	s.Run("Destination file not exists", func() {
		err := s.syncer.Sync(sourcePath, "foo", destinationPath, "foo", nil)
		s.NoError(err)
		s.FileExists(filepath.Join(destinationPath, "foo"))
		content, _ := os.ReadFile(filepath.Join(destinationPath, "foo"))
		s.Equal("bar", string(content))
	})

	s.Run("Destination file exists and same", func() {
		err := s.syncer.Sync(sourcePath, "foo", destinationPath, "file_bar", nil)
		s.NoError(err)
		s.FileExists(filepath.Join(destinationPath, "file_bar"))
		content, _ := os.ReadFile(filepath.Join(destinationPath, "file_bar"))
		s.Equal("bar", string(content))
	})

	s.Run("Destination file exists and different", func() {
		err := s.syncer.Sync(sourcePath, "foo", destinationPath, "file_foo", nil)
		s.NoError(err)
		s.FileExists(filepath.Join(destinationPath, "file_foo"))
		content, _ := os.ReadFile(filepath.Join(destinationPath, "file_foo"))
		s.Equal("bar", string(content))
	})

	s.Run("Source file over destination directory empty", func() {
		err := s.syncer.Sync(sourcePath, "foo", destinationPath, "dir_empty", nil)
		s.NoError(err)
		s.FileExists(filepath.Join(destinationPath, "dir_empty"))
		content, _ := os.ReadFile(filepath.Join(destinationPath, "dir_empty"))
		s.Equal("bar", string(content))
	})

	s.Run("Source file over destination directory", func() {
		err := s.syncer.Sync(sourcePath, "foo", destinationPath, "dir", nil)
		s.NoError(err)
		s.FileExists(filepath.Join(destinationPath, "dir"))
		content, _ := os.ReadFile(filepath.Join(destinationPath, "dir"))
		s.Equal("bar", string(content))
	})

	s.Run("Destination directory not exists", func() {
		err := s.syncer.Sync(sourcePath, "bar", destinationPath, "bar", nil)
		s.NoError(err)
		s.FileExists(filepath.Join(destinationPath, "bar", "foo"))
		content, _ := os.ReadFile(filepath.Join(destinationPath, "bar", "foo"))
		s.Equal("baz", string(content))
	})

	s.Run("Destination directory exists", func() {
		err := s.syncer.Sync(sourcePath, "bar", destinationPath, "dir", nil)
		s.NoError(err)
		s.FileExists(filepath.Join(destinationPath, "dir", "foo"))
		content, _ := os.ReadFile(filepath.Join(destinationPath, "dir", "foo"))
		s.Equal("baz", string(content))
	})

	s.Run("Destination file directory not exists", func() {
		err := s.syncer.Sync(sourcePath, "foo", destinationPath, filepath.Join("baz", "foo"), nil)
		s.NoError(err)
		s.FileExists(filepath.Join(destinationPath, "baz", "foo"))
		content, _ := os.ReadFile(filepath.Join(destinationPath, "baz", "foo"))
		s.Equal("bar", string(content))
	})
}

func (s *SyncerSuite) TestSyncExecutable() {
	// Irrelevant on Windows
	// noinspection GoBoolExpressions
	if runtime.GOOS == "windows" {
		s.T().Skip()
	}

	sourcePath := filepath.Join(syncerTestPath, "sync_executable", "source")
	destinationPath := filepath.Join(syncerTestPath, "sync_executable", "destination")

	_ = os.RemoveAll(destinationPath)
	_ = os.Mkdir(destinationPath, 0755)
	_ = os.WriteFile(filepath.Join(destinationPath, "executable_true"), []byte(""), 0777)
	_ = os.WriteFile(filepath.Join(destinationPath, "executable_false"), []byte(""), 0666)

	s.Run("Source true", func() {
		err := s.syncer.Sync(sourcePath, "executable_true", destinationPath, "executable", nil)
		s.NoError(err)
		stat, _ := os.Stat(filepath.Join(destinationPath, "executable"))
		s.Equal(true, (stat.Mode()&0100) != 0)
	})

	s.Run("Source false", func() {
		err := s.syncer.Sync(sourcePath, "executable_false", destinationPath, "executable", nil)
		s.NoError(err)
		stat, _ := os.Stat(filepath.Join(destinationPath, "executable"))
		s.Equal(false, (stat.Mode()&0100) != 0)
	})

	s.Run("Source false destination false", func() {
		err := s.syncer.Sync(sourcePath, "executable_false", destinationPath, "executable_false", nil)
		s.NoError(err)
		stat, _ := os.Stat(filepath.Join(destinationPath, "executable_false"))
		s.Equal(false, (stat.Mode()&0100) != 0)
	})

	s.Run("Source true destination false", func() {
		err := s.syncer.Sync(sourcePath, "executable_true", destinationPath, "executable_false", nil)
		s.NoError(err)
		stat, _ := os.Stat(filepath.Join(destinationPath, "executable_false"))
		s.Equal(true, (stat.Mode()&0100) != 0)
	})

	s.Run("Source false destination true", func() {
		err := s.syncer.Sync(sourcePath, "executable_false", destinationPath, "executable_true", nil)
		s.NoError(err)
		stat, _ := os.Stat(filepath.Join(destinationPath, "executable_true"))
		s.Equal(false, (stat.Mode()&0100) != 0)
	})

	s.Run("Source true destination true", func() {
		err := s.syncer.Sync(sourcePath, "executable_true", destinationPath, "executable_true", nil)
		s.NoError(err)
		stat, _ := os.Stat(filepath.Join(destinationPath, "executable_true"))
		s.Equal(true, (stat.Mode()&0100) != 0)
	})
}

func (s *SyncerSuite) TestSyncTemplate() {
	sourcePath := filepath.Join(syncerTestPath, "sync_template", "source")
	destinationPath := filepath.Join(syncerTestPath, "sync_template", "destination")

	_ = os.RemoveAll(destinationPath)
	_ = os.Mkdir(destinationPath, 0755)
	_ = os.WriteFile(filepath.Join(destinationPath, "file_foo"), []byte("foo"), 0666)
	_ = os.WriteFile(filepath.Join(destinationPath, "file_bar"), []byte("bar"), 0666)

	s.Run("Source not exists", func() {
		err := s.syncer.Sync(sourcePath, "baz.tmpl", destinationPath, "baz", s.templateProvider)
		s.Error(err, "no source file or directory")
	})

	s.Run("Destination file not exists", func() {
		err := s.syncer.Sync(sourcePath, "foo.tmpl", destinationPath, "foo", s.templateProvider)
		s.NoError(err)
		s.FileExists(filepath.Join(destinationPath, "foo"))
		content, _ := os.ReadFile(filepath.Join(destinationPath, "foo"))
		s.Equal("bar", string(content))
	})

	s.Run("Destination file exists and same", func() {
		err := s.syncer.Sync(sourcePath, "foo.tmpl", destinationPath, "file_bar", s.templateProvider)
		s.NoError(err)
		s.FileExists(filepath.Join(destinationPath, "file_bar"))
		content, _ := os.ReadFile(filepath.Join(destinationPath, "file_bar"))
		s.Equal("bar", string(content))
	})

	s.Run("Destination file exists and different", func() {
		err := s.syncer.Sync(sourcePath, "foo.tmpl", destinationPath, "file_foo", s.templateProvider)
		s.NoError(err)
		s.FileExists(filepath.Join(destinationPath, "file_foo"))
		content, _ := os.ReadFile(filepath.Join(destinationPath, "file_foo"))
		s.Equal("bar", string(content))
	})

	s.Run("Invalid", func() {
		err := s.syncer.Sync(sourcePath, "invalid.tmpl", destinationPath, "invalid", s.templateProvider)
		s.Error(err, "template execution error")
	})
}
