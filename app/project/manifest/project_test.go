package manifest_test

import (
	_ "embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/project/manifest"

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
		On("Repository").Return(repositoryMock)

	m := manifest.New()

	dir := filepath.FromSlash("testdata/ProjectSuite/Test")
	content, _ := os.ReadFile(filepath.Join(dir, "manifest.yaml"))

	err := m.Unmarshal(content)
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

	s.Run("Watches", func() {
		watches, err := project.Watches()

		s.Require().NoError(err)
		s.Equal([]string{
			filepath.Join(dir, ".manala.yaml"),
		}, watches)
	})
}
