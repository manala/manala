package manifest

import (
	"bytes"
	_ "embed"
	"manala/app"
	"manala/internal/template"
	"manala/internal/testing/heredoc"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
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
		On("URL").Return("repository")

	recipeMock := &app.RecipeMock{}
	recipeMock.
		On("Name").Return("recipe").
		On("Description").Return("description").
		On("Icon").Return("icon").
		On("Vars").Return(map[string]any{
		"foo": "recipe",
		"bar": "recipe",
	}).
		On("Repository").Return(repositoryMock).
		On("Template").Return(template.NewTemplate())

	dir := "dir"

	manifest := &Manifest{
		vars: map[string]any{
			"bar": "project",
			"baz": "project",
		},
	}

	project := NewProject(
		dir,
		manifest,
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

		s.Require().NoError(err)
		heredoc.Equal(s.T(), `
			bar: project
			baz: project
			foo: recipe`, out)
	})

	s.Run("Watches", func() {
		watches, err := project.Watches()

		s.Equal([]string{
			filepath.Join(dir, filename),
		}, watches)
		s.NoError(err)
	})
}
