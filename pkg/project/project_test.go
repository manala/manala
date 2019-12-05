package project

import (
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
	"io"
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
	prj := New("testdata/new")
	s.IsType(&Project{}, prj)
	s.Equal("testdata/new", prj.Dir)
	s.Equal(".manala.yaml", prj.ConfigFile)
	s.True(prj.IsExist())
}

func (s *NewTestSuite) TestNewNotExists() {
	prj := New("testdata/new_not_exists")
	s.IsType(&Project{}, prj)
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
	prj, err := Load("testdata/load", "bar")
	s.NoError(err)
	s.IsType(&Project{}, prj)
	s.Equal("testdata/load", prj.Dir)
	s.Equal(".manala.yaml", prj.ConfigFile)
	s.Equal("foo", prj.Config.Recipe)
	s.Equal("bar", prj.Config.Repository)
	s.Equal("bar", prj.Vars["foo"])
}

func (s *LoadTestSuite) TestLoadNotFound() {
	prj, err := Load("testdata/load_not_found", "bar")
	s.IsType(&os.PathError{}, err)
	s.Nil(prj)
}

func (s *LoadTestSuite) TestLoadEmpty() {
	prj, err := Load("testdata/load_empty", "bar")
	s.Equal(io.EOF, err)
	s.Nil(prj)
}

func (s *LoadTestSuite) TestLoadInvalid() {
	prj, err := Load("testdata/load_invalid", "bar")
	s.IsType(&yaml.TypeError{}, err)
	s.Nil(prj)
}

func (s *LoadTestSuite) TestLoadWithoutRecipe() {
	prj, err := Load("testdata/load_without_recipe", "bar")
	s.IsType(validator.ValidationErrors{}, err)
	s.Nil(prj)
}

func (s *LoadTestSuite) TestLoadWithRepository() {
	prj, err := Load("testdata/load_with_repository", "bar")
	s.NoError(err)
	s.IsType(&Project{}, prj)
	s.Equal("baz", prj.Config.Repository)
}
