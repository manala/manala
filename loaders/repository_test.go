package loaders

import (
	"github.com/stretchr/testify/suite"
	"manala/config"
	"manala/logger"
	"manala/models"
	"os"
	"runtime"
	"testing"
)

/**********************/
/* Repository - Suite */
/**********************/

type RepositoryTestSuite struct {
	suite.Suite
	ld RepositoryLoaderInterface
}

func TestRepositoryTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(RepositoryTestSuite))
}

func (s *RepositoryTestSuite) SetupTest() {
	cacheDir := "testdata/repository/.cache"
	_ = os.RemoveAll(cacheDir)
	_ = os.Mkdir(cacheDir, 0755)

	conf := config.New(config.WithCacheDir(cacheDir))

	log := logger.New()

	s.ld = NewRepositoryLoader(log, conf)
}

/**********************/
/* Repository - Tests */
/**********************/

func (s *RepositoryTestSuite) TestRepository() {
	s.Implements((*RepositoryLoaderInterface)(nil), s.ld)
}

func (s *RepositoryTestSuite) TestRepositoryLoadDir() {
	repo, err := s.ld.Load("testdata/repository/load_dir")
	s.NoError(err)
	s.Implements((*models.RepositoryInterface)(nil), repo)
	s.Equal("testdata/repository/load_dir", repo.Source())
}

func (s *RepositoryTestSuite) TestRepositoryLoadDirNotFound() {
	repo, err := s.ld.Load("testdata/repository/load_dir_not_found")
	s.Error(err)
	s.Equal("\"testdata/repository/load_dir_not_found\" directory does not exists", err.Error())
	s.Nil(repo)
}

func (s *RepositoryTestSuite) TestRepositoryLoadDirFile() {
	repo, err := s.ld.Load("testdata/repository/load_dir_file")
	s.Error(err)
	s.Equal("\"testdata/repository/load_dir_file\" is not a directory", err.Error())
	s.Nil(repo)
}

func (s *RepositoryTestSuite) TestRepositoryLoadGit() {
	repo, err := s.ld.Load("https://github.com/octocat/Hello-World.git")
	s.NoError(err)
	s.Implements((*models.RepositoryInterface)(nil), repo)
	s.Equal("https://github.com/octocat/Hello-World.git", repo.Source())

	s.DirExists("testdata/repository/.cache/repositories")
	stat, _ := os.Stat("testdata/repository/.cache/repositories")
	if runtime.GOOS == "windows" {
		s.Equal(os.FileMode(0777), stat.Mode().Perm())
	} else {
		s.Equal(os.FileMode(0700), stat.Mode().Perm())
	}

	s.DirExists("testdata/repository/.cache/repositories/1d60d0a17c4d14e9bda84ee53ee51311")
	stat, _ = os.Stat("testdata/repository/.cache/repositories/1d60d0a17c4d14e9bda84ee53ee51311")
	if runtime.GOOS == "windows" {
		s.Equal(os.FileMode(0777), stat.Mode().Perm())
	} else {
		s.Equal(os.FileMode(0700), stat.Mode().Perm())
	}
}

func (s *RepositoryTestSuite) TestRepositoryLoadGitNotExist() {
	repo, err := s.ld.Load("https://github.com/octocat/Foo-Bar.git")
	s.Error(err)
	s.Equal("unable to clone repository: authentication required", err.Error())
	s.Nil(repo)
}
