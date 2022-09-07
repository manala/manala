package recipe

import (
	"bytes"
	_ "embed"
	"github.com/stretchr/testify/suite"
	"manala/core"
	internalSyncer "manala/internal/syncer"
	internalTesting "manala/internal/testing"
	"path/filepath"
	"testing"
)

type RecipeSuite struct{ suite.Suite }

func TestRecipeSuite(t *testing.T) {
	suite.Run(t, new(RecipeSuite))
}

func (s *RecipeSuite) Test() {
	repo := core.NewRepositoryMock()

	recManifest := core.NewRecipeManifestMock().
		WithDir(internalTesting.DataPath(s, "recipe")).
		WithDescription("Description").
		WithTemplate(filepath.Join("templates", "foo.tmpl")).
		WithVars(map[string]interface{}{"foo": "bar"}).
		WithInitVars(map[string]interface{}{"foo": "baz"})

	rec := NewRecipe(
		"recipe",
		recManifest,
		repo,
	)

	s.Equal(internalTesting.DataPath(s, "recipe"), rec.Path())
	s.Equal("recipe", rec.Name())
	s.Equal("Description", rec.Description())
	s.Equal(map[string]interface{}{"foo": "bar"}, rec.Vars())
	s.Equal([]internalSyncer.UnitInterface(nil), rec.Sync())
	s.Equal(map[string]interface{}{}, rec.Schema())
	s.Equal(repo, rec.Repository())

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
