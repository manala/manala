package manifest_test

import (
	"bytes"
	_ "embed"
	"github.com/manala/manala/app"
	"github.com/manala/manala/app/recipe/manifest"
	"github.com/manala/manala/internal/schema"
	"github.com/manala/manala/internal/sync"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type RecipeSuite struct{ suite.Suite }

func TestRecipeSuite(t *testing.T) {
	suite.Run(t, new(RecipeSuite))
}

func (s *RecipeSuite) Test() {
	m := manifest.New()

	dir := filepath.FromSlash("testdata/RecipeSuite/Test")

	mFile, _ := os.Open(filepath.Join(dir, "manifest.yaml"))
	_, err := m.ReadFrom(mFile)

	s.Require().NoError(err)

	repositoryMock := &app.RepositoryMock{}

	recipe := manifest.NewRecipe(
		dir,
		"recipe",
		m,
		repositoryMock,
	)

	s.Equal(dir, recipe.Dir())
	s.Equal("recipe", recipe.Name())
	s.Equal("description", recipe.Description())
	s.Equal("icon", recipe.Icon())
	s.Equal(map[string]any{"foo": "bar"}, recipe.Vars())
	s.Equal([]sync.UnitInterface{}, recipe.Sync())
	s.Equal(schema.Schema{
		"additionalProperties": false,
		"properties": map[string]interface{}{
			"foo": map[string]interface{}{
				"type": "string",
			},
		},
		"type": "object",
	}, recipe.Schema())
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

		s.Require().NoError(err)
		s.Equal([]string{
			dir,
			filepath.Join(dir, "templates"),
		}, watches)
	})
}
