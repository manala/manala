package internal

import (
	"github.com/stretchr/testify/suite"
	"io"
	internalLog "manala/internal/log"
	internalTesting "manala/internal/testing"
	"testing"
)

type ProjectLoadersSuite struct{ suite.Suite }

func TestProjectLoadersSuite(t *testing.T) {
	suite.Run(t, new(ProjectLoadersSuite))
}

func (s *ProjectLoadersSuite) TestDir() {
	log := internalLog.New(io.Discard)
	loader := &ProjectDirLoader{
		Log:              log,
		RepositoryLoader: &RepositoryDirLoader{Log: log},
	}

	s.Run("LoadProjectManifest Not Found", func() {
		projectManifest, err := loader.LoadProjectManifest(internalTesting.DataPath(s, "project"))

		var _err *NotFoundProjectManifestError
		s.ErrorAs(err, &_err)

		s.ErrorAs(err, &internalError)
		s.Equal("project manifest not found", internalError.Message)
		s.Nil(projectManifest)
	})

	s.Run("LoadProjectManifest Wrong", func() {
		projectManifest, err := loader.LoadProjectManifest(internalTesting.DataPath(s, "project"))

		s.ErrorAs(err, &internalError)
		s.Equal("wrong project manifest", internalError.Message)
		s.Nil(projectManifest)
	})

	s.Run("LoadProjectManifest", func() {
		projectManifest, err := loader.LoadProjectManifest(internalTesting.DataPath(s, "project"))

		s.NoError(err)
		s.Equal("recipe", projectManifest.Recipe)
		s.Equal("repository", projectManifest.Repository)
		s.Equal(map[string]interface{}{"foo": "bar"}, projectManifest.Vars)
	})

	s.Run("LoadProject", func() {
		project, err := loader.LoadProject(
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
