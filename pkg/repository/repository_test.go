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
	s.IsType(&repository{}, repo)
	s.Equal("foo", repo.GetSrc())
}

/****************/
/* Load - Suite */
/****************/

type LoadTestSuite struct{ suite.Suite }

func TestLoadTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(LoadTestSuite))
}

func (s *LoadTestSuite) SetupTest() {
	dir := "testdata/load/cache"
	_ = os.RemoveAll(dir)
	_ = os.Mkdir(dir, 0755)
}

/****************/
/* Load - Tests */
/****************/

func (s *LoadTestSuite) TestLoadDir() {
	repo := New("foo")
	err := repo.Load("bar")
	s.NoError(err)
	s.Equal("foo", repo.GetDir())
}

func (s *LoadTestSuite) TestLoadGit() {
	repo := New("https://github.com/octocat/Hello-World.git")
	err := repo.Load("testdata/load/cache")
	s.NoError(err)
	s.Equal("https://github.com/octocat/Hello-World.git", repo.GetSrc())
	s.Equal("testdata/load/cache/repositories/1d60d0a17c4d14e9bda84ee53ee51311", repo.GetDir())

	s.DirExists("testdata/load/cache/repositories")
	stat, _ := os.Stat("testdata/load/cache/repositories")
	s.Equal(os.FileMode(0700), stat.Mode().Perm())

	s.DirExists("testdata/load/cache/repositories/1d60d0a17c4d14e9bda84ee53ee51311")
	stat, _ = os.Stat("testdata/load/cache/repositories/1d60d0a17c4d14e9bda84ee53ee51311")
	s.Equal(os.FileMode(0700), stat.Mode().Perm())
}
