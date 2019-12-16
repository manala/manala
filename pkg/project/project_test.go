package project

import (
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
	"io"
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
	prj := New("testdata/new")
	s.IsType(&project{}, prj)
	s.Equal("testdata/new", prj.GetDir())
	s.Equal("testdata/new/.manala.yaml", prj.GetConfigFile())
	s.True(prj.IsExist())
}

func (s *NewTestSuite) TestNewNotExists() {
	prj := New("testdata/new_not_exists")
	s.False(prj.IsExist())
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
	prj := New("testdata/load")
	err := prj.Load(Config{
		Repository: "bar",
	})
	s.NoError(err)
	s.Equal("testdata/load", prj.GetDir())
	s.Equal("testdata/load/.manala.yaml", prj.GetConfigFile())
	s.Equal("foo", prj.GetConfig().Recipe)
	s.Equal("bar", prj.GetConfig().Repository)
	s.Equal("bar", prj.GetVars()["foo"])
}

func (s *LoadTestSuite) TestLoadNotFound() {
	prj := New("testdata/load_not_found")
	err := prj.Load(Config{})
	s.Error(err, "project not found")
}

func (s *LoadTestSuite) TestLoadEmpty() {
	prj := New("testdata/load_empty")
	err := prj.Load(Config{})
	s.Equal(io.EOF, err)
}

func (s *LoadTestSuite) TestLoadInvalid() {
	prj := New("testdata/load_invalid")
	err := prj.Load(Config{})
	s.IsType(&yaml.TypeError{}, err)
}

func (s *LoadTestSuite) TestLoadWithoutRecipe() {
	prj := New("testdata/load_without_recipe")
	err := prj.Load(Config{})
	s.IsType(validator.ValidationErrors{}, err)
}

func (s *LoadTestSuite) TestLoadWithRepository() {
	prj := New("testdata/load_with_repository")
	err := prj.Load(Config{})
	s.NoError(err)
	s.Equal("baz", prj.GetConfig().Repository)
}
