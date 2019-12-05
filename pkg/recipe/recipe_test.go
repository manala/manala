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
	s.IsType(&Recipe{}, rec)
	s.Equal("foo", rec.Name)
	s.Equal(".manala.yaml", rec.ConfigFile)
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
	rec, err := Load(repo, "foo")
	s.NoError(err)
	s.IsType(&Recipe{}, rec)
	s.Equal("testdata/load/repository/foo", rec.Dir)
	s.Equal("foo", rec.Name)
	s.Equal("Foo bar", rec.Config.Description)
	s.Equal("bar", rec.Vars["foo"])
}

func (s *LoadTestSuite) TestLoadSync() {
	repo, _ := repository.Load("testdata/load_sync/repository", "")
	rec, err := Load(repo, "foo")
	s.NoError(err)
	s.IsType(&Recipe{}, rec)
	s.Equal([]SyncUnit{
		{Source: "foo", Destination: "foo"},
		{Source: "foo", Destination: "bar"},
	}, rec.Config.Sync)
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

	err := Walk(repo, func(rec *Recipe) {
		results[rec.Name] = rec.Config.Description
	})

	s.NoError(err)
	s.Len(results, 3)
	s.Equal("Foo bar", results["foo"])
	s.Equal("Bar bar", results["bar"])
	s.Equal("Baz bar", results["baz"])
}
