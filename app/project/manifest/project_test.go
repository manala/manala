package manifest_test

import (
	"bytes"
	_ "embed"
	"os"
	"path/filepath"
	"testing"

	"manala/app"
	"manala/app/project/manifest"
	"manala/internal/template"
	"manala/internal/testing/heredoc"

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

	dir := filepath.FromSlash("testdata/ProjectSuite/Test")

	m := manifest.New()

	mFile, _ := os.Open(filepath.Join(dir, "manifest.yaml"))
	_, err := m.ReadFrom(mFile)

	s.Require().NoError(err)

	project := manifest.NewProject(
		dir,
		m,
		recipeMock,
	)

	s.Equal(dir, project.Dir())
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

		s.Require().NoError(err)
		s.Equal([]string{
			filepath.Join(dir, ".manala.yaml"),
		}, watches)
	})
}
