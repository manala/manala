package recipe

import (
	"bytes"
	_ "embed"
	"github.com/stretchr/testify/suite"
	"manala/app"
	"manala/internal/schema"
	"manala/internal/syncer"
	"path/filepath"
	"testing"
)

type RecipeSuite struct{ suite.Suite }

func TestRecipeSuite(t *testing.T) {
	suite.Run(t, new(RecipeSuite))
}

func (s *RecipeSuite) Test() {
	dir := filepath.FromSlash("testdata/RecipeSuite/Test/recipe")

	mock := &app.RepositoryMock{}

	manifestMock := &app.RecipeManifestMock{}
	manifestMock.
		On("Description").Return("Description").
		On("Icon").Return("Icon").
		On("Template").Return(filepath.Join("templates", "foo.tmpl")).
		On("Vars").Return(map[string]any{"foo": "bar"}).
		On("Sync").Return([]syncer.UnitInterface(nil)).
		On("Schema").Return(schema.Schema{})

	recipe := New(
		dir,
		"recipe",
		manifestMock,
		mock,
	)

	s.Equal(dir, recipe.Dir())
	s.Equal("recipe", recipe.Name())
	s.Equal("Description", recipe.Description())
	s.Equal("Icon", recipe.Icon())
	s.Equal(map[string]any{"foo": "bar"}, recipe.Vars())
	s.Equal([]syncer.UnitInterface(nil), recipe.Sync())
	s.Equal(schema.Schema{}, recipe.Schema())
	s.Equal(mock, recipe.Repository())

	s.Run("Template", func() {
		template := recipe.Template()

		out := &bytes.Buffer{}
		err := template.
			WithDefaultContent(`{{ template "_helpers" }}`).
			WriteTo(out)

		s.NoError(err)

		s.Equal("_helpers", out.String())
	})

	s.Run("ProjectManifestTemplate", func() {
		template := recipe.ProjectManifestTemplate()

		out := &bytes.Buffer{}
		err := template.
			WriteTo(out)

		s.NoError(err)

		s.Equal("bar", out.String())
	})
}
