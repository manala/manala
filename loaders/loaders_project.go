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
	"path/filepath"
)

func NewProjectLoader(repositoryLoader RepositoryLoaderInterface, recipeLoader RecipeLoaderInterface, forceRepositorySrc string, forceRecipe string) ProjectLoaderInterface {
	return &projectLoader{
		repositoryLoader:   repositoryLoader,
		recipeLoader:       recipeLoader,
		forceRepositorySrc: forceRepositorySrc,
		forceRecipe:        forceRecipe,
	}
}

type ProjectLoaderInterface interface {
	Find(dir string, traverse bool) (*os.File, error)
	Load(file *os.File) (models.ProjectInterface, error)
}

var projectConfigFile = ".manala.yaml"

type projectConfig struct {
	Recipe     string `validate:"required"`
	Repository string
}

type projectLoader struct {
	repositoryLoader   RepositoryLoaderInterface
	recipeLoader       RecipeLoaderInterface
	forceRepositorySrc string
	forceRecipe        string
}

func (ld *projectLoader) Find(dir string, traverse bool) (*os.File, error) {
	log.WithField("dir", dir).Debug("Searching project...")

	file, err := os.Open(path.Join(dir, projectConfigFile))

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

	parentDdir := path.Join(dir, "..")
	parentAbs, err := filepath.Abs(parentDdir)
	if err != nil {
		return nil, err
	}

	// If absolute path equals to parent absolute path,
	// we have reached the top of the filesystem
	if abs == parentAbs {
		return nil, nil
	}

	return ld.Find(parentDdir, true)
}

func (ld *projectLoader) Load(file *os.File) (models.ProjectInterface, error) {
	// Get dir
	dir := filepath.Dir(file.Name())

	log.WithField("dir", dir).Debug("Loading project...")

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

	// Force repository
	if ld.forceRepositorySrc != "" {
		cfg.Repository = ld.forceRepositorySrc
	}

	// Force recipe
	if ld.forceRecipe != "" {
		cfg.Recipe = ld.forceRecipe
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
