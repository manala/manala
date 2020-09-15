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
	s.repositoryLoader = NewRepositoryLoader(
		cacheDir,
		"testdata/project/_repository_default",
	)
	s.recipeLoader = NewRecipeLoader()
}

/*******************/
/* Project - Tests */
/*******************/

func (s *ProjectTestSuite) TestProject() {
	ld := NewProjectLoader(s.repositoryLoader, s.recipeLoader, "", "")
	s.Implements((*ProjectLoaderInterface)(nil), ld)
}

func (s *ProjectTestSuite) TestProjectFind() {
	for _, t := range []struct {
		test        string
		dir         string
		prjFileName string
	}{
		{
			test:        "Default",
			dir:         "testdata/project/find/default",
			prjFileName: "testdata/project/find/default/.manala.yaml",
		},
		{
			test: "Not found",
			dir:  "testdata/project/find/not_found",
		},
	} {
		s.Run(t.test, func() {
			ld := NewProjectLoader(s.repositoryLoader, s.recipeLoader, "", "")
			prjFile, err := ld.Find(t.dir, false)
			s.NoError(err)
			if t.prjFileName != "" {
				s.NotNil(prjFile)
				s.Equal(t.prjFileName, prjFile.Name())
			} else {
				s.Nil(prjFile)
			}
		})
	}
}

func (s *ProjectTestSuite) TestProjectFindTraverse() {
	for _, t := range []struct {
		test        string
		dir         string
		prjFileName string
	}{
		{
			test:        "Default",
			dir:         "testdata/project/find_traverse/default",
			prjFileName: "testdata/project/find_traverse/default/.manala.yaml",
		},
		{
			test: "Not found",
			dir:  "testdata/project/find_traverse/not_found",
		},
		{
			test:        "Level one",
			dir:         "testdata/project/find_traverse/traverse/level",
			prjFileName: "testdata/project/find_traverse/traverse/.manala.yaml",
		},
	} {
		s.Run(t.test, func() {
			ld := NewProjectLoader(s.repositoryLoader, s.recipeLoader, "", "")
			prjFile, err := ld.Find(t.dir, true)
			s.NoError(err)
			if t.prjFileName != "" {
				s.NotNil(prjFile)
				s.Equal(t.prjFileName, prjFile.Name())
			} else {
				s.Nil(prjFile)
			}
		})
	}
}

func (s *ProjectTestSuite) TestProjectLoad() {
	for _, t := range []struct {
		test               string
		forceRepositorySrc string
		forceRecipe        string
		recipeName         string
		recipeDescription  string
	}{
		{
			test:               "Default",
			forceRepositorySrc: "",
			forceRecipe:        "",
			recipeName:         "foo",
			recipeDescription:  "Default foo",
		},
		{
			test:               "Force repository",
			forceRepositorySrc: "testdata/project/_repository_force",
			forceRecipe:        "",
			recipeName:         "foo",
			recipeDescription:  "Force foo",
		},
		{
			test:               "Force recipe",
			forceRepositorySrc: "",
			forceRecipe:        "bar",
			recipeName:         "bar",
			recipeDescription:  "Default bar",
		},
		{
			test:               "Force repository force recipe",
			forceRepositorySrc: "testdata/project/_repository_force",
			forceRecipe:        "bar",
			recipeName:         "bar",
			recipeDescription:  "Force bar",
		},
	} {
		s.Run(t.test, func() {
			ld := NewProjectLoader(s.repositoryLoader, s.recipeLoader, t.forceRepositorySrc, t.forceRecipe)
			prjFile, err := ld.Find("testdata/project/load", false)
			s.NoError(err)
			prj, err := ld.Load(prjFile)
			s.NoError(err)
			s.Implements((*models.ProjectInterface)(nil), prj)
			s.Equal("testdata/project/load", prj.Dir())
			s.Equal(t.recipeName, prj.Recipe().Name())
			s.Equal(t.recipeDescription, prj.Recipe().Description())
		})
	}
}

func (s *ProjectTestSuite) TestProjectLoadEmpty() {
	ld := NewProjectLoader(s.repositoryLoader, s.recipeLoader, "", "")
	prjFile, err := ld.Find("testdata/project/load_empty", false)
	s.NoError(err)
	prj, err := ld.Load(prjFile)
	s.Error(err)
	s.Equal("empty project config \"testdata/project/load_empty/.manala.yaml\"", err.Error())
	s.Nil(prj)
}

func (s *ProjectTestSuite) TestProjectLoadIncorrect() {
	ld := NewProjectLoader(s.repositoryLoader, s.recipeLoader, "", "")
	prjFile, err := ld.Find("testdata/project/load_incorrect", false)
	s.NoError(err)
	prj, err := ld.Load(prjFile)
	s.Error(err)
	s.Equal("invalid project config \"testdata/project/load_incorrect/.manala.yaml\" (yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `foo` into map[string]interface {})", err.Error())
	s.Nil(prj)
}

func (s *ProjectTestSuite) TestProjectLoadNoRecipe() {
	ld := NewProjectLoader(s.repositoryLoader, s.recipeLoader, "", "")
	prjFile, err := ld.Find("testdata/project/load_no_recipe", false)
	s.NoError(err)
	prj, err := ld.Load(prjFile)
	s.Error(err)
	s.Equal("Key: 'projectConfig.Recipe' Error:Field validation for 'Recipe' failed on the 'required' tag", err.Error())
	s.Nil(prj)
}

func (s *ProjectTestSuite) TestProjectLoadRepository() {
	for _, t := range []struct {
		test               string
		forceRepositorySrc string
		forceRecipe        string
		recipeName         string
		recipeDescription  string
	}{
		{
			test:               "Default",
			forceRepositorySrc: "",
			forceRecipe:        "",
			recipeName:         "foo",
			recipeDescription:  "Custom foo",
		},
		{
			test:               "Force repository",
			forceRepositorySrc: "testdata/project/_repository_force",
			forceRecipe:        "",
			recipeName:         "foo",
			recipeDescription:  "Force foo",
		},
		{
			test:               "Force recipe",
			forceRepositorySrc: "",
			forceRecipe:        "bar",
			recipeName:         "bar",
			recipeDescription:  "Custom bar",
		},
		{
			test:               "Force repository force recipe",
			forceRepositorySrc: "testdata/project/_repository_force",
			forceRecipe:        "bar",
			recipeName:         "bar",
			recipeDescription:  "Force bar",
		},
	} {
		s.Run(t.test, func() {
			ld := NewProjectLoader(s.repositoryLoader, s.recipeLoader, t.forceRepositorySrc, t.forceRecipe)
			prjFile, err := ld.Find("testdata/project/load_repository", false)
			s.NoError(err)
			prj, err := ld.Load(prjFile)
			s.NoError(err)
			s.Implements((*models.ProjectInterface)(nil), prj)
			s.Equal("testdata/project/load_repository", prj.Dir())
			s.Equal(t.recipeName, prj.Recipe().Name())
			s.Equal(t.recipeDescription, prj.Recipe().Description())
		})
	}
}

func (s *ProjectTestSuite) TestProjectLoadVars() {
	ld := NewProjectLoader(s.repositoryLoader, s.recipeLoader, "", "")
	prjFile, err := ld.Find("testdata/project/load_vars", false)
	s.NoError(err)
	prj, err := ld.Load(prjFile)
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
