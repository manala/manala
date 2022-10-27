package application

import (
	"bytes"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/suite"
	"manala/core"
	internalConfig "manala/internal/config"
	internalLog "manala/internal/log"
	internalTesting "manala/internal/testing"
	"os"
	"path/filepath"
	"testing"
)

type ApplicationSuite struct {
	suite.Suite
	goldie *goldie.Goldie
}

func TestApplicationSuite(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	suite.Run(t, new(ApplicationSuite))
}

func (s *ApplicationSuite) SetupTest() {
	s.goldie = goldie.New(s.T())
}

func (s *ApplicationSuite) TestCreateProject() {
	path := internalTesting.DataPath(s)
	projPath := filepath.Join(path, "project")
	repoPath := filepath.Join(path, "repository")

	_ = os.RemoveAll(projPath)

	stderr := &bytes.Buffer{}

	app := NewApplication(
		internalConfig.New(),
		internalLog.New(stderr),
	)

	repo, _ := app.Repository(repoPath)

	proj, err := app.CreateProject(projPath, repo,
		// Recipe selector
		func(recWalker core.RecipeWalker) (core.Recipe, error) {
			var rec core.Recipe
			recWalker.WalkRecipes(func(_rec core.Recipe) {
				rec = _rec
			})
			return rec, nil
		},
		// Options selector
		func(rec core.Recipe, options []core.RecipeOption) error {
			// String float int
			options[0].Set("3.0")
			// String asterisk
			options[1].Set("*")
			return nil
		},
	)

	s.NotNil(proj)
	s.Nil(err)

	s.goldie.Assert(s.T(), internalTesting.Path(s, "stderr"), stderr.Bytes())
	manifestContent, _ := os.ReadFile(filepath.Join(projPath, ".manala.yaml"))
	s.goldie.Assert(s.T(), internalTesting.Path(s, "manifest"), manifestContent)
}
