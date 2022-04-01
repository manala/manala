package loaders

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apex/log"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-yaml"
	yamlAst "github.com/goccy/go-yaml/ast"
	yamlParser "github.com/goccy/go-yaml/parser"
	"github.com/imdario/mergo"
	"github.com/mitchellh/mapstructure"
	"io"
	"io/fs"
	yamlDoc "manala/internal/yaml/doc"
	"manala/models"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
)

func NewRecipeLoader(log log.Interface, fsManager models.FsManagerInterface) RecipeLoaderInterface {
	return &recipeLoader{
		log:       log,
		fsManager: fsManager,
	}
}

type RecipeLoaderInterface interface {
	Load(name string, repository models.RepositoryInterface) (models.RecipeInterface, error)
	Walk(repository models.RepositoryInterface, fn recipeWalkFunc) error
}

type recipeConfig struct {
	Description string `validate:"required"`
	Template    string
	Sync        []models.RecipeSyncUnit
}

type recipeLoader struct {
	log       log.Interface
	fsManager models.FsManagerInterface
}

func (ld *recipeLoader) Load(name string, repository models.RepositoryInterface) (models.RecipeInterface, error) {
	var recipe models.RecipeInterface

	if err := ld.Walk(repository, func(rec models.RecipeInterface) {
		if rec.Name() == name {
			recipe = rec
		}
	}); err != nil {
		return nil, err
	}

	if recipe != nil {
		return recipe, nil
	}

	return nil, fmt.Errorf("recipe not found")
}

func (ld *recipeLoader) Walk(repository models.RepositoryInterface, fn recipeWalkFunc) error {
	// Repository file system
	repoFs := ld.fsManager.NewModelFs(repository)

	files, err := repoFs.ReadDir("")
	if err != nil {
		return err
	}

	for _, file := range files {
		// Exclude dot files
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}

		// Keep dirs only
		if !file.IsDir() {
			continue
		}

		manifest, err := repoFs.Open(filepath.Join(file.Name(), models.RecipeManifestFile))
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				// Dir does not contain a manifest
				continue
			}
			return err
		}

		rec, err := ld.loadDir(file.Name(), manifest, repository)
		if err != nil {
			return err
		}

		fn(rec)
	}

	return nil
}

type recipeWalkFunc func(rec models.RecipeInterface)

func (ld *recipeLoader) loadDir(dir string, manifest fs.File, repository models.RepositoryInterface) (models.RecipeInterface, error) {
	ld.log.WithField("dir", dir).Debug("Loading recipe...")

	// Parse manifest
	manifestContent, err := io.ReadAll(manifest)
	if err != nil {
		return nil, err
	}

	manifestNode, err := yamlParser.ParseBytes(manifestContent, yamlParser.ParseComments)
	if err != nil {
		return nil, err
	}

	var vars map[string]interface{}
	if err := yaml.NewDecoder(manifestNode).Decode(&vars); err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("empty recipe manifest \"%s\"", dir)
		}
		return nil, fmt.Errorf("incorrect recipe manifest \"%s\" %s", dir, yaml.FormatError(err, true, true))
	}

	// Map config
	cfg := recipeConfig{}
	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:     &cfg,
		DecodeHook: recipeStringToSyncUnitHookFunc(),
	})
	if err := decoder.Decode(vars["manala"]); err != nil {
		return nil, err
	}

	// Validate
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, fmt.Errorf("invalid recipe manifest config \"%s\" (%w)", dir, err)
	}

	// Cleanup vars
	delete(vars, "manala")

	// Parse manifest node
	var options []models.RecipeOption
	schema, err := ld.parseManifestNode(manifestNode.Docs[0], &options, "")
	if err != nil {
		return nil, err
	}

	return models.NewRecipe(
		dir,
		cfg.Description,
		cfg.Template,
		dir,
		repository,
		vars,
		cfg.Sync,
		schema,
		options,
	), nil
}

func (ld *recipeLoader) parseManifestNode(node yamlAst.Node, options *[]models.RecipeOption, root string) (map[string]interface{}, error) {
	var nodes []*yamlAst.MappingValueNode

	switch node := node.(type) {
	case *yamlAst.DocumentNode:
		schema, err := ld.parseManifestNode(node.Body, options, "/")
		if err != nil {
			return nil, err
		}
		return schema, nil
	case *yamlAst.MappingValueNode:
		nodes = []*yamlAst.MappingValueNode{node}
	case *yamlAst.MappingNode:
		nodes = node.Values
	}

	schemaProperties := map[string]interface{}{}

	for _, n := range nodes {
		nPath := path.Join(root, n.Key.String())

		// Exclude "manala" config
		if nPath == "/manala" {
			continue
		}

		var schema map[string]interface{} = nil

		switch n.Value.(type) {
		case yamlAst.ScalarNode:
			schema = map[string]interface{}{}
		case *yamlAst.MappingValueNode, *yamlAst.MappingNode:
			var err error
			schema, err = ld.parseManifestNode(n.Value, options, nPath)
			if err != nil {
				return nil, err
			}
		case *yamlAst.SequenceNode:
			schema = map[string]interface{}{
				"type": "array",
			}
		default:
			return nil, fmt.Errorf("unknown node type: %s", reflect.TypeOf(n))
		}

		if comment := n.GetComment(); comment != nil {
			tags := yamlDoc.ParseCommentTags(comment.String())
			// Handle schema tags
			for _, tag := range tags.Filter("schema") {
				var tagSchema map[string]interface{}
				if err := json.Unmarshal([]byte(tag.Value), &tagSchema); err != nil {
					return nil, fmt.Errorf("invalid recipe schema tag at \"%s\": %w", nPath, err)
				}
				if err := mergo.Merge(&schema, tagSchema, mergo.WithOverride); err != nil {
					return nil, fmt.Errorf("unable to merge recipe schema tag at \"%s\": %w", nPath, err)
				}
			}
			// Handle option tags
			for _, tag := range tags.Filter("option") {
				option := &models.RecipeOption{
					Path:   nPath,
					Schema: schema,
				}
				if err := json.Unmarshal([]byte(tag.Value), &option); err != nil {
					return nil, fmt.Errorf("invalid recipe option tag at \"%s\": %w", nPath, err)
				}
				validate := validator.New()
				if err := validate.Struct(option); err != nil {
					return nil, fmt.Errorf("incorrect recipe option tag at \"%s\": %w", nPath, err)
				}
				*options = append(*options, *option)
			}
		}

		schemaProperties[n.Key.String()] = schema
	}

	// Allow additional properties for empty mappings only
	schemaAdditionalProperties := false
	if len(nodes) == 0 {
		schemaAdditionalProperties = true
	}

	return map[string]interface{}{
		"type":                 "object",
		"additionalProperties": schemaAdditionalProperties,
		"properties":           schemaProperties,
	}, nil
}

// Returns a DecodeHookFunc that converts strings to syncUnit
func recipeStringToSyncUnitHookFunc() mapstructure.DecodeHookFunc {
	return func(rf reflect.Type, rt reflect.Type, data interface{}) (interface{}, error) {
		if rf.Kind() != reflect.String {
			return data, nil
		}
		if rt != reflect.TypeOf(models.RecipeSyncUnit{}) {
			return data, nil
		}

		src := data.(string)
		dst := src

		// Separate source / destination
		u := strings.Split(src, " ")
		if len(u) > 1 {
			src = u[0]
			dst = u[1]
		}

		return models.RecipeSyncUnit{
			Source:      src,
			Destination: dst,
		}, nil
	}
}
