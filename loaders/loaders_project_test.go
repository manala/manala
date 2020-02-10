package loaders

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/discard"
	"github.com/stretchr/testify/suite"
	"manala/models"
	"os"
	"testing"
)

/*******************/
/* Project - Suite */
/*******************/

type ProjectTestSuite struct {
	suite.Suite
	repositoryLoader RepositoryLoaderInterface
	recipeLoader     RecipeLoaderInterface
	repositorySrc    string
}

func TestProjectTestSuite(t *testing.T) {
	// Discard logs
	log.SetHandler(discard.Default)
	// Run
	suite.Run(t, new(ProjectTestSuite))
}

func (s *ProjectTestSuite) SetupTest() {
	cacheDir := "testdata/project/.cache"
	_ = os.RemoveAll(cacheDir)
	_ = os.Mkdir(cacheDir, 0755)
	s.repositoryLoader = NewRepositoryLoader(cacheDir)
	s.recipeLoader = NewRecipeLoader()
	s.repositorySrc = "testdata/project/_repository"
}

/*******************/
/* Project - Tests */
/*******************/

func (s *ProjectTestSuite) TestProject() {
	ld := NewProjectLoader(s.repositoryLoader, s.recipeLoader, s.repositorySrc)
	s.Implements((*ProjectLoaderInterface)(nil), ld)
}

func (s *ProjectTestSuite) TestProjectConfigFile() {
	ld := NewProjectLoader(s.repositoryLoader, s.recipeLoader, s.repositorySrc)
	file, err := ld.ConfigFile("testdata/project/config_file")
	s.NoError(err)
	s.IsType((*os.File)(nil), file)
	s.Equal("testdata/project/config_file/.manala.yaml", file.Name())
}

func (s *ProjectTestSuite) TestProjectConfigFileNotFound() {
	ld := NewProjectLoader(s.repositoryLoader, s.recipeLoader, s.repositorySrc)
	file, err := ld.ConfigFile("testdata/project/config_file_not_found")
	s.Error(err)
	s.Equal("open testdata/project/config_file_not_found/.manala.yaml: no such file or directory", err.Error())
	s.Nil(file)
}

func (s *ProjectTestSuite) TestProjectConfigFileDirectory() {
	ld := NewProjectLoader(s.repositoryLoader, s.recipeLoader, s.repositorySrc)
	file, err := ld.ConfigFile("testdata/project/config_file_directory")
	s.Error(err)
	s.Equal("open testdata/project/config_file_directory/.manala.yaml: is a directory", err.Error())
	s.Nil(file)
}

func (s *ProjectTestSuite) TestProjectLoad() {
	ld := NewProjectLoader(s.repositoryLoader, s.recipeLoader, s.repositorySrc)
	prj, err := ld.Load("testdata/project/load")
	s.NoError(err)
	s.Implements((*models.ProjectInterface)(nil), prj)
	s.Equal("testdata/project/load", prj.Dir())
	s.Equal("foo", prj.Recipe().Name())
}

func (s *ProjectTestSuite) TestProjectLoadNotFound() {
	ld := NewProjectLoader(s.repositoryLoader, s.recipeLoader, s.repositorySrc)
	prj, err := ld.Load("testdata/project/load_not_found")
	s.Error(err)
	s.Equal("open testdata/project/load_not_found/.manala.yaml: no such file or directory", err.Error())
	s.Nil(prj)
}

func (s *ProjectTestSuite) TestProjectLoadEmpty() {
	ld := NewProjectLoader(s.repositoryLoader, s.recipeLoader, s.repositorySrc)
	prj, err := ld.Load("testdata/project/load_empty")
	s.Error(err)
	s.Equal("empty project config \"testdata/project/load_empty/.manala.yaml\"", err.Error())
	s.Nil(prj)
}

func (s *ProjectTestSuite) TestProjectLoadIncorrect() {
	ld := NewProjectLoader(s.repositoryLoader, s.recipeLoader, s.repositorySrc)
	prj, err := ld.Load("testdata/project/load_incorrect")
	s.Error(err)
	s.Equal("invalid project config \"testdata/project/load_incorrect/.manala.yaml\" (yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `foo` into map[string]interface {})", err.Error())
	s.Nil(prj)
}

func (s *ProjectTestSuite) TestProjectLoadNoRecipe() {
	ld := NewProjectLoader(s.repositoryLoader, s.recipeLoader, s.repositorySrc)
	prj, err := ld.Load("testdata/project/load_no_recipe")
	s.Error(err)
	s.Equal("Key: 'projectConfig.Recipe' Error:Field validation for 'Recipe' failed on the 'required' tag", err.Error())
	s.Nil(prj)
}

func (s *ProjectTestSuite) TestProjectLoadRepository() {
	ld := NewProjectLoader(s.repositoryLoader, s.recipeLoader, s.repositorySrc)
	prj, err := ld.Load("testdata/project/load_repository")
	s.NoError(err)
	s.Equal("bar", prj.Recipe().Name())
}

func (s *ProjectTestSuite) TestProjectLoadVars() {
	ld := NewProjectLoader(s.repositoryLoader, s.recipeLoader, s.repositorySrc)
	prj, err := ld.Load("testdata/project/load_vars")
	s.NoError(err)
	s.Equal(
		map[string]interface{}{
			"foo":  map[string]interface{}{"foo": "bar", "bar": "baz", "baz": []interface{}{}},
			"bar":  map[string]interface{}{"bar": "baz"},
			"baz":  "qux",
			"qux":  map[string]interface{}{"qux": "quux"},
			"quux": map[string]interface{}{"qux": "quux", "quux": "corge"},
		},
		prj.Vars(),
	)
}
