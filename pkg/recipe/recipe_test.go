package recipe

import (
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/suite"
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
	rec := New("testdata/new")
	s.Implements((*Interface)(nil), rec)
	s.Equal("new", rec.GetName())
	s.Equal("testdata/new", rec.GetDir())
	s.Equal("testdata/new/.manala.yaml", rec.GetConfigFile())
	s.True(rec.IsExist())
}

func (s *NewTestSuite) TestNewNotExists() {
	rec := New("testdata/new_not_exists")
	s.False(rec.IsExist())
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
	rec := New("testdata/load")
	err := rec.Load(Config{})
	s.NoError(err)
	s.Implements((*Interface)(nil), rec)
	s.Equal("load", rec.GetName())
	s.Equal("testdata/load", rec.GetDir())
	s.Equal("testdata/load/.manala.yaml", rec.GetConfigFile())
	s.Equal("Foo bar", rec.GetConfig().Description)
	s.Equal("bar", rec.GetVars()["foo"])
}

func (s *LoadTestSuite) TestLoadNotFound() {
	rec := New("testdata/load_not_found")
	err := rec.Load(Config{})
	s.Error(err, "recipe not found")
}

func (s *LoadTestSuite) TestLoadEmpty() {
	rec := New("testdata/load_empty")
	err := rec.Load(Config{})
	s.Error(err)
	s.Contains(err.Error(), "empty recipe config")
}

func (s *LoadTestSuite) TestLoadInvalid() {
	rec := New("testdata/load_invalid")
	err := rec.Load(Config{})
	s.Error(err)
	s.Contains(err.Error(), "invalid recipe config")
}

func (s *LoadTestSuite) TestLoadWithoutDescription() {
	rec := New("testdata/load_without_description")
	err := rec.Load(Config{})
	s.IsType(validator.ValidationErrors{}, err)
}

func (s *LoadTestSuite) TestLoadVars() {
	rec := New("testdata/load_vars")
	_ = rec.Load(Config{})
	s.IsType(map[string]interface{}{}, rec.GetVars(), "vars should be a string map")
}

func (s *LoadTestSuite) TestLoadVarsManala() {
	rec := New("testdata/load_vars")
	_ = rec.Load(Config{})
	_, exists := rec.GetVars()["manala"]
	s.False(exists, "vars should not contain \"manala\" key")
}

func (s *LoadTestSuite) TestLoadVarsStringMap() {
	rec := New("testdata/load_vars")
	_ = rec.Load(Config{})
	s.IsType(map[string]interface{}{}, rec.GetVars()["foo"], "yaml mapping should be mapped as a string map")
	s.IsType(map[string]interface{}{}, rec.GetVars()["bar"], "yaml mapping with anchor should be mapped as a string map")
	s.IsType(map[string]interface{}{}, rec.GetVars()["baz"], "yaml mapping with alias should be mapped as a string map")
}

func (s *LoadTestSuite) TestLoadSync() {
	rec := New("testdata/load_sync")
	err := rec.Load(Config{})
	s.NoError(err)
	s.Equal([]SyncUnit{
		{Source: "foo", Destination: "foo"},
		{Source: "foo", Destination: "bar"},
	}, rec.GetConfig().Sync)
}
