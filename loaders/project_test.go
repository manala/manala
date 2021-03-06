package loaders

import (
	"github.com/stretchr/testify/suite"
	"manala/config"
	"manala/fs"
	"manala/logger"
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
	ld ProjectLoaderInterface
}

func TestProjectTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(ProjectTestSuite))
}

func (s *ProjectTestSuite) SetupTest() {
	cacheDir := "testdata/project/.cache"
	_ = os.RemoveAll(cacheDir)
	_ = os.Mkdir(cacheDir, 0755)

	conf := config.New(
		config.WithMainRepository("testdata/project/_repository_default"),
		config.WithCacheDir(cacheDir),
	)

	log := logger.New()

	fsManager := fs.NewManager()
	modelFsManager := models.NewFsManager(fsManager)

	repositoryLoader := NewRepositoryLoader(log, conf)
	recipeLoader := NewRecipeLoader(log, modelFsManager)

	s.ld = NewProjectLoader(log, conf, repositoryLoader, recipeLoader)
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
		test                 string
		withRepositorySource string
		withRecipeName       string
		recipeName           string
		recipeDescription    string
	}{
		{
			test:                 "Default",
			withRepositorySource: "",
			withRecipeName:       "",
			recipeName:           "foo",
			recipeDescription:    "Default foo",
		},
		{
			test:                 "With repository",
			withRepositorySource: "testdata/project/_repository_with",
			withRecipeName:       "",
			recipeName:           "foo",
			recipeDescription:    "With foo",
		},
		{
			test:                 "With recipe",
			withRepositorySource: "",
			withRecipeName:       "bar",
			recipeName:           "bar",
			recipeDescription:    "Default bar",
		},
		{
			test:                 "With repository with recipe",
			withRepositorySource: "testdata/project/_repository_with",
			withRecipeName:       "bar",
			recipeName:           "bar",
			recipeDescription:    "With bar",
		},
	} {
		s.Run(t.test, func() {
			prjManifest, err := s.ld.Find("testdata/project/load", false)
			s.NoError(err)
			prj, err := s.ld.Load(prjManifest, t.withRepositorySource, t.withRecipeName)
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
	prj, err := s.ld.Load(prjManifest, "", "")
	s.Error(err)
	s.Equal("empty project manifest \""+filepath.Join("testdata", "project", "load_empty", ".manala.yaml")+"\"", err.Error())
	s.Nil(prj)
}

func (s *ProjectTestSuite) TestProjectLoadIncorrect() {
	prjManifest, err := s.ld.Find("testdata/project/load_incorrect", false)
	s.NoError(err)
	prj, err := s.ld.Load(prjManifest, "", "")
	s.Error(err)
	s.Equal("incorrect project manifest \""+filepath.Join("testdata", "project", "load_incorrect", ".manala.yaml")+"\" \x1b[91m[1:1] string was used where mapping is expected\x1b[0m\n>  1 | \x1b[92mfoo\x1b[0m\n       ^\n", err.Error())
	s.Nil(prj)
}

func (s *ProjectTestSuite) TestProjectLoadNoRecipe() {
	prjManifest, err := s.ld.Find("testdata/project/load_no_recipe", false)
	s.NoError(err)
	prj, err := s.ld.Load(prjManifest, "", "")
	s.Error(err)
	s.Equal("Key: 'projectConfig.Recipe' Error:Field validation for 'Recipe' failed on the 'required' tag", err.Error())
	s.Nil(prj)
}

func (s *ProjectTestSuite) TestProjectLoadRepository() {
	for _, t := range []struct {
		test                 string
		withRepositorySource string
		withRecipeName       string
		recipeName           string
		recipeDescription    string
	}{
		{
			test:                 "Default",
			withRepositorySource: "",
			withRecipeName:       "",
			recipeName:           "foo",
			recipeDescription:    "Custom foo",
		},
		{
			test:                 "With repository",
			withRepositorySource: "testdata/project/_repository_with",
			withRecipeName:       "",
			recipeName:           "foo",
			recipeDescription:    "With foo",
		},
		{
			test:                 "With recipe",
			withRepositorySource: "",
			withRecipeName:       "bar",
			recipeName:           "bar",
			recipeDescription:    "Custom bar",
		},
		{
			test:                 "With repository with recipe",
			withRepositorySource: "testdata/project/_repository_with",
			withRecipeName:       "bar",
			recipeName:           "bar",
			recipeDescription:    "With bar",
		},
	} {
		s.Run(t.test, func() {
			prjManifest, err := s.ld.Find("testdata/project/load_repository", false)
			s.NoError(err)
			prj, err := s.ld.Load(prjManifest, t.withRepositorySource, t.withRecipeName)
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
	prj, err := s.ld.Load(prjManifest, "", "")
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
