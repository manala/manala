package internal

import (
	"github.com/stretchr/testify/suite"
	"io"
	internalLog "manala/internal/log"
	"path/filepath"
	"testing"
)

type ProjectLoadersSuite struct{ suite.Suite }

func TestProjectLoadersSuite(t *testing.T) {
	suite.Run(t, new(ProjectLoadersSuite))
}

var projectLoadersTestPath = filepath.Join("testdata", "project_loaders")

func (s *ProjectLoadersSuite) TestDir() {
	log := internalLog.New(io.Discard)
	loader := &ProjectDirLoader{
		Log:              log,
		RepositoryLoader: &RepositoryDirLoader{Log: log},
	}

	s.Run("LoadProjectManifest Not Found", func() {
		projectManifest, err := loader.LoadProjectManifest(filepath.Join(projectLoadersTestPath, "project_not_found"))

		var _err *NotFoundProjectManifestError
		s.ErrorAs(err, &_err)

		s.ErrorAs(err, &internalError)
		s.Equal("project manifest not found", internalError.Message)
		s.Nil(projectManifest)
	})

	s.Run("LoadProjectManifest Wrong", func() {
		projectManifest, err := loader.LoadProjectManifest(filepath.Join(projectLoadersTestPath, "project_wrong"))

		s.ErrorAs(err, &internalError)
		s.Equal("wrong project manifest", internalError.Message)
		s.Nil(projectManifest)
	})

	s.Run("LoadProjectManifest", func() {
		projectManifest, err := loader.LoadProjectManifest(filepath.Join(projectLoadersTestPath, "project"))

		s.NoError(err)
		s.Equal("recipe", projectManifest.Recipe)
		s.Equal("repository", projectManifest.Repository)
		s.Equal(map[string]interface{}{"foo": "bar"}, projectManifest.Vars)
	})

	s.Run("LoadProject", func() {
		project, err := loader.LoadProject(
			filepath.Join(projectLoadersTestPath, "project"),
			filepath.Join(projectLoadersTestPath, "repository"),
			"recipe",
		)

		s.NoError(err)
		s.Equal(filepath.Join(projectLoadersTestPath, "project"), project.Path())
		s.Equal(filepath.Join(projectLoadersTestPath, "repository"), project.Recipe().Repository().Path())
		s.Equal(filepath.Join(projectLoadersTestPath, "repository", "recipe"), project.Recipe().Path())
	})
}
