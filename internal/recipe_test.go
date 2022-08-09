package internal

import (
	"bytes"
	_ "embed"
	"github.com/stretchr/testify/suite"
	"io"
	internalTesting "manala/internal/testing"
	"os"
	"path/filepath"
	"testing"
)

type RecipeSuite struct{ suite.Suite }

func TestRecipeSuite(t *testing.T) {
	suite.Run(t, new(RecipeSuite))
}

func (s *RecipeSuite) Test() {
	repository := &Repository{path: "repository"}
	recipeManifest := NewRecipeManifest(internalTesting.DataPath(s, "recipe"))
	recipeManifest.Description = "Recipe"
	recipeManifest.Vars = map[string]interface{}{"foo": "bar"}
	recipeManifest.Sync = []RecipeManifestSyncUnit{}
	recipeManifest.Schema = map[string]interface{}{}
	recipeManifest.Options = []RecipeManifestOption{}
	recipe := &Recipe{
		name:       "recipe",
		manifest:   recipeManifest,
		repository: repository,
	}

	s.Equal(internalTesting.DataPath(s, "recipe"), recipe.Path())
	s.Equal("recipe", recipe.Name())
	s.Equal("Recipe", recipe.Description())
	s.Equal(map[string]interface{}{"foo": "bar"}, recipe.Vars())
	s.Equal([]RecipeManifestSyncUnit{}, recipe.Sync())
	s.Equal(map[string]interface{}{}, recipe.Schema())
	s.Equal([]RecipeManifestOption{}, recipe.Options())
	s.Equal(repository, recipe.Repository())

	s.Run("Template", func() {
		template := recipe.Template()

		out := &bytes.Buffer{}
		err := template.
			WithDefaultContent(`{{ template "_helpers" }}`).
			Write(out)

		s.NoError(err)
		s.Equal("_helpers", out.String())
	})

	s.Run("ProjectManifestTemplate", func() {
		recipeManifest.Template = filepath.Join("templates", "foo.tmpl")

		template := recipe.ProjectManifestTemplate()

		out := &bytes.Buffer{}
		err := template.
			Write(out)

		s.NoError(err)
		s.Equal("bar", out.String())
	})

	s.Run("NewProject", func() {
		project := recipe.NewProject("dir")

		s.Equal("dir", project.Path())
		s.Equal("repository", project.Manifest().Repository)
		s.Equal("recipe", project.Manifest().Recipe)
		s.Equal(map[string]interface{}{"foo": "bar"}, project.Manifest().Vars)
	})
}

func (s *RecipeSuite) TestManifest() {
	recipeManifest := NewRecipeManifest("dir")

	s.Equal(filepath.Join("dir", ".manala.yaml"), recipeManifest.path)

	s.Run("Write", func() {
		length, err := recipeManifest.Write([]byte("foo"))
		s.NoError(err)
		s.Equal(3, length)
		s.Equal("foo", string(recipeManifest.content))
	})
}

func (s *RecipeSuite) TestManifestLoad() {

	s.Run("Valid", func() {
		recipeManifest := NewRecipeManifest("")

		file, _ := os.Open(internalTesting.DataPath(s, "manifest.yaml"))
		_, _ = io.Copy(recipeManifest, file)
		err := recipeManifest.Load()

		s.NoError(err)
		s.Equal("description", recipeManifest.Description)
		s.Equal("template", recipeManifest.Template)
		s.Equal(map[string]interface{}{
			"foo": map[string]interface{}{
				"bar": "baz",
			},
			"underscore_key": "ok",
		}, recipeManifest.Vars)
		s.Equal([]RecipeManifestSyncUnit{
			{Source: "sync", Destination: "sync"},
		}, recipeManifest.Sync)
		s.Equal(map[string]interface{}{
			"type":                 "object",
			"additionalProperties": false,
			"properties": map[string]interface{}{
				"foo": map[string]interface{}{
					"properties": map[string]interface{}{
						"bar": map[string]interface{}{
							"type":      "string",
							"minLength": float64(1),
						},
					},
				},
				"underscore_key": map[string]interface{}{
					"type": "string",
				},
			},
		}, recipeManifest.Schema)
		s.Equal([]RecipeManifestOption{
			{
				Label: "Bar",
				Path:  "$.foo.bar",
				Schema: map[string]interface{}{
					"type":      "string",
					"minLength": float64(1),
				},
			},
		}, recipeManifest.Options)
	})

	s.Run("Invalid Yaml", func() {
		recipeManifest := NewRecipeManifest("")

		file, _ := os.Open(internalTesting.DataPath(s, "manifest.yaml"))
		_, _ = io.Copy(recipeManifest, file)
		err := recipeManifest.Load()

		s.ErrorAs(err, &internalError)
		s.Equal("yaml processing error", internalError.Message)
	})

	s.Run("Empty", func() {
		recipeManifest := NewRecipeManifest("")

		file, _ := os.Open(internalTesting.DataPath(s, "manifest.yaml"))
		_, _ = io.Copy(recipeManifest, file)
		err := recipeManifest.Load()

		s.ErrorAs(err, &internalError)
		s.Equal("empty recipe manifest", internalError.Message)
	})

	s.Run("Wrong", func() {
		recipeManifest := NewRecipeManifest("")

		file, _ := os.Open(internalTesting.DataPath(s, "manifest.yaml"))
		_, _ = io.Copy(recipeManifest, file)
		err := recipeManifest.Load()

		s.ErrorAs(err, &internalError)
		s.Equal("yaml processing error", internalError.Message)
	})

	s.Run("Invalid", func() {
		recipeManifest := NewRecipeManifest("")

		file, _ := os.Open(internalTesting.DataPath(s, "manifest.yaml"))
		_, _ = io.Copy(recipeManifest, file)
		err := recipeManifest.Load()

		s.ErrorAs(err, &internalError)
		s.Equal("recipe validation error", internalError.Message)
	})
}
