package recipe

import (
	"github.com/stretchr/testify/suite"
	"manala/pkg/repository"
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
	rec := New("foo")
	s.IsType(&recipe{}, rec)
	s.Equal("foo", rec.GetName())
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
	repo, _ := repository.Load("testdata/load/repository", "")
	rec := New("foo")
	err := rec.Load(repo)
	s.NoError(err)
	s.IsType(&recipe{}, rec)
	s.Equal("testdata/load/repository/foo", rec.GetDir())
	s.Equal("foo", rec.GetName())
	s.Equal("Foo bar", rec.GetConfig().Description)
	s.Equal("bar", rec.GetVars()["foo"])
}

func (s *LoadTestSuite) TestLoadSync() {
	repo, _ := repository.Load("testdata/load_sync/repository", "")
	rec := New("foo")
	err := rec.Load(repo)
	s.NoError(err)
	s.Equal([]SyncUnit{
		{Source: "foo", Destination: "foo"},
		{Source: "foo", Destination: "bar"},
	}, rec.GetConfig().Sync)
}

/****************/
/* Walk - Suite */
/****************/

type WalkTestSuite struct{ suite.Suite }

func TestWalkTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(WalkTestSuite))
}

/****************/
/* Walk - Tests */
/****************/

func (s *WalkTestSuite) TestWalk() {
	repo, _ := repository.Load("testdata/walk/repository", "")

	results := make(map[string]string)

	err := Walk(repo, func(rec Interface) {
		results[rec.GetName()] = rec.GetConfig().Description
	})

	s.NoError(err)
	s.Len(results, 3)
	s.Equal("Foo bar", results["foo"])
	s.Equal("Bar bar", results["bar"])
	s.Equal("Baz bar", results["baz"])
}
