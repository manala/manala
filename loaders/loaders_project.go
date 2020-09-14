package loaders

import (
	"fmt"
	"github.com/apex/log"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
	"io"
	"manala/models"
	"manala/yaml/cleaner"
	"os"
	"path"
)

func NewProjectLoader(repositoryLoader RepositoryLoaderInterface, recipeLoader RecipeLoaderInterface, defaultRepositorySrc string) ProjectLoaderInterface {
	return &projectLoader{
		repositoryLoader:     repositoryLoader,
		recipeLoader:         recipeLoader,
		defaultRepositorySrc: defaultRepositorySrc,
	}
}

type ProjectLoaderInterface interface {
	ConfigFile(dir string) (*os.File, error)
	Load(dir string) (models.ProjectInterface, error)
}

var projectConfigFile = ".manala.yaml"

type projectConfig struct {
	Recipe     string `validate:"required"`
	Repository string
}

type projectLoader struct {
	repositoryLoader     RepositoryLoaderInterface
	recipeLoader         RecipeLoaderInterface
	defaultRepositorySrc string
}

func (ld *projectLoader) ConfigFile(dir string) (*os.File, error) {
	file, err := os.Open(path.Join(dir, projectConfigFile))
	if err != nil {
		return nil, err
	}

	stat, err := file.Stat()
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("\"%s\" file does not exists", file.Name())
		}
		return nil, err
	} else if stat.IsDir() {
		return nil, fmt.Errorf("\"%s\" is not a file", file.Name())
	}

	return file, nil
}

func (ld *projectLoader) Load(dir string) (models.ProjectInterface, error) {
	// Get config file
	cfgFile, err := ld.ConfigFile(dir)
	if err != nil {
		return nil, err
	}

	log.WithField("dir", dir).Debug("Loading project...")

	// Parse config file
	var vars map[string]interface{}
	if err := yaml.NewDecoder(cfgFile).Decode(&vars); err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("empty project config \"%s\"", cfgFile.Name())
		}
		return nil, fmt.Errorf("invalid project config \"%s\" (%w)", cfgFile.Name(), err)
	}

	// See: https://github.com/go-yaml/yaml/issues/139
	vars = cleaner.Clean(vars)

	// Map config
	cfg := projectConfig{}
	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result: &cfg,
	})
	if err := decoder.Decode(vars["manala"]); err != nil {
		return nil, err
	}

	// Validate config
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, err
	}

	// Default repository
	if cfg.Repository == "" {
		cfg.Repository = ld.defaultRepositorySrc
	}

	// Cleanup vars
	delete(vars, "manala")

	log.WithFields(log.Fields{
		"recipe":     cfg.Recipe,
		"repository": cfg.Repository,
	}).Info("Project loaded")

	// Load repository
	repo, err := ld.repositoryLoader.Load(cfg.Repository)
	if err != nil {
		return nil, err
	}

	log.Info("Repository loaded")

	rec, err := ld.recipeLoader.Load(cfg.Recipe, repo)
	if err != nil {
		return nil, err
	}

	log.Info("Recipe loaded")

	prj := models.NewProject(
		dir,
		rec,
	)
	prj.MergeVars(&vars)

	return prj, nil
}
