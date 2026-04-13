package manifest_test

import (
	_ "embed"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/recipe/manifest"
	"github.com/manala/manala/internal/schema"
	"github.com/manala/manala/internal/sync"

	"github.com/stretchr/testify/suite"
)

type RecipeSuite struct{ suite.Suite }

func TestRecipeSuite(t *testing.T) {
	suite.Run(t, new(RecipeSuite))
}

func (s *RecipeSuite) Test() {
	m := manifest.New()

	dir := filepath.FromSlash("testdata/RecipeSuite/Test")

	reader, _ := os.Open(filepath.Join(dir, "manifest.yaml"))
	content, _ := io.ReadAll(reader)

	err := m.UnmarshalYAML(content)

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
		"properties": map[string]any{
			"foo": map[string]any{
				"type": "string",
			},
		},
		"type": "object",
	}, recipe.Schema())
	s.Equal(repositoryMock, recipe.Repository())

	s.Run("Template", func() {
		s.Equal(filepath.Join(dir, "templates", "foo.tmpl"), recipe.Template())
	})

	s.Run("Partials", func() {
		s.Empty(recipe.Partials())
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

func (s *RecipeSuite) TestPartials() {
	m := manifest.New()

	dir := filepath.FromSlash("testdata/RecipeSuite/TestPartials")

	reader, _ := os.Open(filepath.Join(dir, "manifest.yaml"))
	content, _ := io.ReadAll(reader)

	err := m.UnmarshalYAML(content)
	s.Require().NoError(err)

	recipe := manifest.NewRecipe(
		dir,
		"recipe",
		m,
		&app.RepositoryMock{},
	)

	s.Equal([]string{
		filepath.Join(dir, "partial.tmpl"),
		filepath.Join(dir, "dir/partial.tmpl"),
	}, recipe.Partials())
}

func (s *RecipeSuite) TestPartialsHelpers() {
	m := manifest.New()

	dir := filepath.FromSlash("testdata/RecipeSuite/TestPartialsHelpers")

	reader, _ := os.Open(filepath.Join(dir, "manifest.yaml"))
	content, _ := io.ReadAll(reader)

	err := m.UnmarshalYAML(content)
	s.Require().NoError(err)

	repositoryMock := &app.RepositoryMock{}

	recipe := manifest.NewRecipe(
		dir,
		"recipe",
		m,
		repositoryMock,
	)

	s.Equal([]string{
		filepath.Join(dir, "_helpers.tmpl"),
	}, recipe.Partials())
}

func (s *RecipeSuite) TestPartialsNoHelpers() {
	m := manifest.New()

	dir := filepath.FromSlash("testdata/RecipeSuite/TestPartialsNoHelpers")

	reader, _ := os.Open(filepath.Join(dir, "manifest.yaml"))
	content, _ := io.ReadAll(reader)

	err := m.UnmarshalYAML(content)
	s.Require().NoError(err)

	repositoryMock := &app.RepositoryMock{}

	recipe := manifest.NewRecipe(
		dir,
		"recipe",
		m,
		repositoryMock,
	)

	s.Empty(recipe.Partials())
}
