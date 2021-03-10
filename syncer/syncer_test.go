package syncer

import (
	"bytes"
	"github.com/stretchr/testify/suite"
	"manala/config"
	"manala/logger"
	"manala/template"
	"os"
	"testing"
)

/****************/
/* Sync - Suite */
/****************/

type SyncTestSuite struct {
	suite.Suite
	sync *Syncer
}

func TestSyncTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(SyncTestSuite))
}

func (s *SyncTestSuite) SetupTest() {
	dir := "testdata/sync/destination"
	_ = os.RemoveAll(dir)
	_ = os.Mkdir(dir, 0755)
	_ = os.WriteFile(dir+"/file_foo", []byte("foo"), 0666)
	_ = os.WriteFile(dir+"/file_bar", []byte("bar"), 0666)
	_ = os.Mkdir(dir+"/dir_empty", 0755)
	_ = os.Mkdir(dir+"/dir", 0755)
	_, _ = os.Create(dir + "/dir/foo")
	_ = os.WriteFile(dir+"/dir/foo", []byte("bar"), 0666)
	_ = os.Mkdir(dir+"/dir/bar", 0755)
	_, _ = os.Create(dir + "/dir/bar/foo")

	conf := config.New("test", "foo")

	log := logger.New(conf)
	log.SetOut(bytes.NewBufferString(""))

	tmpl := template.New()

	s.sync = New(log, tmpl)
}

/****************/
/* Sync - Tests */
/****************/

func (s *SyncTestSuite) TestSyncSourceNotExists() {
	err := s.sync.Sync("testdata/sync/source/baz", "testdata/sync/destination/baz", nil)
	s.IsType(&SourceNotExistError{}, err)
}

func (s *SyncTestSuite) TestSyncDestinationFileNotExists() {
	err := s.sync.Sync("testdata/sync/source/foo", "testdata/sync/destination/foo", nil)
	s.NoError(err)
	s.FileExists("testdata/sync/destination/foo")
	content, _ := os.ReadFile("testdata/sync/destination/foo")
	s.Equal("bar", string(content))
}

func (s *SyncTestSuite) TestSyncDestinationFileExistsAndSame() {
	err := s.sync.Sync("testdata/sync/source/foo", "testdata/sync/destination/file_bar", nil)
	s.NoError(err)
	s.FileExists("testdata/sync/destination/file_bar")
	content, _ := os.ReadFile("testdata/sync/destination/file_bar")
	s.Equal("bar", string(content))
}

func (s *SyncTestSuite) TestSyncDestinationFileExistsAndDifferent() {
	err := s.sync.Sync("testdata/sync/source/foo", "testdata/sync/destination/file_foo", nil)
	s.NoError(err)
	s.FileExists("testdata/sync/destination/file_foo")
	content, _ := os.ReadFile("testdata/sync/destination/file_foo")
	s.Equal("bar", string(content))
}

func (s *SyncTestSuite) TestSyncSourceFileOverDestinationDirectoryEmpty() {
	err := s.sync.Sync("testdata/sync/source/foo", "testdata/sync/destination/dir_empty", nil)
	s.NoError(err)
	s.FileExists("testdata/sync/destination/dir_empty")
	content, _ := os.ReadFile("testdata/sync/destination/dir_empty")
	s.Equal("bar", string(content))
}

func (s *SyncTestSuite) TestSyncSourceFileOverDestinationDirectory() {
	err := s.sync.Sync("testdata/sync/source/foo", "testdata/sync/destination/dir", nil)
	s.NoError(err)
	s.FileExists("testdata/sync/destination/dir")
	content, _ := os.ReadFile("testdata/sync/destination/dir")
	s.Equal("bar", string(content))
}

func (s *SyncTestSuite) TestSyncDestinationDirectoryNotExists() {
	err := s.sync.Sync("testdata/sync/source/bar", "testdata/sync/destination/bar", nil)
	s.NoError(err)
	s.FileExists("testdata/sync/destination/bar/foo")
	content, _ := os.ReadFile("testdata/sync/destination/bar/foo")
	s.Equal("baz", string(content))
}

func (s *SyncTestSuite) TestSyncDestinationDirectoryExists() {
	err := s.sync.Sync("testdata/sync/source/bar", "testdata/sync/destination/dir", nil)
	s.NoError(err)
	s.FileExists("testdata/sync/destination/dir/foo")
	content, _ := os.ReadFile("testdata/sync/destination/dir/foo")
	s.Equal("baz", string(content))
}

/***************************/
/* Sync Executable - Suite */
/***************************/

type SyncExecutableTestSuite struct {
	suite.Suite
	sync *Syncer
}

func TestSyncExecutableTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(SyncExecutableTestSuite))
}

func (s *SyncExecutableTestSuite) SetupTest() {
	dir := "testdata/sync_executable/destination"
	_ = os.RemoveAll(dir)
	_ = os.Mkdir(dir, 0755)
	_ = os.WriteFile(dir+"/executable_true", []byte(""), 0777)
	_ = os.WriteFile(dir+"/executable_false", []byte(""), 0666)

	conf := config.New("test", "foo")

	log := logger.New(conf)
	log.SetOut(bytes.NewBufferString(""))

	tmpl := template.New()

	s.sync = New(log, tmpl)
}

/***************************/
/* Sync Executable - Tests */
/***************************/

func (s *SyncExecutableTestSuite) TestSyncExecutableSourceTrue() {
	err := s.sync.Sync("testdata/sync_executable/source/executable_true", "testdata/sync_executable/destination/executable", nil)
	s.NoError(err)
	stat, _ := os.Stat("testdata/sync_executable/destination/executable")
	s.Equal(true, (stat.Mode()&0100) != 0)
}

func (s *SyncExecutableTestSuite) TestSyncExecutableSourceFalse() {
	err := s.sync.Sync("testdata/sync_executable/source/executable_false", "testdata/sync_executable/destination/executable", nil)
	s.NoError(err)
	stat, _ := os.Stat("testdata/sync_executable/destination/executable")
	s.Equal(false, (stat.Mode()&0100) != 0)
}

func (s *SyncExecutableTestSuite) TestSyncExecutableSourceFalseDestinationFalse() {
	err := s.sync.Sync("testdata/sync_executable/source/executable_false", "testdata/sync_executable/destination/executable_false", nil)
	s.NoError(err)
	stat, _ := os.Stat("testdata/sync_executable/destination/executable_false")
	s.Equal(false, (stat.Mode()&0100) != 0)
}

func (s *SyncExecutableTestSuite) TestSyncExecutableSourceTrueDestinationFalse() {
	err := s.sync.Sync("testdata/sync_executable/source/executable_true", "testdata/sync_executable/destination/executable_false", nil)
	s.NoError(err)
	stat, _ := os.Stat("testdata/sync_executable/destination/executable_false")
	s.Equal(true, (stat.Mode()&0100) != 0)
}

func (s *SyncExecutableTestSuite) TestSyncExecutableSourceFalseDestinationTrue() {
	err := s.sync.Sync("testdata/sync_executable/source/executable_false", "testdata/sync_executable/destination/executable_true", nil)
	s.NoError(err)
	stat, _ := os.Stat("testdata/sync_executable/destination/executable_true")
	s.Equal(false, (stat.Mode()&0100) != 0)
}

func (s *SyncExecutableTestSuite) TestSyncExecutableSourceTrueDestinationTrue() {
	err := s.sync.Sync("testdata/sync_executable/source/executable_true", "testdata/sync_executable/destination/executable_true", nil)
	s.NoError(err)
	stat, _ := os.Stat("testdata/sync_executable/destination/executable_true")
	s.Equal(true, (stat.Mode()&0100) != 0)
}

/*************************/
/* Sync Template - Suite */
/*************************/

type SyncTemplateTestSuite struct {
	suite.Suite
	sync *Syncer
}

func TestSyncTemplateTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(SyncTemplateTestSuite))
}

func (s *SyncTemplateTestSuite) SetupTest() {
	dir := "testdata/sync_template/destination"
	_ = os.RemoveAll(dir)
	_ = os.Mkdir(dir, 0755)

	conf := config.New("test", "foo")

	log := logger.New(conf)
	log.SetOut(bytes.NewBufferString(""))

	tmpl := template.New()

	s.sync = New(log, tmpl)
}

/*************************/
/* Sync Template - Tests */
/*************************/

func (s *SyncTemplateTestSuite) TestSyncTemplateBase() {
	err := s.sync.Sync("testdata/sync_template/source/base.tmpl", "testdata/sync_template/destination/base", nil)
	s.NoError(err)
	s.FileExists("testdata/sync_template/destination/base")
	content, _ := os.ReadFile("testdata/sync_template/destination/base")
	s.Equal(`foo
`, string(content))
}

func (s *SyncTemplateTestSuite) TestSyncTemplateInvalid() {
	err := s.sync.Sync("testdata/sync_template/source/invalid.tmpl", "testdata/sync_template/destination/invalid", nil)
	s.Error(err)
	s.Contains(err.Error(), "invalid template")
}
