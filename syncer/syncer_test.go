package syncer

import (
	"github.com/stretchr/testify/suite"
	"manala/fs"
	"manala/logger"
	"manala/models"
	"manala/template"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

/****************/
/* Sync - Suite */
/****************/

type SyncTestSuite struct {
	suite.Suite
	sync            *Syncer
	fsManager       fs.ManagerInterface
	templateManager template.ManagerInterface
}

func TestSyncTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(SyncTestSuite))
}

func (s *SyncTestSuite) SetupTest() {
	log := logger.New()

	s.fsManager = fs.NewManager()
	modelFsManager := models.NewFsManager(s.fsManager)
	s.templateManager = template.NewManager()
	modelTemplateManager := models.NewTemplateManager(s.templateManager, modelFsManager)

	s.sync = New(log, modelFsManager, modelTemplateManager)
}

/****************/
/* Sync - Tests */
/****************/

func (s *SyncTestSuite) TestSync() {
	srcDir := filepath.Join("testdata", "sync", "source")
	dstDir := filepath.Join("testdata", "sync", "destination")

	_ = os.RemoveAll(dstDir)
	_ = os.Mkdir(dstDir, 0755)
	_ = os.WriteFile(filepath.Join(dstDir, "file_foo"), []byte("foo"), 0666)
	_ = os.WriteFile(filepath.Join(dstDir, "file_bar"), []byte("bar"), 0666)
	_ = os.Mkdir(filepath.Join(dstDir, "dir_empty"), 0755)
	_ = os.Mkdir(filepath.Join(dstDir, "dir"), 0755)
	f, _ := os.Create(filepath.Join(dstDir, "dir", "foo"))
	f.Close()
	_ = os.WriteFile(filepath.Join(dstDir, "dir", "foo"), []byte("bar"), 0666)
	_ = os.Mkdir(filepath.Join(dstDir, "dir", "bar"), 0755)
	f, _ = os.Create(filepath.Join(dstDir+"dir", "bar", "foo"))
	f.Close()

	srcFs := s.fsManager.NewDirFs(srcDir)
	dstFs := s.fsManager.NewDirFs(dstDir)

	s.Run("Source not exists", func() {
		err := s.sync.Sync(srcFs, "baz", nil, "baz", dstFs, nil)
		s.IsType(&SourceNotExistError{}, err)
	})

	s.Run("Destination file not exists", func() {
		err := s.sync.Sync(srcFs, "foo", nil, "foo", dstFs, nil)
		s.NoError(err)
		s.FileExists(filepath.Join(dstDir, "foo"))
		content, _ := os.ReadFile(filepath.Join(dstDir, "foo"))
		s.Equal("bar", string(content))
	})

	s.Run("Destination file exists and same", func() {
		err := s.sync.Sync(srcFs, "foo", nil, "file_bar", dstFs, nil)
		s.NoError(err)
		s.FileExists(filepath.Join(dstDir, "file_bar"))
		content, _ := os.ReadFile(filepath.Join(dstDir, "file_bar"))
		s.Equal("bar", string(content))
	})

	s.Run("Destination file exists and different", func() {
		err := s.sync.Sync(srcFs, "foo", nil, "file_foo", dstFs, nil)
		s.NoError(err)
		s.FileExists(filepath.Join(dstDir, "file_foo"))
		content, _ := os.ReadFile(filepath.Join(dstDir, "file_foo"))
		s.Equal("bar", string(content))
	})

	s.Run("Source file over destination directory empty", func() {
		err := s.sync.Sync(srcFs, "foo", nil, "dir_empty", dstFs, nil)
		s.NoError(err)
		s.FileExists(filepath.Join(dstDir, "dir_empty"))
		content, _ := os.ReadFile(filepath.Join(dstDir, "dir_empty"))
		s.Equal("bar", string(content))
	})

	s.Run("Source file over destination directory", func() {
		err := s.sync.Sync(srcFs, "foo", nil, "dir", dstFs, nil)
		s.NoError(err)
		s.FileExists(filepath.Join(dstDir, "dir"))
		content, _ := os.ReadFile(filepath.Join(dstDir, "dir"))
		s.Equal("bar", string(content))
	})

	s.Run("Destination directory not exists", func() {
		err := s.sync.Sync(srcFs, "bar", nil, "bar", dstFs, nil)
		s.NoError(err)
		s.FileExists(filepath.Join(dstDir, "bar", "foo"))
		content, _ := os.ReadFile(filepath.Join(dstDir, "bar", "foo"))
		s.Equal("baz", string(content))
	})

	s.Run("Destination directory exists", func() {
		err := s.sync.Sync(srcFs, "bar", nil, "dir", dstFs, nil)
		s.NoError(err)
		s.FileExists(filepath.Join(dstDir, "dir", "foo"))
		content, _ := os.ReadFile(filepath.Join(dstDir, "dir", "foo"))
		s.Equal("baz", string(content))
	})

	s.Run("Destination file directory not exists", func() {
		err := s.sync.Sync(srcFs, "foo", nil, filepath.Join("baz", "foo"), dstFs, nil)
		s.NoError(err)
		s.FileExists(filepath.Join(dstDir, "baz", "foo"))
		content, _ := os.ReadFile(filepath.Join(dstDir, "baz", "foo"))
		s.Equal("bar", string(content))
	})
}

func (s *SyncTestSuite) TestSyncExecutable() {
	// Irrelevant on windows
	if runtime.GOOS == "windows" {
		s.T().Skip()
	}

	srcDir := filepath.Join("testdata", "sync_executable", "source")
	dstDir := filepath.Join("testdata", "sync_executable", "destination")

	_ = os.RemoveAll(dstDir)
	_ = os.Mkdir(dstDir, 0755)
	_ = os.WriteFile(filepath.Join(dstDir, "executable_true"), []byte(""), 0777)
	_ = os.WriteFile(filepath.Join(dstDir, "executable_false"), []byte(""), 0666)

	srcFs := s.fsManager.NewDirFs(srcDir)
	dstFs := s.fsManager.NewDirFs(dstDir)

	s.Run("Source true", func() {
		err := s.sync.Sync(srcFs, "executable_true", nil, "executable", dstFs, nil)
		s.NoError(err)
		stat, _ := os.Stat(filepath.Join(dstDir, "executable"))
		s.Equal(true, (stat.Mode()&0100) != 0)
	})

	s.Run("Source false", func() {
		err := s.sync.Sync(srcFs, "executable_false", nil, "executable", dstFs, nil)
		s.NoError(err)
		stat, _ := os.Stat(filepath.Join(dstDir, "executable"))
		s.Equal(false, (stat.Mode()&0100) != 0)
	})

	s.Run("Source false destination false", func() {
		err := s.sync.Sync(srcFs, "executable_false", nil, "executable_false", dstFs, nil)
		s.NoError(err)
		stat, _ := os.Stat(filepath.Join(dstDir, "executable_false"))
		s.Equal(false, (stat.Mode()&0100) != 0)
	})

	s.Run("Source true destination false", func() {
		err := s.sync.Sync(srcFs, "executable_true", nil, "executable_false", dstFs, nil)
		s.NoError(err)
		stat, _ := os.Stat(filepath.Join(dstDir, "executable_false"))
		s.Equal(true, (stat.Mode()&0100) != 0)
	})

	s.Run("Source false destination true", func() {
		err := s.sync.Sync(srcFs, "executable_false", nil, "executable_true", dstFs, nil)
		s.NoError(err)
		stat, _ := os.Stat(filepath.Join(dstDir, "executable_true"))
		s.Equal(false, (stat.Mode()&0100) != 0)
	})

	s.Run("Source true destination true", func() {
		err := s.sync.Sync(srcFs, "executable_true", nil, "executable_true", dstFs, nil)
		s.NoError(err)
		stat, _ := os.Stat(filepath.Join(dstDir, "executable_true"))
		s.Equal(true, (stat.Mode()&0100) != 0)
	})
}

func (s *SyncTestSuite) TestSyncTemplate() {
	srcDir := filepath.Join("testdata", "sync_template", "source")
	dstDir := filepath.Join("testdata", "sync_template", "destination")

	_ = os.RemoveAll(dstDir)
	_ = os.Mkdir(dstDir, 0755)
	_ = os.WriteFile(filepath.Join(dstDir, "file_foo"), []byte("foo"), 0666)
	_ = os.WriteFile(filepath.Join(dstDir, "file_bar"), []byte("bar"), 0666)

	srcFs := s.fsManager.NewDirFs(srcDir)
	dstFs := s.fsManager.NewDirFs(dstDir)

	srcTmpl := s.templateManager.NewFsTemplate(srcFs)

	s.Run("Source not exists", func() {
		err := s.sync.Sync(srcFs, "baz.tmpl", srcTmpl, "baz", dstFs, nil)
		s.IsType(&SourceNotExistError{}, err)
	})

	s.Run("Destination file not exists", func() {
		err := s.sync.Sync(srcFs, "foo.tmpl", srcTmpl, "foo", dstFs, nil)
		s.NoError(err)
		s.FileExists(filepath.Join(dstDir, "foo"))
		content, _ := os.ReadFile(filepath.Join(dstDir, "foo"))
		s.Equal("bar", string(content))
	})

	s.Run("Destination file exists and same", func() {
		err := s.sync.Sync(srcFs, "foo.tmpl", srcTmpl, "file_bar", dstFs, nil)
		s.NoError(err)
		s.FileExists(filepath.Join(dstDir, "file_bar"))
		content, _ := os.ReadFile(filepath.Join(dstDir, "file_bar"))
		s.Equal("bar", string(content))
	})

	s.Run("Destination file exists and different", func() {
		err := s.sync.Sync(srcFs, "foo.tmpl", srcTmpl, "file_foo", dstFs, nil)
		s.NoError(err)
		s.FileExists(filepath.Join(dstDir, "file_foo"))
		content, _ := os.ReadFile(filepath.Join(dstDir, "file_foo"))
		s.Equal("bar", string(content))
	})

	s.Run("Invalid", func() {
		err := s.sync.Sync(srcFs, "invalid.tmpl", srcTmpl, "invalid", dstFs, nil)
		s.Error(err)
		s.Contains(err.Error(), "invalid template")
	})
}
