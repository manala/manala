package manifest

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
	recipeDir := filepath.FromSlash("testdata/RecipeSuite/Test/recipe")

	repositoryMock := &app.RepositoryMock{}

	manifest := &Manifest{
		config: &config{
			Description: "description",
			Icon:        "icon",
			Template:    filepath.Join("templates", "foo.tmpl"),
			Sync:        []syncer.UnitInterface(nil),
		},
		vars:   map[string]any{"foo": "bar"},
		schema: schema.Schema{},
	}

	recipe := NewRecipe(
		recipeDir,
		"recipe",
		manifest,
		repositoryMock,
	)

	s.Equal(recipeDir, recipe.Dir())
	s.Equal("recipe", recipe.Name())
	s.Equal("description", recipe.Description())
	s.Equal("icon", recipe.Icon())
	s.Equal(map[string]any{"foo": "bar"}, recipe.Vars())
	s.Equal([]syncer.UnitInterface(nil), recipe.Sync())
	s.Equal(schema.Schema{}, recipe.Schema())
	s.Equal(repositoryMock, recipe.Repository())

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
