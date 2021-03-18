package fs

import (
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

/*********/
/* Suite */
/*********/

type FsTestSuite struct {
	suite.Suite
	manager ManagerInterface
}

func TestFsTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(FsTestSuite))
}

func (s *FsTestSuite) SetupTest() {
	s.manager = NewManager()
}

/*********/
/* Tests */
/*********/

func (s *FsTestSuite) TestRead() {
	dir := "testdata/read"

	fs := s.manager.NewDirFs(dir)

	s.Run("Open", func() {
		_, err := fs.Open("file")
		s.NoError(err)
	})

	s.Run("Stat File", func() {
		stat, err := fs.Stat("file")
		s.NoError(err)
		s.Equal("file", stat.Name())
		s.False(stat.IsDir())
	})

	s.Run("Stat Dir", func() {
		stat, err := fs.Stat("dir")
		s.NoError(err)
		s.Equal("dir", stat.Name())
		s.True(stat.IsDir())
	})

	s.Run("ReadFile", func() {
		content, err := fs.ReadFile("file")
		s.NoError(err)
		s.Equal("file", string(content))
	})

	s.Run("ReadDir", func() {
		files, err := fs.ReadDir("dir")
		s.NoError(err)

		var entries []struct {
			name string
			dir  bool
		}
		for _, file := range files {
			entries = append(entries, struct {
				name string
				dir  bool
			}{
				file.Name(),
				file.IsDir(),
			})
		}
		s.Equal([]struct {
			name string
			dir  bool
		}{
			{"dir", true},
			{"file", false},
		}, entries)
	})
}

func (s *FsTestSuite) TestWrite() {
	dir := "testdata/write"

	// Clean
	_ = os.RemoveAll(dir)
	_ = os.Mkdir(dir, 0755)
	_ = os.WriteFile(dir+"/file_chmod", []byte(""), 0666)
	_ = os.WriteFile(dir+"/file_remove", []byte(""), 0666)
	_ = os.Mkdir(dir+"/dir_remove", 0755)

	fs := s.manager.NewDirFs(dir)

	s.Run("OpenFile", func() {
		_, err := fs.OpenFile("file", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		s.NoError(err)
		s.FileExists(dir + "/file")
	})

	s.Run("Chmod", func() {
		err := fs.Chmod("file_chmod", 0777)
		s.NoError(err)
		stat, _ := os.Stat(dir + "/file_chmod")
		s.Equal("-rwxrwxrwx", stat.Mode().String())
	})

	s.Run("Remove", func() {
		err := fs.Remove("file_remove")
		s.NoError(err)
		s.NoFileExists(dir + "/file_remove")
	})

	s.Run("MkdirAll", func() {
		err := fs.MkdirAll("dir", 0755)
		s.NoError(err)
		s.DirExists(dir + "/dir")
		stat, _ := os.Stat(dir + "/dir")
		s.Equal("drwxr-xr-x", stat.Mode().String())
	})

	s.Run("RemoveAll", func() {
		err := fs.RemoveAll("dir_remove")
		s.NoError(err)
		s.NoDirExists(dir + "/dir_remove")
	})
}
