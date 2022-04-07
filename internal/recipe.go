package internal

import (
	_ "embed"
	"encoding/json"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"github.com/imdario/mergo"
	"github.com/ohler55/ojg/jp"
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
	visitor := &recipeManifestVisitor{manifest: manifest}
	ast.Walk(visitor, contentFile.Docs[0])
	if visitor.err != nil {
		return err
	}

	return nil
}

type recipeManifestVisitor struct {
	manifest *RecipeManifest
	err      error
}

func (visitor *recipeManifestVisitor) Visit(node ast.Node) ast.Visitor {
	nodePath := node.GetPath()

	// Exclude "manala" part
	if strings.HasPrefix(nodePath, "$.manala") {
		return visitor
	}

	switch node := node.(type) {
	case *ast.DocumentNode:
		// Schema root
		visitor.manifest.Schema = map[string]interface{}{
			"type":                 "object",
			"additionalProperties": false,
			"properties":           map[string]interface{}{},
		}
	case *ast.MappingNode:
		// For obscure reasons, comments of the first MappingValue node
		// are set on its Mapping node. Let's copy them from it.
		if len(node.Values) > 0 {
			node.Values[0].Comment = node.Comment
		}
	case *ast.MappingValueNode:
		schema := map[string]interface{}{}

		// Schema type based on node value type
		switch nodeValue := node.Value.(type) {
		case *ast.StringNode, *ast.LiteralNode:
			schema["type"] = "string"
		case *ast.IntegerNode:
			schema["type"] = "integer"
		case *ast.FloatNode:
			schema["type"] = "number"
		case *ast.BoolNode:
			schema["type"] = "boolean"
		case *ast.SequenceNode:
			schema["type"] = "array"
		case *ast.MappingNode:
			schema["type"] = "object"
			schema["properties"] = map[string]interface{}{}
			schema["additionalProperties"] = false
			// Allow additional properties for empty mappings only
			if len(nodeValue.Values) == 0 {
				schema["additionalProperties"] = true
			}
		}

		// Parse comment tags
		if node.Comment != nil {
			tags := &internalYaml.Tags{}
			internalYaml.ParseCommentTags(node.Comment.String(), tags)

			// Handle schema tags
			for _, tag := range *tags.Filter("schema") {
				var tagSchema map[string]interface{}
				if err := json.Unmarshal([]byte(tag.Value), &tagSchema); err != nil {
					visitor.err = internalYaml.CommentTagError(nodePath, err)
					return nil
				}
				if err := mergo.Merge(&schema, tagSchema, mergo.WithOverride); err != nil {
					visitor.err = internalYaml.CommentTagError(nodePath, err)
					return nil
				}
			}

			// Handle options tags
			for _, tag := range *tags.Filter("option") {
				option := &RecipeManifestOption{
					Path:   nodePath,
					Schema: schema,
				}
				if err := json.Unmarshal([]byte(tag.Value), &option); err != nil {
					visitor.err = internalYaml.CommentTagError(nodePath, err)
					return nil
				}
				if err, errs, ok := internalValidator.Validate(recipeManifestOptionSchema, option); !ok {
					visitor.err = ValidationRecipeManifestOptionError(err, errs)
					return nil
				}
				visitor.manifest.Options = append(visitor.manifest.Options, *option)
			}
		}

		// Compute schema path
		_schemaPath, err := jp.ParseString(nodePath)
		if err != nil {
			visitor.err = err
			return nil
		}
		schemaPath := jp.R()
		for _, child := range _schemaPath {
			if _, ok := child.(jp.Child); ok {
				schemaPath = append(schemaPath, jp.Child("properties"), child)
			}
		}

		// Apply schema
		if err := schemaPath.SetOne(visitor.manifest.Schema, schema); err != nil {
			visitor.err = err
			return nil
		}
	}

	return visitor
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
