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
		s.Equal("Description", recipeManifest.Description)
		s.Equal("template", recipeManifest.Template)
		s.Equal(map[string]interface{}{
			"string":      "string",
			"string_null": nil,
			"sequence": []interface{}{
				"first",
			},
			"sequence_string_empty": []interface{}{},
			"boolean":               true,
			"integer":               uint64(123),
			"float":                 1.2,
			"map": map[string]interface{}{
				"string": "string",
				"map": map[string]interface{}{
					"string": "string",
				},
			},
			"map_empty": map[string]interface{}{},
			"map_single": map[string]interface{}{
				"first": "foo",
			},
			"map_multiple": map[string]interface{}{
				"first":  "foo",
				"second": "foo",
			},
			"enum":           nil,
			"underscore_key": "ok",
			"hyphen-key":     "ok",
			"dot.key":        "ok",
		}, recipeManifest.Vars)
		s.Equal([]RecipeManifestSyncUnit{
			{Source: "file", Destination: "file"},
			{Source: "dir/file", Destination: "dir/file"},
			{Source: "file", Destination: "dir/file"},
			{Source: "dir/file", Destination: "file"},
			{Source: "src_file", Destination: "dst_file"},
			{Source: "src_dir/file", Destination: "dst_dir/file"},
		}, recipeManifest.Sync)
		s.Equal(map[string]interface{}{
			"type":                 "object",
			"additionalProperties": false,
			"properties": map[string]interface{}{
				"string": map[string]interface{}{
					"type": "string",
				},
				"string_null": map[string]interface{}{
					"type": "string",
				},
				"sequence": map[string]interface{}{
					"type": "array",
				},
				"sequence_string_empty": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"type": "string",
					},
				},
				"boolean": map[string]interface{}{
					"type": "boolean",
				},
				"integer": map[string]interface{}{
					"type": "integer",
				},
				"float": map[string]interface{}{
					"type": "number",
				},
				"map": map[string]interface{}{
					"type":                 "object",
					"additionalProperties": false,
					"properties": map[string]interface{}{
						"string": map[string]interface{}{
							"type": "string",
						},
						"map": map[string]interface{}{
							"type":                 "object",
							"additionalProperties": false,
							"properties": map[string]interface{}{
								"string": map[string]interface{}{
									"type": "string",
								},
							},
						},
					},
				},
				"map_empty": map[string]interface{}{
					"type":                 "object",
					"additionalProperties": true,
				},
				"map_single": map[string]interface{}{
					"type":                 "object",
					"additionalProperties": false,
					"properties": map[string]interface{}{
						"first": map[string]interface{}{
							"type":      "string",
							"minLength": float64(1),
						},
					},
				},
				"map_multiple": map[string]interface{}{
					"type":                 "object",
					"additionalProperties": false,
					"properties": map[string]interface{}{
						"first": map[string]interface{}{
							"type":      "string",
							"minLength": float64(1),
						},
						"second": map[string]interface{}{
							"type":      "string",
							"minLength": float64(1),
						},
					},
				},
				"enum": map[string]interface{}{
					"enum": []interface{}{
						nil,
						true,
						false,
						"string",
						float64(12),
						2.3,
						3.0,
						"3.0",
					},
				},
				"underscore_key": map[string]interface{}{
					"type": "string",
				},
				"hyphen-key": map[string]interface{}{
					"type": "string",
				},
				"dot.key": map[string]interface{}{
					"type": "string",
				},
			},
		}, recipeManifest.Schema)
		s.Equal([]RecipeManifestOption{
			{
				Label: "String",
				Path:  "$.string",
				Schema: map[string]interface{}{
					"type": "string",
				},
			},
			{
				Label: "String null",
				Path:  "$.string_null",
				Schema: map[string]interface{}{
					"type": "string",
				},
			},
			{
				Label: "Map single first",
				Path:  "$.map_single.first",
				Schema: map[string]interface{}{
					"type":      "string",
					"minLength": float64(1),
				},
			},
			{
				Label: "Map multiple first",
				Path:  "$.map_multiple.first",
				Schema: map[string]interface{}{
					"type":      "string",
					"minLength": float64(1),
				},
			},
			{
				Label: "Map multiple second",
				Path:  "$.map_multiple.second",
				Schema: map[string]interface{}{
					"type":      "string",
					"minLength": float64(1),
				},
			},
			{
				Label: "Enum null",
				Path:  "$.enum",
				Schema: map[string]interface{}{
					"enum": []interface{}{
						nil,
						true,
						false,
						"string",
						float64(12),
						2.3,
						3.0,
						"3.0",
					},
				},
			},
			{
				Label: "Underscore key",
				Path:  "$.underscore_key",
				Schema: map[string]interface{}{
					"type": "string",
				},
			},
			{
				Label: "Hyphen key",
				Path:  "$.hyphen-key",
				Schema: map[string]interface{}{
					"type": "string",
				},
			},
			{
				Label: "Dot key",
				Path:  "$.'dot.key'",
				Schema: map[string]interface{}{
					"type": "string",
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
