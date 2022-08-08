package internal

import (
	"github.com/stretchr/testify/suite"
	"io"
	internalLog "manala/internal/log"
	internalTesting "manala/internal/testing"
	"testing"
)

type ProjectManagerSuite struct{ suite.Suite }

func TestProjectManagerSuite(t *testing.T) {
	suite.Run(t, new(ProjectManagerSuite))
}

func (s *ProjectManagerSuite) Test() {
	log := internalLog.New(io.Discard)
	repositoryLoader := &RepositoryDirLoader{Log: log}
	projectLoader := &ProjectDirLoader{
		Log:              log,
		RepositoryLoader: repositoryLoader,
	}

	manager := NewProjectManager(
		log,
		repositoryLoader,
		projectLoader,
	)

	s.Run("LoadProjectManifest", func() {
		projectManifest, err := manager.LoadProjectManifest(internalTesting.DataPath(s, "project"))

		s.NoError(err)
		s.Equal("recipe", projectManifest.Recipe)
		s.Equal("repository", projectManifest.Repository)
		s.Equal(map[string]interface{}{"foo": "bar"}, projectManifest.Vars)
	})

	s.Run("LoadProject", func() {
		project, err := manager.LoadProject(
			internalTesting.DataPath(s, "project"),
			internalTesting.DataPath(s, "repository"),
			"recipe",
		)

		s.NoError(err)
		s.Equal(internalTesting.DataPath(s, "project"), project.Path())
		s.Equal(internalTesting.DataPath(s, "repository"), project.Recipe().Repository().Path())
		s.Equal(internalTesting.DataPath(s, "repository", "recipe"), project.Recipe().Path())
	})
}
