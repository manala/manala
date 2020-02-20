package loaders

import (
	"encoding/json"
	"fmt"
	"github.com/apex/log"
	"github.com/go-playground/validator/v10"
	"github.com/imdario/mergo"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
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

func NewRecipeLoader() RecipeLoaderInterface {
	return &recipeLoader{}
}

var recipeConfigFile = ".manala.yaml"

type RecipeLoaderInterface interface {
	ConfigFile(dir string) (*os.File, error)
	Load(name string, repository models.RepositoryInterface) (models.RecipeInterface, error)
	Walk(repository models.RepositoryInterface, fn recipeWalkFunc) error
}

type recipeConfig struct {
	Description string `validate:"required"`
	Sync        []models.RecipeSyncUnit
}

type recipeLoader struct {
}

func (ld *recipeLoader) ConfigFile(dir string) (*os.File, error) {
	file, err := os.Open(path.Join(dir, recipeConfigFile))
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return nil, fmt.Errorf("open %s: is a directory", file.Name())
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
	files, err := ioutil.ReadDir(repository.Dir())
	if err != nil {
		return err
	}

	for _, file := range files {
		// Exclude dot files
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}
		if file.IsDir() {
			rec, err := ld.loadDir(
				file.Name(),
				filepath.Join(repository.Dir(), file.Name()),
				repository,
			)
			if err != nil {
				return err
			}
			fn(rec)
		}
	}

	return nil
}

type recipeWalkFunc func(rec models.RecipeInterface)

func (ld *recipeLoader) loadDir(name string, dir string, repository models.RepositoryInterface) (models.RecipeInterface, error) {
	// Get config file
	cfgFile, err := ld.ConfigFile(dir)
	if err != nil {
		return nil, err
	}

	log.WithField("name", name).Debug("Loading recipe...")

	// Parse config
	node := yaml.Node{}
	if err := yaml.NewDecoder(cfgFile).Decode(&node); err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("empty recipe config \"%s\"", cfgFile.Name())
		}
		return nil, fmt.Errorf("invalid recipe config \"%s\" (%w)", cfgFile.Name(), err)
	}

	var vars map[string]interface{}
	if err := node.Decode(&vars); err != nil {
		return nil, fmt.Errorf("incorrect recipe config \"%s\" (%w)", cfgFile.Name(), err)
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
