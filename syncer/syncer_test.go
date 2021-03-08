package syncer

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/discard"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

/****************/
/* Sync - Suite */
/****************/

type SyncTestSuite struct{ suite.Suite }

func TestSyncTestSuite(t *testing.T) {
	// Discard logs
	log.SetHandler(discard.Default)
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
}

/****************/
/* Sync - Tests */
/****************/

func (s *SyncTestSuite) TestSyncSourceNotExists() {
	err := Sync("testdata/sync/source/baz", "testdata/sync/destination/baz", NewTemplate(), nil)
	s.IsType(&SourceNotExistError{}, err)
}

func (s *SyncTestSuite) TestSyncDestinationFileNotExists() {
	err := Sync("testdata/sync/source/foo", "testdata/sync/destination/foo", NewTemplate(), nil)
	s.NoError(err)
	s.FileExists("testdata/sync/destination/foo")
	content, _ := os.ReadFile("testdata/sync/destination/foo")
	s.Equal("bar", string(content))
}

func (s *SyncTestSuite) TestSyncDestinationFileExistsAndSame() {
	err := Sync("testdata/sync/source/foo", "testdata/sync/destination/file_bar", NewTemplate(), nil)
	s.NoError(err)
	s.FileExists("testdata/sync/destination/file_bar")
	content, _ := os.ReadFile("testdata/sync/destination/file_bar")
	s.Equal("bar", string(content))
}

func (s *SyncTestSuite) TestSyncDestinationFileExistsAndDifferent() {
	err := Sync("testdata/sync/source/foo", "testdata/sync/destination/file_foo", NewTemplate(), nil)
	s.NoError(err)
	s.FileExists("testdata/sync/destination/file_foo")
	content, _ := os.ReadFile("testdata/sync/destination/file_foo")
	s.Equal("bar", string(content))
}

func (s *SyncTestSuite) TestSyncSourceFileOverDestinationDirectoryEmpty() {
	err := Sync("testdata/sync/source/foo", "testdata/sync/destination/dir_empty", NewTemplate(), nil)
	s.NoError(err)
	s.FileExists("testdata/sync/destination/dir_empty")
	content, _ := os.ReadFile("testdata/sync/destination/dir_empty")
	s.Equal("bar", string(content))
}

func (s *SyncTestSuite) TestSyncSourceFileOverDestinationDirectory() {
	err := Sync("testdata/sync/source/foo", "testdata/sync/destination/dir", NewTemplate(), nil)
	s.NoError(err)
	s.FileExists("testdata/sync/destination/dir")
	content, _ := os.ReadFile("testdata/sync/destination/dir")
	s.Equal("bar", string(content))
}

func (s *SyncTestSuite) TestSyncDestinationDirectoryNotExists() {
	err := Sync("testdata/sync/source/bar", "testdata/sync/destination/bar", NewTemplate(), nil)
	s.NoError(err)
	s.FileExists("testdata/sync/destination/bar/foo")
	content, _ := os.ReadFile("testdata/sync/destination/bar/foo")
	s.Equal("baz", string(content))
}

func (s *SyncTestSuite) TestSyncDestinationDirectoryExists() {
	err := Sync("testdata/sync/source/bar", "testdata/sync/destination/dir", NewTemplate(), nil)
	s.NoError(err)
	s.FileExists("testdata/sync/destination/dir/foo")
	content, _ := os.ReadFile("testdata/sync/destination/dir/foo")
	s.Equal("baz", string(content))
}

/***************************/
/* Sync Executable - Suite */
/***************************/

type SyncExecutableTestSuite struct{ suite.Suite }

func TestSyncExecutableTestSuite(t *testing.T) {
	// Discard logs
	log.SetHandler(discard.Default)
	// Run
	suite.Run(t, new(SyncExecutableTestSuite))
}

func (s *SyncExecutableTestSuite) SetupTest() {
	dir := "testdata/sync_executable/destination"
	_ = os.RemoveAll(dir)
	_ = os.Mkdir(dir, 0755)
	_ = os.WriteFile(dir+"/executable_true", []byte(""), 0777)
	_ = os.WriteFile(dir+"/executable_false", []byte(""), 0666)
}

/***************************/
/* Sync Executable - Tests */
/***************************/

func (s *SyncExecutableTestSuite) TestSyncExecutableSourceTrue() {
	err := Sync("testdata/sync_executable/source/executable_true", "testdata/sync_executable/destination/executable", NewTemplate(), nil)
	s.NoError(err)
	stat, _ := os.Stat("testdata/sync_executable/destination/executable")
	s.Equal(true, (stat.Mode()&0100) != 0)
}

func (s *SyncExecutableTestSuite) TestSyncExecutableSourceFalse() {
	err := Sync("testdata/sync_executable/source/executable_false", "testdata/sync_executable/destination/executable", NewTemplate(), nil)
	s.NoError(err)
	stat, _ := os.Stat("testdata/sync_executable/destination/executable")
	s.Equal(false, (stat.Mode()&0100) != 0)
}

func (s *SyncExecutableTestSuite) TestSyncExecutableSourceFalseDestinationFalse() {
	err := Sync("testdata/sync_executable/source/executable_false", "testdata/sync_executable/destination/executable_false", NewTemplate(), nil)
	s.NoError(err)
	stat, _ := os.Stat("testdata/sync_executable/destination/executable_false")
	s.Equal(false, (stat.Mode()&0100) != 0)
}

func (s *SyncExecutableTestSuite) TestSyncExecutableSourceTrueDestinationFalse() {
	err := Sync("testdata/sync_executable/source/executable_true", "testdata/sync_executable/destination/executable_false", NewTemplate(), nil)
	s.NoError(err)
	stat, _ := os.Stat("testdata/sync_executable/destination/executable_false")
	s.Equal(true, (stat.Mode()&0100) != 0)
}

func (s *SyncExecutableTestSuite) TestSyncExecutableSourceFalseDestinationTrue() {
	err := Sync("testdata/sync_executable/source/executable_false", "testdata/sync_executable/destination/executable_true", NewTemplate(), nil)
	s.NoError(err)
	stat, _ := os.Stat("testdata/sync_executable/destination/executable_true")
	s.Equal(false, (stat.Mode()&0100) != 0)
}

func (s *SyncExecutableTestSuite) TestSyncExecutableSourceTrueDestinationTrue() {
	err := Sync("testdata/sync_executable/source/executable_true", "testdata/sync_executable/destination/executable_true", NewTemplate(), nil)
	s.NoError(err)
	stat, _ := os.Stat("testdata/sync_executable/destination/executable_true")
	s.Equal(true, (stat.Mode()&0100) != 0)
}

/*************************/
/* Sync Template - Suite */
/*************************/

type SyncTemplateTestSuite struct{ suite.Suite }

func TestSyncTemplateTestSuite(t *testing.T) {
	// Discard logs
	log.SetHandler(discard.Default)
	// Run
	suite.Run(t, new(SyncTemplateTestSuite))
}

func (s *SyncTemplateTestSuite) SetupTest() {
	dir := "testdata/sync_template/destination"
	_ = os.RemoveAll(dir)
	_ = os.Mkdir(dir, 0755)
}

/*************************/
/* Sync Template - Tests */
/*************************/

func (s *SyncTemplateTestSuite) TestSyncTemplateBase() {
	err := Sync("testdata/sync_template/source/base.tmpl", "testdata/sync_template/destination/base", NewTemplate(), nil)
	s.NoError(err)
	s.FileExists("testdata/sync_template/destination/base")
	content, _ := os.ReadFile("testdata/sync_template/destination/base")
	s.Equal(`foo
`, string(content))
}

func (s *SyncTemplateTestSuite) TestSyncTemplateInvalid() {
	err := Sync("testdata/sync_template/source/invalid.tmpl", "testdata/sync_template/destination/invalid", NewTemplate(), nil)
	s.Error(err)
	s.Contains(err.Error(), "invalid template")
}

func (s *SyncTemplateTestSuite) TestSyncTemplateToYaml() {
	err := Sync("testdata/sync_template/source/to_yaml.tmpl", "testdata/sync_template/destination/to_yaml", NewTemplate(), map[string]interface{}{
		"foo": map[string]interface{}{
			"bar": "string",
			"baz": struct {
				Foo string
				Bar int
			}{
				Foo: "foo",
				Bar: 123,
			},
			"qux":    123,
			"quux":   true,
			"corge":  false,
			"grault": 1.23,
			"garply": map[string]interface{}{},
			"waldo": map[string]interface{}{
				"foo": "bar",
				"bar": "baz",
			},
			"fred": []interface{}{},
			"plugh": []interface{}{
				"foo",
				"bar",
			},
			"xyzzy": nil,
			"thud":  "123",
		},
	})
	s.NoError(err)
	s.FileExists("testdata/sync_template/destination/to_yaml")
	content, _ := os.ReadFile("testdata/sync_template/destination/to_yaml")
	s.Equal(`foo:
    bar: string
    baz:
        foo: foo
        bar: 123
    corge: false
    fred: []
    garply: {}
    grault: 1.23
    plugh:
        - foo
        - bar
    quux: true
    qux: 123
    thud: "123"
    waldo:
        bar: baz
        foo: bar
    xyzzy: null
`, string(content))
}

func (s *SyncTemplateTestSuite) TestSyncTemplateCases() {
	err := Sync("testdata/sync_template/source/cases.tmpl", "testdata/sync_template/destination/cases", NewTemplate(), map[string]interface{}{
		"foo": map[string]interface{}{
			"bar":  true,
			"BAZ":  true,
			"qUx":  true,
			"QuuX": true,
		},
	})
	s.NoError(err)
	s.FileExists("testdata/sync_template/destination/cases")
	content, _ := os.ReadFile("testdata/sync_template/destination/cases")
	s.Equal(`foo:
    BAZ: true
    QuuX: true
    bar: true
    qUx: true
`, string(content))
}

func (s *SyncTemplateTestSuite) TestSyncTemplateDict() {
	err := Sync("testdata/sync_template/source/dict.tmpl", "testdata/sync_template/destination/dict", NewTemplate(), map[string]interface{}{
		"foo": map[string]interface{}{
			"bar": true,
			"baz": true,
			"qux": true,
		},
	})
	s.NoError(err)
	s.FileExists("testdata/sync_template/destination/dict")
	content, _ := os.ReadFile("testdata/sync_template/destination/dict")
	s.Equal(`bar: true
qux: true
`, string(content))
}

func (s *SyncTemplateTestSuite) TestSyncTemplateInclude() {
	err := Sync("testdata/sync_template/source/include.tmpl", "testdata/sync_template/destination/include", NewTemplate(), nil)
	s.NoError(err)
	s.FileExists("testdata/sync_template/destination/include")
	content, _ := os.ReadFile("testdata/sync_template/destination/include")
	s.Equal(`foo: bar`, string(content))
}

func (s *SyncTemplateTestSuite) TestSyncTemplateHelpers() {
	tmpl := NewTemplate()
	_, _ = tmpl.ParseFiles("testdata/sync_template/source/_helpers.tmpl")
	err := Sync("testdata/sync_template/source/helpers.tmpl", "testdata/sync_template/destination/helpers", tmpl, nil)
	s.NoError(err)
	s.FileExists("testdata/sync_template/destination/helpers")
	content, _ := os.ReadFile("testdata/sync_template/destination/helpers")
	s.Equal(`bar: foo`, string(content))
}
