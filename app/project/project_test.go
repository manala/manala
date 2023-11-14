package project

import (
	"bytes"
	_ "embed"
	"github.com/stretchr/testify/suite"
	"manala/app"
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
	repositoryMock := &app.RepositoryMock{}
	repositoryMock.
		On("Url").Return("repository")

	recipeMock := &app.RecipeMock{}
	recipeMock.
		On("Name").Return("recipe").
		On("Vars").Return(map[string]any{
		"foo": "recipe",
		"bar": "recipe",
	}).
		On("Repository").Return(repositoryMock).
		On("Template").Return(template.NewTemplate())

	dir := filepath.Join("dir")

	manifestMock := &app.ProjectManifestMock{}
	manifestMock.
		On("Vars").Return(map[string]any{
		"bar": "project",
		"baz": "project",
	})

	project := New(
		dir,
		manifestMock,
		recipeMock,
	)

	s.Equal("dir", project.Dir())
	s.Equal(recipeMock, project.Recipe())

	s.Run("Vars", func() {
		s.Equal(map[string]any{
			"foo": "recipe",
			"bar": "project",
			"baz": "project",
		}, project.Vars())
	})

	s.Run("Template", func() {
		template := project.Template()

		out := &bytes.Buffer{}
		err := template.
			WithDefaultContent(`{{ .Vars | toYaml }}`).
			WriteTo(out)

		s.NoError(err)

		heredoc.Equal(s.Assert(), `
			bar: project
			baz: project
			foo: recipe`,
			out.String(),
		)
	})
}
