package manifest

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/project"
	"github.com/manala/manala/app/recipe"
	"github.com/manala/manala/app/repository"
	"github.com/manala/manala/internal/errors/serror"
	"github.com/manala/manala/internal/errors/source"
	"github.com/manala/manala/internal/errors/std"
	"github.com/manala/manala/internal/filepath/backwalk"
	"github.com/manala/manala/internal/log"
	"github.com/manala/manala/internal/validation"
	yamlerrors "github.com/manala/manala/internal/yaml/errors"
	yamlmapping "github.com/manala/manala/internal/yaml/mapping"
	yamlparser "github.com/manala/manala/internal/yaml/parser"
	yamlvalidation "github.com/manala/manala/internal/yaml/validation"

	"dario.cat/mergo"
	"github.com/goccy/go-yaml"
)

const filename = ".manala.yaml"

type LoaderHandler struct {
	log              *log.Log
	repositoryLoader *repository.Loader
	recipeLoader     *recipe.Loader
}

func NewLoaderHandler(log *log.Log, repositoryLoader *repository.Loader, recipeLoader *recipe.Loader) *LoaderHandler {
	return &LoaderHandler{
		log:              log,
		repositoryLoader: repositoryLoader,
		recipeLoader:     recipeLoader,
	}
}

func (handler *LoaderHandler) Handle(query *project.LoaderQuery, chain project.LoaderHandlerChain) (app.Project, error) {
	dir := query.Dir
	file := filepath.Join(dir, filename)

	handler.log.Debug("handle project manifest", "handler", "manifest", "file", file)

	// Stat file
	if fileInfo, err := os.Stat(file); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Chain
			return chain.Next(query)
		}

		return nil, serror.New("unable to stat project manifest").
			With("file", file).
			WithErr(std.From(err))
	} else if fileInfo.IsDir() {
		return nil, serror.New("project manifest is a directory").
			With("dir", file)
	}

	// Read file
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, serror.New("unable to read project manifest").
			With("file", file).
			WithErr(std.From(err))
	}

	// Prepare source error origin
	origin := source.Origin{
		File:     file,
		Source:   string(content),
		Language: "yaml",
	}

	// Init project
	project := &Project{
		dir: dir,
	}

	// Parse content
	node, err := yamlparser.Parse(content)
	if err != nil {
		return nil, serror.New("unable to parse project manifest").
			WithErr(source.From(err, origin))
	}

	// Pop config node
	configNode, found := yamlmapping.Pop(node, "manala")
	if !found {
		return nil, serror.New("invalid project manifest").
			WithErr(source.From(yamlerrors.New(
				errors.New("missing \"manala\" property"),
				node.GetToken(),
			), origin))
	}

	// Decode config
	config := &Config{}
	if err := yaml.NodeToValue(configNode, config); err != nil {
		return nil, serror.New("unable to decode project manifest config").
			WithErr(source.From(err, origin))
	}

	handler.log.Debug("project manifest loaded", "handler", "manifest",
		"file", file,
		"repository", config.Repository,
		"recipe", config.Recipe,
	)

	// Load repository
	repository, err := handler.repositoryLoader.Load(config.Repository)
	if err != nil {
		return nil, err
	}

	// Load recipe
	project.recipe, err = handler.recipeLoader.Load(repository, config.Recipe)
	if err != nil {
		return nil, err
	}

	// Decode vars
	var vars map[string]any
	if err := yaml.NodeToValue(node, &vars); err != nil {
		return nil, serror.New("unable to decode project manifest vars").
			WithErr(source.From(yamlerrors.From(err), origin))
	}

	// Merge vars
	_ = mergo.Merge(&project.vars, project.recipe.Vars())
	_ = mergo.Merge(&project.vars, vars, mergo.WithOverride)

	// Validate vars
	validator, err := validation.NewValidator(project.recipe.Schema())
	if err != nil {
		return nil, err
	}
	if violations, err := validator.Validate(project.vars, yamlvalidation.WithLocator(node)); err != nil {
		return nil, serror.New("unable to validate project manifest vars").
			With("file", file).WithErr(err)
	} else if violations != nil {
		return nil, serror.New("invalid project manifest vars").
			WithErr(source.From(violations, origin))
	}

	return project, nil
}

type FromLoaderHandler struct {
	log *log.Log
}

func NewFromLoaderHandler(log *log.Log) *FromLoaderHandler {
	return &FromLoaderHandler{
		log: log,
	}
}

func (handler *FromLoaderHandler) Handle(query *project.LoaderQuery, chain project.LoaderHandlerChain) (app.Project, error) {
	dir := query.Dir

	handler.log.Debug("handle project manifest from", "handler", "manifest.from", "dir", dir)

	var project app.Project

	// Backwalk from dir
	if err := backwalk.WalkDir(dir,
		func(path string, _ os.DirEntry, err error) error {
			if err != nil {
				return serror.New("file system error").
					With("path", path).
					WithErr(std.From(err))
			}

			// Update query
			query.Dir = path

			// Load project
			project, err = chain.Next(query)
			if err != nil {
				if _, ok := errors.AsType[*app.NotFoundProjectError](err); ok {
					return nil
				}

				return err
			}

			// Stop backwalk
			return filepath.SkipAll
		},
	); err != nil {
		return nil, err
	}

	if project == nil {
		query.Dir = dir

		return chain.Last(query)
	}

	return project, nil
}
