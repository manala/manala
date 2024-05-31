package syncer

import (
	"manala/internal/log"
	"manala/internal/serrors"
	"manala/internal/template"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	syncer           *Syncer
	templateProvider template.ProviderInterface
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) SetupTest() {
	s.syncer = New(log.Discard)
	s.templateProvider = &template.Provider{}
}

func (s *Suite) TestSync() {
	sourcePath := filepath.FromSlash("testdata/TestSync/source")
	destinationPath := filepath.FromSlash("testdata/TestSync/destination")

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

func (s *Suite) TestSyncExecutable() {
	// Irrelevant on Windows
	//goland:noinspection GoBoolExpressions
	if runtime.GOOS == "windows" {
		s.T().Skip()
	}

	sourcePath := filepath.FromSlash("testdata/TestSyncExecutable/source")
	destinationPath := filepath.FromSlash("testdata/TestSyncExecutable/destination")

	_ = os.RemoveAll(destinationPath)
	_ = os.Mkdir(destinationPath, 0755)
	_ = os.WriteFile(filepath.Join(destinationPath, "executable_true"), []byte(""), 0777)
	_ = os.WriteFile(filepath.Join(destinationPath, "executable_false"), []byte(""), 0666)

	s.Run("SourceTrue", func() {
		err := s.syncer.Sync(sourcePath, "executable_true", destinationPath, "executable", nil)
		s.Require().NoError(err)

		stat, _ := os.Stat(filepath.Join(destinationPath, "executable"))
		s.NotEqual(0, int(stat.Mode()&0100))
	})

	s.Run("SourceFalse", func() {
		err := s.syncer.Sync(sourcePath, "executable_false", destinationPath, "executable", nil)
		s.Require().NoError(err)

		stat, _ := os.Stat(filepath.Join(destinationPath, "executable"))
		s.Equal(0, int(stat.Mode()&0100))
	})

	s.Run("SourceFalseDestinationFalse", func() {
		err := s.syncer.Sync(sourcePath, "executable_false", destinationPath, "executable_false", nil)
		s.Require().NoError(err)

		stat, _ := os.Stat(filepath.Join(destinationPath, "executable_false"))
		s.Equal(0, int(stat.Mode()&0100))
	})

	s.Run("SourceTrueDestinationFalse", func() {
		err := s.syncer.Sync(sourcePath, "executable_true", destinationPath, "executable_false", nil)
		s.Require().NoError(err)

		stat, _ := os.Stat(filepath.Join(destinationPath, "executable_false"))
		s.NotEqual(0, int(stat.Mode()&0100))
	})

	s.Run("SourceFalseDestinationTrue", func() {
		err := s.syncer.Sync(sourcePath, "executable_false", destinationPath, "executable_true", nil)
		s.Require().NoError(err)

		stat, _ := os.Stat(filepath.Join(destinationPath, "executable_true"))
		s.Equal(0, int(stat.Mode()&0100))
	})

	s.Run("SourceTrueDestinationTrue", func() {
		err := s.syncer.Sync(sourcePath, "executable_true", destinationPath, "executable_true", nil)
		s.Require().NoError(err)

		stat, _ := os.Stat(filepath.Join(destinationPath, "executable_true"))
		s.NotEqual(0, int(stat.Mode()&0100))
	})
}

func (s *Suite) TestSyncTemplate() {
	sourcePath := filepath.FromSlash("testdata/TestSyncTemplate/source")
	destinationPath := filepath.FromSlash("testdata/TestSyncTemplate/destination")

	_ = os.RemoveAll(destinationPath)
	_ = os.Mkdir(destinationPath, 0755)
	_ = os.WriteFile(filepath.Join(destinationPath, "file_foo"), []byte("foo"), 0666)
	_ = os.WriteFile(filepath.Join(destinationPath, "file_bar"), []byte("bar"), 0666)

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
