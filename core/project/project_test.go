package project

import (
	"bytes"
	_ "embed"
	"github.com/stretchr/testify/suite"
	"manala/app/mocks"
	"manala/internal/template"
	"manala/internal/testing/heredoc"
	"path/filepath"
	"testing"
)

type ProjectSuite struct {
	suite.Suite
}

func TestProjectSuite(t *testing.T) {
	suite.Run(t, new(ProjectSuite))
}

func (s *ProjectSuite) Test() {
	repoMock := &mocks.RepositoryMock{}
	repoMock.
		On("Url").Return("repository")

	recMock := &mocks.RecipeMock{}
	recMock.
		On("Name").Return("recipe").
		On("Vars").Return(map[string]interface{}{
		"foo": "recipe",
		"bar": "recipe",
	}).
		On("Repository").Return(repoMock).
		On("Template").Return(template.NewTemplate())

	projDir := filepath.Join("dir")

	projManifestMock := &mocks.ProjectManifestMock{}
	projManifestMock.
		On("Vars").Return(map[string]interface{}{
		"bar": "project",
		"baz": "project",
	})

	proj := NewProject(
		projDir,
		projManifestMock,
		recMock,
	)

	s.Equal("dir", proj.Dir())
	s.Equal(recMock, proj.Recipe())

	s.Run("Vars", func() {
		s.Equal(map[string]interface{}{
			"foo": "recipe",
			"bar": "project",
			"baz": "project",
		}, proj.Vars())
	})

	s.Run("Template", func() {
		template := proj.Template()

		out := &bytes.Buffer{}
		err := template.
			WithDefaultContent(`{{ .Vars | toYaml }}`).
			WriteTo(out)

		s.NoError(err)

		s.Equal(heredoc.Doc(`
			bar: project
			baz: project
			foo: recipe`,
		), out.String())
	})
}
