package loaders

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/discard"
	"github.com/stretchr/testify/suite"
	"manala/fs"
	"manala/models"
	"os"
	"path/filepath"
	"testing"
)

/*******************/
/* Project - Suite */
/*******************/

type ProjectTestSuite struct {
	suite.Suite
	cacheDir string
	ld       ProjectLoaderInterface
}

func TestProjectTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(ProjectTestSuite))
}

func (s *ProjectTestSuite) SetupTest() {
	s.cacheDir = "testdata/project/.cache"
	_ = os.RemoveAll(s.cacheDir)
	_ = os.Mkdir(s.cacheDir, 0755)

	logger := &log.Logger{
		Handler: discard.Default,
	}

	fsManager := fs.NewManager()
	modelFsManager := models.NewFsManager(fsManager)

	repositoryLoader := NewRepositoryLoader(logger)
	recipeLoader := NewRecipeLoader(logger, modelFsManager)

	s.ld = NewProjectLoader(logger, repositoryLoader, recipeLoader)
}

/*******************/
/* Project - Tests */
/*******************/

func (s *ProjectTestSuite) TestProject() {
	s.Implements((*ProjectLoaderInterface)(nil), s.ld)
}

func (s *ProjectTestSuite) TestProjectFind() {
	for _, t := range []struct {
		test            string
		dir             string
		prjManifestName string
	}{
		{
			test:            "Default",
			dir:             "testdata/project/find/default",
			prjManifestName: filepath.Join("testdata", "project", "find", "default", ".manala.yaml"),
		},
		{
			test: "Not found",
			dir:  "testdata/project/find/not_found",
		},
	} {
		s.Run(t.test, func() {
			prjManifest, err := s.ld.Find(t.dir, false)
			s.NoError(err)
			if t.prjManifestName != "" {
				s.NotNil(prjManifest)
				s.Equal(t.prjManifestName, prjManifest.Name())
			} else {
				s.Nil(prjManifest)
			}
		})
	}
}

func (s *ProjectTestSuite) TestProjectFindTraverse() {
	for _, t := range []struct {
		test            string
		dir             string
		prjManifestName string
	}{
		{
			test:            "Default",
			dir:             "testdata/project/find_traverse/default",
			prjManifestName: filepath.Join("testdata", "project", "find_traverse", "default", ".manala.yaml"),
		},
		{
			test: "Not found",
			dir:  "testdata/project/find_traverse/not_found",
		},
		{
			test:            "Level one",
			dir:             "testdata/project/find_traverse/traverse/level",
			prjManifestName: filepath.Join("testdata", "project", "find_traverse", "traverse", ".manala.yaml"),
		},
	} {
		s.Run(t.test, func() {
			prjManifest, err := s.ld.Find(t.dir, true)
			s.NoError(err)
			if t.prjManifestName != "" {
				s.NotNil(prjManifest)
				s.Equal(t.prjManifestName, prjManifest.Name())
			} else {
				s.Nil(prjManifest)
			}
		})
	}
}

func (s *ProjectTestSuite) TestProjectLoad() {
	for _, t := range []struct {
		test              string
		defaultRepository string
		withRecipeName    string
		recipeName        string
		recipeDescription string
	}{
		{
			test:              "With default repository",
			defaultRepository: "testdata/project/_repository_with",
			withRecipeName:    "",
			recipeName:        "foo",
			recipeDescription: "With foo",
		},
		{
			test:              "With default repository and recipe",
			defaultRepository: "testdata/project/_repository_with",
			withRecipeName:    "bar",
			recipeName:        "bar",
			recipeDescription: "With bar",
		},
	} {
		s.Run(t.test, func() {
			prjManifest, err := s.ld.Find("testdata/project/load", false)
			s.NoError(err)
			prj, err := s.ld.Load(prjManifest, t.defaultRepository, t.withRecipeName, s.cacheDir)
			s.NoError(err)
			s.Implements((*models.ProjectInterface)(nil), prj)
			s.Equal(t.recipeName, prj.Recipe().Name())
			s.Equal(t.recipeDescription, prj.Recipe().Description())
		})
	}
}

func (s *ProjectTestSuite) TestProjectLoadEmpty() {
	prjManifest, err := s.ld.Find("testdata/project/load_empty", false)
	s.NoError(err)
	prj, err := s.ld.Load(prjManifest, "", "", s.cacheDir)
	s.Error(err)
	s.Equal("empty project manifest \""+filepath.Join("testdata", "project", "load_empty", ".manala.yaml")+"\"", err.Error())
	s.Nil(prj)
}

func (s *ProjectTestSuite) TestProjectLoadIncorrect() {
	prjManifest, err := s.ld.Find("testdata/project/load_incorrect", false)
	s.NoError(err)
	prj, err := s.ld.Load(prjManifest, "", "", s.cacheDir)
	s.Error(err)
	s.Equal("incorrect project manifest \""+filepath.Join("testdata", "project", "load_incorrect", ".manala.yaml")+"\" \x1b[91m[1:1] string was used where mapping is expected\x1b[0m\n>  1 | \x1b[92mfoo\x1b[0m\n       ^\n", err.Error())
	s.Nil(prj)
}

func (s *ProjectTestSuite) TestProjectLoadNoRecipe() {
	prjManifest, err := s.ld.Find("testdata/project/load_no_recipe", false)
	s.NoError(err)
	prj, err := s.ld.Load(prjManifest, "", "", s.cacheDir)
	s.Error(err)
	s.Equal("Key: 'projectConfig.Recipe' Error:Field validation for 'Recipe' failed on the 'required' tag", err.Error())
	s.Nil(prj)
}

func (s *ProjectTestSuite) TestProjectLoadRepository() {
	for _, t := range []struct {
		test              string
		defaultRepository string
		withRecipeName    string
		recipeName        string
		recipeDescription string
	}{
		{
			test:              "With default repository",
			defaultRepository: "testdata/project/_repository_with",
			withRecipeName:    "",
			recipeName:        "foo",
			recipeDescription: "Custom foo",
		},
		{
			test:              "With default repository and recipe",
			defaultRepository: "testdata/project/_repository_with",
			withRecipeName:    "bar",
			recipeName:        "bar",
			recipeDescription: "Custom bar",
		},
	} {
		s.Run(t.test, func() {
			prjManifest, err := s.ld.Find("testdata/project/load_repository", false)
			s.NoError(err)
			prj, err := s.ld.Load(prjManifest, t.defaultRepository, t.withRecipeName, s.cacheDir)
			s.NoError(err)
			s.Implements((*models.ProjectInterface)(nil), prj)
			s.Equal(t.recipeName, prj.Recipe().Name())
			s.Equal(t.recipeDescription, prj.Recipe().Description())
		})
	}
}

func (s *ProjectTestSuite) TestProjectLoadVars() {
	prjManifest, err := s.ld.Find("testdata/project/load_vars", false)
	s.NoError(err)
	prj, err := s.ld.Load(prjManifest, "testdata/project/_repository_default", "", s.cacheDir)
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
