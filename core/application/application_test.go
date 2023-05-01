package application

import (
	"bytes"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/suite"
	"manala/app/interfaces"
	"manala/app/mocks"
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
	suite.Run(t, new(ApplicationSuite))
}

func (s *ApplicationSuite) SetupTest() {
	s.goldie = goldie.New(s.T())
}

func (s *ApplicationSuite) TestCreateProject() {
	path := internalTesting.DataPath(s)
	projDir := filepath.Join(path, "project")
	repoUrl := filepath.Join(path, "repository")

	_ = os.RemoveAll(projDir)

	stderr := &bytes.Buffer{}

	confMock := mocks.MockConfig()
	confMock.
		On("Fields").Return(map[string]interface{}{}).
		On("CacheDir").Return("").
		On("Repository").Return("")

	app := NewApplication(
		confMock,
		internalLog.New(stderr),
		WithRepositoryUrl(repoUrl),
	)

	proj, err := app.CreateProject(
		projDir,
		// Recipe selector
		func(recWalker func(walker func(rec interfaces.Recipe) error) error) (interfaces.Recipe, error) {
			var rec interfaces.Recipe
			_ = recWalker(func(_rec interfaces.Recipe) error {
				rec = _rec
				return nil
			})
			return rec, nil
		},
		// Options selector
		func(rec interfaces.Recipe, options []interfaces.RecipeOption) error {
			// String float int
			_ = options[0].Set("3.0")
			// String asterisk
			_ = options[1].Set("*")
			return nil
		},
	)

	s.NotNil(proj)
	s.NoError(err)

	s.goldie.Assert(s.T(), internalTesting.Path(s, "stderr"), stderr.Bytes())
	manContent, _ := os.ReadFile(filepath.Join(projDir, ".manala.yaml"))
	s.goldie.Assert(s.T(), internalTesting.Path(s, "manifest"), manContent)
}
