package loaders

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
	"io"
	"manala/logger"
	"manala/models"
	"manala/yaml/cleaner"
	"os"
	"path/filepath"
)

func NewProjectLoader(log *logger.Logger, repositoryLoader RepositoryLoaderInterface, recipeLoader RecipeLoaderInterface) ProjectLoaderInterface {
	return &projectLoader{
		log:              log,
		repositoryLoader: repositoryLoader,
		recipeLoader:     recipeLoader,
	}
}

type ProjectLoaderInterface interface {
	Find(dir string, traverse bool) (*os.File, error)
	Load(file *os.File, withRepositorySrc string, withRecipeName string) (models.ProjectInterface, error)
}

var projectConfigFile = ".manala.yaml"

type projectConfig struct {
	Recipe     string `validate:"required"`
	Repository string
}

type projectLoader struct {
	log              *logger.Logger
	repositoryLoader RepositoryLoaderInterface
	recipeLoader     RecipeLoaderInterface
}

func (ld *projectLoader) Find(dir string, traverse bool) (*os.File, error) {
	ld.log.DebugWithField("Searching project...", "dir", dir)

	file, err := os.Open(filepath.Join(dir, projectConfigFile))

	// Return all errors but non existing file ones
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if file != nil || traverse == false {
		return file, nil
	}

	// Traversal mode
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	parentDir := filepath.Join(dir, "..")
	parentAbs, err := filepath.Abs(parentDir)
	if err != nil {
		return nil, err
	}

	// If absolute path equals to parent absolute path,
	// we have reached the top of the filesystem
	if abs == parentAbs {
		return nil, nil
	}

	return ld.Find(parentDir, true)
}

func (ld *projectLoader) Load(file *os.File, withRepositorySrc string, withRecipeName string) (models.ProjectInterface, error) {
	// Get dir
	dir := filepath.Dir(file.Name())

	ld.log.DebugWithField("Loading project...", "dir", dir)

	// Reset file pointer
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	// Parse config file
	var vars map[string]interface{}
	if err := yaml.NewDecoder(file).Decode(&vars); err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("empty project config \"%s\"", file.Name())
		}
		return nil, fmt.Errorf("invalid project config \"%s\" (%w)", file.Name(), err)
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

	// With repository
	if withRepositorySrc != "" {
		cfg.Repository = withRepositorySrc
	}

	// With recipe
	if withRecipeName != "" {
		cfg.Recipe = withRecipeName
	}

	// Cleanup vars
	delete(vars, "manala")

	ld.log.InfoWithFields("Project loaded", logger.Fields{
		"recipe":     cfg.Recipe,
		"repository": cfg.Repository,
	})

	// Load repository
	repo, err := ld.repositoryLoader.Load(cfg.Repository)
	if err != nil {
		return nil, err
	}

	ld.log.Info("Repository loaded")

	rec, err := ld.recipeLoader.Load(cfg.Recipe, repo)
	if err != nil {
		return nil, err
	}

	ld.log.Info("Recipe loaded")

	prj := models.NewProject(
		dir,
		rec,
	)
	prj.MergeVars(&vars)

	return prj, nil
}
