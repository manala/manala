package project

import (
	"bytes"
	_ "embed"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/suite"
	"manala/core"
	internalTesting "manala/internal/testing"
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
	repo := core.NewRepositoryMock().
		WithPath("repository")

	rec := core.NewRecipeMock().
		WithName("recipe").
		WithVars(map[string]interface{}{
			"foo": "recipe",
			"bar": "recipe",
		}).
		WithRepository(repo)

	projManifest := core.NewProjectManifestMock().
		WithDir("dir").
		WithVars(map[string]interface{}{
			"bar": "project",
			"baz": "project",
		})

	proj := NewProject(
		projManifest,
		rec,
	)

	s.Equal("dir", proj.Path())
	s.Equal(projManifest, proj.Manifest())
	s.Equal(rec, proj.Recipe())

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
