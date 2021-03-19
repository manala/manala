package loaders

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/imdario/mergo"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
	"io"
	"manala/logger"
	"manala/models"
	"manala/yaml/cleaner"
	"manala/yaml/doc"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

func NewRecipeLoader(log *logger.Logger) RecipeLoaderInterface {
	return &recipeLoader{
		log: log,
	}
}

var recipeConfigFile = ".manala.yaml"

type RecipeLoaderInterface interface {
	Find(dir string) (*os.File, error)
	Load(name string, repository models.RepositoryInterface) (models.RecipeInterface, error)
	Walk(repository models.RepositoryInterface, fn recipeWalkFunc) error
}

type recipeConfig struct {
	Description string `validate:"required"`
	Sync        []models.RecipeSyncUnit
}

type recipeLoader struct {
	log *logger.Logger
}

func (ld *recipeLoader) Find(dir string) (*os.File, error) {
	ld.log.DebugWithField("Searching recipe...", "dir", dir)

	file, err := os.Open(filepath.Join(dir, recipeConfigFile))

	// Return all errors but non existing file ones
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return file, nil
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
	files, err := os.ReadDir(repository.Dir())
	if err != nil {
		return err
	}

	for _, file := range files {
		// Exclude dot files
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}
		if !file.IsDir() {
			continue
		}

		recFile, err := ld.Find(filepath.Join(repository.Dir(), file.Name()))
		if err != nil {
			return err
		}
		if recFile == nil {
			continue
		}

		rec, err := ld.loadDir(file.Name(), recFile, repository)
		if err != nil {
			return err
		}

		fn(rec)
	}

	return nil
}

type recipeWalkFunc func(rec models.RecipeInterface)

func (ld *recipeLoader) loadDir(name string, file *os.File, repository models.RepositoryInterface) (models.RecipeInterface, error) {
	// Get dir
	dir := filepath.Dir(file.Name())

	ld.log.DebugWithField("Loading recipe...", "name", name)

	// Reset file pointer
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	// Parse config file
	node := yaml.Node{}
	if err := yaml.NewDecoder(file).Decode(&node); err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("empty recipe config \"%s\"", file.Name())
		}
		return nil, fmt.Errorf("invalid recipe config \"%s\" (%w)", file.Name(), err)
	}

	var vars map[string]interface{}
	if err := node.Decode(&vars); err != nil {
		return nil, fmt.Errorf("incorrect recipe config \"%s\" (%w)", file.Name(), err)
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

	rec := models.NewRecipe(
		name,
		cfg.Description,
		dir,
		repository,
	)

	// Handle config
	rec.MergeVars(&vars)
	rec.AddSyncUnits(cfg.Sync)

	// Parse config node
	var options []models.RecipeOption
	schema, err := ld.parseConfigNode(&node, &options, "")
	if err != nil {
		return nil, err
	}
	rec.MergeSchema(&schema)
	rec.AddOptions(options)

	return rec, nil
}

func (ld *recipeLoader) parseConfigNode(node *yaml.Node, options *[]models.RecipeOption, path string) (map[string]interface{}, error) {
	var nodeKey *yaml.Node = nil
	schemaProperties := map[string]interface{}{}

	for _, nodeChild := range node.Content {
		// Do we have a current node key ?
		if nodeKey != nil {
			nodePath := filepath.Join(path, nodeKey.Value)

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
