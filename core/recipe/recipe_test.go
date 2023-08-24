package recipe

import (
	"bytes"
	_ "embed"
	"github.com/stretchr/testify/suite"
	"manala/app/mocks"
	"manala/internal/syncer"
	"path/filepath"
	"testing"
)

type RecipeSuite struct{ suite.Suite }

func TestRecipeSuite(t *testing.T) {
	suite.Run(t, new(RecipeSuite))
}

func (s *RecipeSuite) Test() {
	recDir := filepath.FromSlash("testdata/RecipeSuite/Test/recipe")

	repoMock := &mocks.RepositoryMock{}

	recManifestMock := &mocks.RecipeManifestMock{}
	recManifestMock.
		On("Description").Return("Description").
		On("Template").Return(filepath.Join("templates", "foo.tmpl")).
		On("Vars").Return(map[string]interface{}{"foo": "bar"}).
		On("Sync").Return([]syncer.UnitInterface(nil)).
		On("Schema").Return(map[string]interface{}{}).
		On("InitVars").Return(map[string]interface{}{"foo": "baz"})

	rec := NewRecipe(
		recDir,
		"recipe",
		recManifestMock,
		repoMock,
	)

	s.Equal(recDir, rec.Dir())
	s.Equal("recipe", rec.Name())
	s.Equal("Description", rec.Description())
	s.Equal(map[string]interface{}{"foo": "bar"}, rec.Vars())
	s.Equal([]syncer.UnitInterface(nil), rec.Sync())
	s.Equal(map[string]interface{}{}, rec.Schema())
	s.Equal(repoMock, rec.Repository())

	s.Run("Template", func() {
		template := rec.Template()

		out := &bytes.Buffer{}
		err := template.
			WithDefaultContent(`{{ template "_helpers" }}`).
			WriteTo(out)

		s.NoError(err)

		s.Equal("_helpers", out.String())
	})

	s.Run("ProjectManifestTemplate", func() {
		template := rec.ProjectManifestTemplate()

		out := &bytes.Buffer{}
		err := template.
			WriteTo(out)

		s.NoError(err)

		s.Equal("bar", out.String())
	})
}
