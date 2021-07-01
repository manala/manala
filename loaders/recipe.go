package loaders

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/imdario/mergo"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
	"io"
	"io/fs"
	"manala/logger"
	"manala/models"
	"manala/yaml/cleaner"
	"manala/yaml/doc"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

func NewRecipeLoader(log logger.Logger, fsManager models.FsManagerInterface) RecipeLoaderInterface {
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
	log       logger.Logger
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
	ld.log.Debug("Loading recipe...", ld.log.WithField("dir", dir))

	// Parse manifest
	node := yaml.Node{}

	if err := yaml.NewDecoder(manifest).Decode(&node); err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("empty recipe manifest \"%s\"", dir)
		}
		return nil, fmt.Errorf("invalid recipe manifest \"%s\" (%w)", dir, err)
	}

	var vars map[string]interface{}
	if err := node.Decode(&vars); err != nil {
		return nil, fmt.Errorf("incorrect recipe manifest \"%s\" (%w)", dir, err)
	}

	// See: https://github.com/go-yaml/yaml/issues/139
	vars = cleaner.Clean(vars)

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
		return nil, err
	}

	// Cleanup vars
	delete(vars, "manala")

	// Parse config node
	var options []models.RecipeOption
	schema, err := ld.parseConfigNode(&node, &options, "")
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

func (ld *recipeLoader) parseConfigNode(node *yaml.Node, options *[]models.RecipeOption, root string) (map[string]interface{}, error) {
	var nodeKey *yaml.Node = nil
	schemaProperties := map[string]interface{}{}

	for _, nodeChild := range node.Content {
		// Do we have a current node key ?
		if nodeKey != nil {
			nodePath := path.Join(root, nodeKey.Value)

			// Exclude "manala" config
			if nodePath == "/manala" {
				nodeKey = nil
				continue
			}

			var schema map[string]interface{} = nil

			switch nodeChild.Kind {
			case yaml.ScalarNode:
				// Both key/value node are scalars
				schema = map[string]interface{}{}
			case yaml.MappingNode:
				var err error
				schema, err = ld.parseConfigNode(nodeChild, options, nodePath)
				if err != nil {
					return nil, err
				}
			case yaml.SequenceNode:
				schema = map[string]interface{}{
					"type": "array",
				}
			default:
				return nil, fmt.Errorf("unknown node kind: %s", strconv.Itoa(int(nodeChild.Kind)))
			}

			if nodeKey.HeadComment != "" {
				tags := doc.ParseCommentTags(nodeKey.HeadComment)
				// Handle schema tags
				for _, tag := range tags.Filter("schema") {
					var tagSchema map[string]interface{}
					if err := json.Unmarshal([]byte(tag.Value), &tagSchema); err != nil {
						return nil, fmt.Errorf("invalid recipe schema tag at \"%s\": %w", nodePath, err)
					}
					if err := mergo.Merge(&schema, tagSchema, mergo.WithOverride); err != nil {
						return nil, fmt.Errorf("unable to merge recipe schema tag at \"%s\": %w", nodePath, err)
					}
				}
				// Handle option tags
				for _, tag := range tags.Filter("option") {
					option := &models.RecipeOption{
						Path:   nodePath,
						Schema: schema,
					}
					if err := json.Unmarshal([]byte(tag.Value), &option); err != nil {
						return nil, fmt.Errorf("invalid recipe option tag at \"%s\": %w", nodePath, err)
					}
					validate := validator.New()
					if err := validate.Struct(option); err != nil {
						return nil, fmt.Errorf("incorrect recipe option tag at \"%s\": %w", nodePath, err)
					}
					*options = append(*options, *option)
				}
			}

			schemaProperties[nodeKey.Value] = schema

			// Reset node key
			nodeKey = nil
		} else {
			switch nodeChild.Kind {
			case yaml.ScalarNode:
				// Now we have a node key \o/
				nodeKey = nodeChild
			case yaml.MappingNode:
				// This could only be the root node
				schema, err := ld.parseConfigNode(nodeChild, options, "/")
				if err != nil {
					return nil, err
				}
				return schema, nil
			case yaml.SequenceNode:
				// This could only be the root node
				return map[string]interface{}{
					"type": "array",
				}, nil
			default:
				return nil, fmt.Errorf("unknown node kind: %s", strconv.Itoa(int(nodeChild.Kind)))
			}
		}
	}

	// Allow additional properties for empty mappings only
	schemaAdditionalProperties := false
	if node.Content == nil {
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
