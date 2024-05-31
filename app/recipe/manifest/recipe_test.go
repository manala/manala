package manifest

import (
	"bytes"
	_ "embed"
	"manala/app"
	"manala/internal/schema"
	"manala/internal/syncer"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type RecipeSuite struct{ suite.Suite }

func TestRecipeSuite(t *testing.T) {
	suite.Run(t, new(RecipeSuite))
}

func (s *RecipeSuite) Test() {
	dir := filepath.FromSlash("testdata/RecipeSuite/Test/recipe")

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
		dir,
		"recipe",
		manifest,
		repositoryMock,
	)

	s.Equal(dir, recipe.Dir())
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

		s.Require().NoError(err)
		s.Equal("_helpers", out.String())
	})

	s.Run("ProjectManifestTemplate", func() {
		template := recipe.ProjectManifestTemplate()

		out := &bytes.Buffer{}
		err := template.
			WriteTo(out)

		s.Require().NoError(err)
		s.Equal("bar", out.String())
	})

	s.Run("Watches", func() {
		watches, err := recipe.Watches()

		s.Equal([]string{
			dir,
			filepath.Join(dir, "templates"),
		}, watches)
		s.NoError(err)
	})
}
