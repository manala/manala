package internal

import (
	_ "embed"
	"encoding/json"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"github.com/imdario/mergo"
	"io"
	internalTemplate "manala/internal/template"
	internalValidator "manala/internal/validator"
	internalYaml "manala/internal/yaml"
	"os"
	"path/filepath"
	"strings"
)

type Recipe struct {
	name       string
	manifest   *RecipeManifest
	repository *Repository
}

func (recipe *Recipe) Path() string {
	return filepath.Dir(recipe.manifest.path)
}

func (recipe *Recipe) Name() string {
	return recipe.name
}

func (recipe *Recipe) Description() string {
	return recipe.manifest.Description
}

func (recipe *Recipe) Vars() map[string]interface{} {
	return recipe.manifest.Vars
}

func (recipe *Recipe) Sync() []RecipeManifestSyncUnit {
	return recipe.manifest.Sync
}

func (recipe *Recipe) Schema() map[string]interface{} {
	return recipe.manifest.Schema
}

func (recipe *Recipe) Options() []RecipeManifestOption {
	return recipe.manifest.Options
}

func (recipe *Recipe) Repository() *Repository {
	return recipe.repository
}

func (recipe *Recipe) Template() *internalTemplate.Template {
	template := &internalTemplate.Template{}

	// Include template helpers if any
	helpersPath := filepath.Join(recipe.Path(), "_helpers.tmpl")
	if _, err := os.Stat(helpersPath); err == nil {
		template.WithDefaultFile(helpersPath)
	}

	return template
}

func (recipe *Recipe) ProjectManifestTemplate() *internalTemplate.Template {
	template := recipe.Template()

	if recipe.manifest.Template != "" {
		template.WithFile(filepath.Join(recipe.Path(), recipe.manifest.Template))
	}

	return template
}

func (recipe *Recipe) NewProject(path string) *Project {
	projectManifest := NewProjectManifest(path)
	projectManifest.Repository = recipe.repository.Path()
	projectManifest.Recipe = recipe.name
	projectManifest.Vars = recipe.manifest.Vars

	return &Project{
		manifest: projectManifest,
		recipe:   recipe,
	}
}

/************/
/* Manifest */
/************/

//go:embed recipe_manifest.schema.json
var recipeManifestSchema string

const recipeManifestFile = ".manala.yaml"

func NewRecipeManifest(dir string) *RecipeManifest {
	return &RecipeManifest{
		path: filepath.Join(dir, recipeManifestFile),
		Vars: map[string]interface{}{},
		Schema: map[string]interface{}{
			"type":                 "object",
			"additionalProperties": false,
			"properties":           map[string]interface{}{},
		},
	}
}

type RecipeManifest struct {
	path        string
	content     []byte
	Description string
	Template    string
	Vars        map[string]interface{}
	Sync        []RecipeManifestSyncUnit
	Schema      map[string]interface{}
	Options     []RecipeManifestOption
}

// Write implement the Writer interface
func (manifest *RecipeManifest) Write(content []byte) (n int, err error) {
	manifest.content = append(manifest.content, content...)
	return len(content), nil
}

// Load parse and decode content
func (manifest *RecipeManifest) Load() error {
	// Parse content
	contentFile, err := parser.ParseBytes(manifest.content, parser.ParseComments)
	if err != nil {
		return internalYaml.Error(manifest.path, err)
	}

	// Decode content
	var contentData map[string]interface{}
	if err := yaml.NewDecoder(contentFile).Decode(&contentData); err != nil {
		if err == io.EOF {
			return EmptyRecipeManifestPathError(manifest.path)
		}
		return internalYaml.Error(manifest.path, err)
	}

	// Validate content
	if err, errs, ok := internalValidator.Validate(recipeManifestSchema, contentData, internalValidator.WithYamlContent(manifest.content)); !ok {
		return ValidationRecipeManifestPathError(manifest.path, err, errs)
	}

	// Decode manifest with "manala" data part
	manalaPath, _ := yaml.PathString("$.manala")
	manalaNode, _ := manalaPath.FilterFile(contentFile)
	if err := yaml.NewDecoder(manalaNode).Decode(manifest); err != nil {
		return internalYaml.Error(manifest.path, err)
	}

	// Keep remaining data for manifest vars
	delete(contentData, "manala")
	manifest.Vars = contentData

	// Parse manifest (options & schema)
	if err := parseRecipeManifest(
		contentFile.Docs[0].Body,
		manifest.Schema["properties"].(map[string]interface{}),
		&manifest.Options,
	); err != nil {
		return err
	}

	return nil
}

func parseRecipeManifest(node ast.Node, properties map[string]interface{}, options *[]RecipeManifestOption) error {
	switch node := node.(type) {
	case *ast.MappingNode:
		// Comments of the first MappingValueNode are set on its MappingNode.
		// Work around by manually copy them from it.
		// See: https://github.com/goccy/go-yaml/issues/311
		if len(node.Values) > 0 {
			node.Values[0].Comment = node.Comment
		}

		// Range over mapping value nodes
		for _, valueNode := range node.Values {
			if err := parseRecipeManifest(valueNode, properties, options); err != nil {
				return err
			}
		}
	case *ast.MappingValueNode:
		nodePath := node.GetPath()

		// Exclude "manala" path
		if strings.HasPrefix(nodePath, "$.manala") {
			return nil
		}

		// Property name based on node key
		propertyName := node.Key.String()

		// Property schema based on node value type
		propertySchema := map[string]interface{}{}

		switch node := node.Value.(type) {
		case *ast.StringNode, *ast.LiteralNode:
			propertySchema["type"] = "string"
		case *ast.IntegerNode:
			propertySchema["type"] = "integer"
		case *ast.FloatNode:
			propertySchema["type"] = "number"
		case *ast.BoolNode:
			propertySchema["type"] = "boolean"
		case *ast.SequenceNode:
			propertySchema["type"] = "array"
		case *ast.MappingValueNode, *ast.MappingNode:
			// Could be either MappingNode or MappingValueNode
			// depending on the number of their items
			// See: https://github.com/goccy/go-yaml/issues/310

			propertySchema["type"] = "object"

			// Allow additional properties for empty mapping nodes only
			if n, ok := node.(*ast.MappingNode); ok && len(n.Values) == 0 {
				propertySchema["additionalProperties"] = true
			} else {
				propertySchema["additionalProperties"] = false
				propertySchema["properties"] = map[string]interface{}{}

				if err := parseRecipeManifest(
					node,
					propertySchema["properties"].(map[string]interface{}),
					options,
				); err != nil {
					return err
				}
			}
		}

		// Parse comment tags
		nodeComment := node.GetComment()
		if nodeComment != nil {
			tags := &internalYaml.Tags{}
			internalYaml.ParseCommentTags(nodeComment.String(), tags)

			// Handle schema tags
			for _, tag := range *tags.Filter("schema") {
				var tagSchemaProperties map[string]interface{}
				if err := json.Unmarshal([]byte(tag.Value), &tagSchemaProperties); err != nil {
					return internalYaml.CommentTagError(nodePath, err)
				}
				if err := mergo.Merge(&propertySchema, tagSchemaProperties, mergo.WithOverride); err != nil {
					return internalYaml.CommentTagError(nodePath, err)
				}
			}

			// Handle options tags
			for _, tag := range *tags.Filter("option") {
				option := &RecipeManifestOption{
					Path:   nodePath,
					Schema: propertySchema,
				}
				if err := json.Unmarshal([]byte(tag.Value), &option); err != nil {
					return internalYaml.CommentTagError(nodePath, err)
				}
				if err, errs, ok := internalValidator.Validate(recipeManifestOptionSchema, option); !ok {
					return ValidationRecipeManifestOptionError(err, errs)
				}
				*options = append(*options, *option)
			}
		}

		properties[propertyName] = propertySchema
	}

	return nil
}

type RecipeManifestSyncUnit struct {
	Source      string
	Destination string
}

func (unit *RecipeManifestSyncUnit) UnmarshalYAML(value []byte) error {
	unit.Source = string(value)
	unit.Destination = unit.Source

	// Separate source / destination
	splits := strings.Split(unit.Source, " ")
	if len(splits) > 1 {
		unit.Source = splits[0]
		unit.Destination = splits[1]
	}

	return nil
}

//go:embed recipe_manifest_option.schema.json
var recipeManifestOptionSchema string

type RecipeManifestOption struct {
	Label  string                 `json:"label" validate:"required"`
	Path   string                 `json:"path"`
	Schema map[string]interface{} `json:"schema"`
}
