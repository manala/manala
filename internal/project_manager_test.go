package internal

import (
	"github.com/stretchr/testify/suite"
	"io"
	internalLog "manala/internal/log"
	"path/filepath"
	"testing"
)

type ProjectManagerSuite struct{ suite.Suite }

func TestProjectManagerSuite(t *testing.T) {
	suite.Run(t, new(ProjectManagerSuite))
}

var projectManagerTestPath = filepath.Join("testdata", "project_manager")

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
		projectManifest, err := manager.LoadProjectManifest(filepath.Join(projectManagerTestPath, "project"))

		s.NoError(err)
		s.Equal("recipe", projectManifest.Recipe)
		s.Equal("repository", projectManifest.Repository)
		s.Equal(map[string]interface{}{"foo": "bar"}, projectManifest.Vars)
	})

	s.Run("LoadProject", func() {
		project, err := manager.LoadProject(
			filepath.Join(projectManagerTestPath, "project"),
			filepath.Join(projectManagerTestPath, "repository"),
			"recipe",
		)

		s.NoError(err)
		s.Equal(filepath.Join(projectManagerTestPath, "project"), project.Path())
		s.Equal(filepath.Join(projectManagerTestPath, "repository"), project.Recipe().Repository().Path())
		s.Equal(filepath.Join(projectManagerTestPath, "repository", "recipe"), project.Recipe().Path())
	})
}
