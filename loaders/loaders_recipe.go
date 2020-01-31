package loaders

import (
	"fmt"
	"github.com/apex/log"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"manala/models"
	"manala/yaml/cleaner"
	"os"
	"path"
	"path/filepath"
	"reflect"
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
	var vars map[string]interface{}
	if err := yaml.NewDecoder(cfgFile).Decode(&vars); err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("empty recipe config \"%s\"", cfgFile.Name())
		}
		return nil, fmt.Errorf("invalid recipe config \"%s\" (%w)", cfgFile.Name(), err)
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

	return rec, nil
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
