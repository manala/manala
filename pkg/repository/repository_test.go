package repository

import (
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

/***************/
/* New - Suite */
/***************/

type NewTestSuite struct{ suite.Suite }

func TestNewTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(NewTestSuite))
}

/***************/
/* New - Tests */
/***************/

func (s *NewTestSuite) TestNew() {
	repo := New("foo")
	s.IsType(&Repository{}, repo)
	s.Equal("foo", repo.Src)
}

/****************/
/* Load - Suite */
/****************/

type LoadTestSuite struct{ suite.Suite }

func TestLoadTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(LoadTestSuite))
}

/****************/
/* Load - Tests */
/****************/

func (s *LoadTestSuite) TestLoad() {
	repo, err := Load("foo", "bar")
	s.NoError(err)
	s.IsType(&Repository{}, repo)
	s.Equal("foo", repo.Src)
	s.Equal("foo", repo.Dir)
}

/********************/
/* Load Git - Suite */
/********************/

type LoadGitTestSuite struct{ suite.Suite }

func TestLoadGitTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(LoadGitTestSuite))
}

func (s *LoadGitTestSuite) SetupTest() {
	dir := "testdata/load_git/cache"
	_ = os.RemoveAll(dir)
	_ = os.Mkdir(dir, 0755)
}

/********************/
/* Load Git - Tests */
/********************/

func (s *LoadGitTestSuite) TestLoadGit() {
	repo, err := Load("https://github.com/octocat/Hello-World.git", "testdata/load_git/cache")
	s.NoError(err)
	s.IsType(&Repository{}, repo)
	s.Equal("https://github.com/octocat/Hello-World.git", repo.Src)
	s.Equal("testdata/load_git/cache/repositories/1d60d0a17c4d14e9bda84ee53ee51311", repo.Dir)

	s.DirExists("testdata/load_git/cache/repositories")
	stat, _ := os.Stat("testdata/load_git/cache/repositories")
	s.Equal(os.FileMode(0700), stat.Mode().Perm())

	s.DirExists("testdata/load_git/cache/repositories/1d60d0a17c4d14e9bda84ee53ee51311")
	stat, _ = os.Stat("testdata/load_git/cache/repositories/1d60d0a17c4d14e9bda84ee53ee51311")
	s.Equal(os.FileMode(0700), stat.Mode().Perm())
}
