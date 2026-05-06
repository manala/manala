package manifest //nolint:testpackage

import (
	_ "embed"
	"path/filepath"
	"testing"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/recipe/option"
	"github.com/manala/manala/app/sync"
	"github.com/manala/manala/app/testing/mocks"

	"github.com/stretchr/testify/suite"
)

type RecipeSuite struct{ suite.Suite }

func TestRecipeSuite(t *testing.T) {
	suite.Run(t, new(RecipeSuite))
}

func (s *RecipeSuite) Test() {
	dir := filepath.FromSlash("testdata/RecipeSuite/Test/repository/recipe")
	name := "name"
	configPartial := "partial"
	config := &Config{
		Description: "description",
		Icon:        "icon",
		Partials:    []string{configPartial},
		Sync:        []sync.Unit{{}},
	}
	repositoryMock := &mocks.Repository{}
	vars := map[string]any{"foo": "bar"}
	schema := map[string]any{"bar": "baz"}
	options := []app.RecipeOption{&option.String{}}

	recipe := &Recipe{
		dir:        dir,
		name:       name,
		config:     config,
		repository: repositoryMock,
		vars:       vars,
		schema:     schema,
		options:    options,
	}

	s.Equal(dir, recipe.Dir())
	s.Equal(name, recipe.Name())
	s.Equal(config.Description, recipe.Description())
	s.Equal(config.Icon, recipe.Icon())
	s.Equal([]string{filepath.Join(dir, configPartial)}, recipe.Partials())
	s.Equal(config.Sync, recipe.Sync())
	s.Equal(repositoryMock, recipe.Repository())
	s.Equal(vars, recipe.Vars())
	s.Equal(schema, recipe.Schema())
	s.Equal(options, recipe.Options())

	watches, err := recipe.Watches()
	s.Require().NoError(err)
	s.Equal([]string{
		dir,
		filepath.Join(dir, "bar"),
		filepath.Join(dir, "bar", "baz"),
	}, watches)
}

func (s *RecipeSuite) TestTemplate() {
	s.Run("WithConfigTemplate", func() {
		dir := "dir"
		config := &Config{
			Template: "template",
		}
		recipe := &Recipe{
			dir:    dir,
			config: config,
		}
		s.Equal(filepath.Join(dir, recipe.config.Template), recipe.Template())
	})
	s.Run("WithoutConfigTemplate", func() {
		dir := "dir"
		config := &Config{}
		recipe := &Recipe{
			dir:    dir,
			config: config,
		}
		s.Empty(recipe.Template())
	})
}

func (s *RecipeSuite) TestLegacyPartials() {
	s.Run("WithConfigPartials", func() {
		dir := filepath.FromSlash("testdata/RecipeSuite/TestLegacyPartials/repository/recipe")
		config := &Config{
			Partials: []string{"foo", "bar"},
		}
		recipe := &Recipe{
			dir:    dir,
			config: config,
		}
		s.Equal([]string{
			filepath.Join(dir, "foo"),
			filepath.Join(dir, "bar"),
		}, recipe.Partials())
	})
	s.Run("WithoutConfigPartials", func() {
		dir := filepath.FromSlash("testdata/RecipeSuite/TestLegacyPartials/repository/recipe")
		config := &Config{}
		recipe := &Recipe{
			dir:    dir,
			config: config,
		}
		s.Equal([]string{
			filepath.Join(dir, "_helpers.tmpl"),
		}, recipe.Partials())
	})
}
