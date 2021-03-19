package loaders

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
	"io"
	"manala/config"
	"manala/logger"
	"manala/models"
	"manala/yaml/cleaner"
	"os"
	"path/filepath"
)

func NewProjectLoader(log *logger.Logger, conf *config.Config, repositoryLoader RepositoryLoaderInterface, recipeLoader RecipeLoaderInterface) ProjectLoaderInterface {
	return &projectLoader{
		log:              log,
		conf:             conf,
		repositoryLoader: repositoryLoader,
		recipeLoader:     recipeLoader,
	}
}

type ProjectLoaderInterface interface {
	Find(dir string, traverse bool) (*os.File, error)
	Load(manifest *os.File, withRepositorySource string, withRecipeName string) (models.ProjectInterface, error)
}

type projectConfig struct {
	Recipe     string `validate:"required"`
	Repository string
}

type projectLoader struct {
	log              *logger.Logger
	conf             *config.Config
	repositoryLoader RepositoryLoaderInterface
	recipeLoader     RecipeLoaderInterface
}

func (ld *projectLoader) Find(dir string, traverse bool) (*os.File, error) {
	ld.log.DebugWithField("Searching project...", "dir", dir)

	manifest, err := os.Open(filepath.Join(dir, models.ProjectManifestFile))

	// Found manifest without errors, return it !
	if err == nil {
		return manifest, nil
	}

	// Encounter serious error, return it !
	if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	// Not found manifest...
	if traverse == false {
		return nil, nil
	}

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

func (ld *projectLoader) Load(manifest *os.File, withRepositorySource string, withRecipeName string) (models.ProjectInterface, error) {
	// Get dir
	dir := filepath.Dir(manifest.Name())

	ld.log.DebugWithField("Loading project...", "dir", dir)

	// Reset manifest pointer
	_, err := manifest.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	// Parse manifest
	var vars map[string]interface{}
	if err := yaml.NewDecoder(manifest).Decode(&vars); err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("empty project manifest \"%s\"", manifest.Name())
		}
		return nil, fmt.Errorf("invalid project manifest \"%s\" (%w)", manifest.Name(), err)
	}

	// See: https://github.com/go-yaml/yaml/issues/139
	vars = cleaner.Clean(vars)

	// Map config
	cfg := projectConfig{
		Repository: ld.conf.Repository(),
	}
	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result: &cfg,
	})
	if err := decoder.Decode(vars["manala"]); err != nil {
		return nil, err
	}

	// Cleanup vars
	delete(vars, "manala")

	// With repository
	if withRepositorySource != "" {
		cfg.Repository = withRepositorySource
	}

	// With recipe
	if withRecipeName != "" {
		cfg.Recipe = withRecipeName
	}

	// Validate config
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, err
	}

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

	return models.NewProject(
		dir,
		rec,
		vars,
	), nil
}
