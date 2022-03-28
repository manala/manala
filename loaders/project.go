package loaders

import (
	"errors"
	"fmt"
	"github.com/apex/log"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-yaml"
	"github.com/mitchellh/mapstructure"
	"io"
	"manala/models"
	"os"
	"path/filepath"
)

func NewProjectLoader(log log.Interface, repositoryLoader RepositoryLoaderInterface, recipeLoader RecipeLoaderInterface) ProjectLoaderInterface {
	return &projectLoader{
		log:              log,
		repositoryLoader: repositoryLoader,
		recipeLoader:     recipeLoader,
	}
}

type ProjectLoaderInterface interface {
	Find(dir string, traverse bool) (*os.File, error)
	Load(manifest *os.File, defaultRepository string, withRecipeName string, cacheDir string) (models.ProjectInterface, error)
}

type projectConfig struct {
	Recipe     string `validate:"required"`
	Repository string
}

type projectLoader struct {
	log              log.Interface
	repositoryLoader RepositoryLoaderInterface
	recipeLoader     RecipeLoaderInterface
}

func (ld *projectLoader) Find(dir string, traverse bool) (*os.File, error) {
	ld.log.WithField("dir", dir).Debug("Searching project...")

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
	if !traverse {
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

func (ld *projectLoader) Load(manifest *os.File, defaultRepository string, withRecipeName string, cacheDir string) (models.ProjectInterface, error) {
	// Get dir
	dir := filepath.Dir(manifest.Name())

	ld.log.WithField("dir", dir).Debug("Loading project...")

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
		return nil, fmt.Errorf("incorrect project manifest \"%s\" %s", manifest.Name(), yaml.FormatError(err, true, true))
	}

	// Map config
	cfg := projectConfig{
		Repository: defaultRepository,
	}
	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result: &cfg,
	})
	if err := decoder.Decode(vars["manala"]); err != nil {
		return nil, err
	}

	// Cleanup vars
	delete(vars, "manala")

	// With recipe
	if withRecipeName != "" {
		cfg.Recipe = withRecipeName
	}

	// Validate config
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, err
	}

	ld.log.WithFields(log.Fields{
		"recipe":     cfg.Recipe,
		"repository": cfg.Repository,
	}).Info("Project loaded")

	// Load repository
	repo, err := ld.repositoryLoader.Load(cfg.Repository, cacheDir)
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
