package loaders

import (
	"github.com/stretchr/testify/suite"
	"manala/models"
	"os"
	"testing"
)

/**********************/
/* Repository - Suite */
/**********************/

type RepositoryTestSuite struct {
	suite.Suite
	cacheDir string
}

func TestRepositoryTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(RepositoryTestSuite))
}

func (s *RepositoryTestSuite) SetupTest() {
	s.cacheDir = "testdata/repository/.cache"
	_ = os.RemoveAll(s.cacheDir)
	_ = os.Mkdir(s.cacheDir, 0755)
}

/**********************/
/* Repository - Tests */
/**********************/

func (s *RepositoryTestSuite) TestRepository() {
	ld := NewRepositoryLoader(s.cacheDir)
	s.Implements((*RepositoryLoaderInterface)(nil), ld)
}

func (s *RepositoryTestSuite) TestRepositoryLoadDir() {
	ld := NewRepositoryLoader(s.cacheDir)
	repo, err := ld.Load("testdata/repository/load_dir")
	s.NoError(err)
	s.Implements((*models.RepositoryInterface)(nil), repo)
	s.Equal("testdata/repository/load_dir", repo.Src())
	s.Equal("testdata/repository/load_dir", repo.Dir())
}

func (s *RepositoryTestSuite) TestRepositoryLoadDirNotFound() {
	ld := NewRepositoryLoader(s.cacheDir)
	repo, err := ld.Load("testdata/repository/load_dir_not_found")
	s.Error(err)
	s.Equal("\"testdata/repository/load_dir_not_found\" directory does not exists", err.Error())
	s.Nil(repo)
}

func (s *RepositoryTestSuite) TestRepositoryLoadDirFile() {
	ld := NewRepositoryLoader(s.cacheDir)
	repo, err := ld.Load("testdata/repository/load_dir_file")
	s.Error(err)
	s.Equal("\"testdata/repository/load_dir_file\" is not a directory", err.Error())
	s.Nil(repo)
}

func (s *RepositoryTestSuite) TestRepositoryLoadGit() {
	ld := NewRepositoryLoader(s.cacheDir)
	repo, err := ld.Load("https://github.com/octocat/Hello-World.git")
	s.NoError(err)
	s.Implements((*models.RepositoryInterface)(nil), repo)
	s.Equal("https://github.com/octocat/Hello-World.git", repo.Src())
	s.Equal("testdata/repository/.cache/repositories/1d60d0a17c4d14e9bda84ee53ee51311", repo.Dir())

	s.DirExists("testdata/repository/.cache/repositories")
	stat, _ := os.Stat("testdata/repository/.cache/repositories")
	s.Equal(os.FileMode(0700), stat.Mode().Perm())

	s.DirExists("testdata/repository/.cache/repositories/1d60d0a17c4d14e9bda84ee53ee51311")
	stat, _ = os.Stat("testdata/repository/.cache/repositories/1d60d0a17c4d14e9bda84ee53ee51311")
	s.Equal(os.FileMode(0700), stat.Mode().Perm())
}

func (s *RepositoryTestSuite) TestRepositoryLoadGitNotExist() {
	ld := NewRepositoryLoader(s.cacheDir)
	repo, err := ld.Load("https://github.com/octocat/Foo-Bar.git")
	s.Error(err)
	s.Equal("unable to clone repository: authentication required", err.Error())
	s.Nil(repo)
}
