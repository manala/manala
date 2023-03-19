package project

import (
	"bytes"
	_ "embed"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/suite"
	"manala/app/mocks"
	internalTemplate "manala/internal/template"
	internalTesting "manala/internal/testing"
	"path/filepath"
	"testing"
)

type ProjectSuite struct {
	suite.Suite
	goldie *goldie.Goldie
}

func TestProjectSuite(t *testing.T) {
	suite.Run(t, new(ProjectSuite))
}

func (s *ProjectSuite) SetupTest() {
	s.goldie = goldie.New(s.T())
}

func (s *ProjectSuite) Test() {
	repoMock := mocks.MockRepository()
	repoMock.
		On("Url").Return("repository")

	recMock := mocks.MockRecipe()
	recMock.
		On("Name").Return("recipe").
		On("Vars").Return(map[string]interface{}{
		"foo": "recipe",
		"bar": "recipe",
	}).
		On("Repository").Return(repoMock).
		On("Template").Return(internalTemplate.NewTemplate())

	projDir := filepath.Join("dir")

	projManifestMock := mocks.MockProjectManifest()
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
		s.goldie.Assert(s.T(), internalTesting.Path(s, "template"), out.Bytes())
	})
}
