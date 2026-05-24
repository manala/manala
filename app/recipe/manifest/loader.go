package manifest

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/recipe"
	"github.com/manala/manala/internal/errors/serror"
	"github.com/manala/manala/internal/errors/source"
	"github.com/manala/manala/internal/errors/std"
	"github.com/manala/manala/internal/log"
	"github.com/manala/manala/internal/validation"
	yamlerrors "github.com/manala/manala/internal/yaml/errors"
	yamlmapping "github.com/manala/manala/internal/yaml/mapping"
	yamlparser "github.com/manala/manala/internal/yaml/parser"
	yamlvalidation "github.com/manala/manala/internal/yaml/validation"

	"github.com/goccy/go-yaml"
)

const filename = ".manala.yaml"

var manifestValidator = validation.MustNewValidator(map[string]any{
	"type": "object",
	"properties": map[string]any{
		"manala": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"description": map[string]any{"type": "string", "minLength": 1, "maxLength": 256},
				"icon":        map[string]any{"type": "string", "minLength": 1, "maxLength": 100},
				"template":    map[string]any{"type": "string", "minLength": 1, "maxLength": 100},
				"partials": map[string]any{
					"type":  "array",
					"items": map[string]any{"type": "string", "minLength": 1, "maxLength": 100},
				},
				"sync": map[string]any{
					"type":  "array",
					"items": map[string]any{"type": "string", "minLength": 1, "maxLength": 256},
				},
			},
			"additionalProperties": false,
			"required":             []any{"description"},
		},
	},
	"required": []any{"manala"},
})

type LoaderHandler struct {
	log *log.Log
}

func NewLoaderHandler(log *log.Log) *LoaderHandler {
	return &LoaderHandler{
		log: log,
	}
}

func (handler *LoaderHandler) Handle(query *recipe.LoaderQuery, chain recipe.LoaderHandlerChain) (app.Recipe, error) {
	dir := filepath.Join(query.Repository.Dir(), query.Name)
	file := filepath.Join(dir, filename)

	handler.log.Debug("handle recipe manifest", "handler", "manifest", "file", file)

	// Stat file
	if fileInfo, err := os.Stat(file); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Chain
			return chain.Next(query)
		}

		return nil, serror.New("unable to stat recipe manifest").
			With("file", file).
			WithErr(std.From(err))
	} else if fileInfo.IsDir() {
		return nil, serror.New("recipe manifest is a directory").
			With("dir", file)
	}

	// Read file
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, serror.New("unable to read recipe manifest").
			With("file", file).
			WithErr(std.From(err))
	}

	// Init recipe
	recipe := &Recipe{
		dir:        dir,
		name:       query.Name,
		config:     &Config{},
		repository: query.Repository,
	}

	// Prepare source error origin
	origin := source.Origin{
		File:     file,
		Source:   string(content),
		Language: "yaml",
	}

	// Parse content
	node, err := yamlparser.Parse(content)
	if err != nil {
		return nil, serror.New("unable to parse recipe manifest").
			WithErr(source.From(err, origin))
	}

	// Decode manifest map
	var manifestMap map[string]any
	if err := yaml.NodeToValue(node, &manifestMap); err != nil {
		return nil, serror.New("unable to decode recipe manifest").
			WithErr(source.From(yamlerrors.From(err), origin))
	}

	// Validate manifest map
	if err := manifestValidator.Validate(manifestMap, yamlvalidation.WithLocator(node)); err != nil {
		if violations, ok := errors.AsType[validation.Violations](err); ok {
			return nil, serror.New("invalid recipe manifest").
				WithErr(source.From(violations, origin))
		}
		return nil, serror.New("unable to validate recipe manifest").
			With("file", file).WithErr(err)
	}

	// Pop config node
	configNode, found := yamlmapping.Pop(node, "manala")
	if !found {
		return nil, serror.New("recipe manifest config not found").
			With("file", file)
	}

	// Decode config
	if err := yaml.NodeToValue(configNode, recipe.config); err != nil {
		return nil, serror.New("unable to decode recipe manifest config").
			WithErr(source.From(err, origin))
	}

	handler.log.Debug("recipe manifest loaded", "handler", "manifest", "file", file)

	// Decode vars
	if err := yaml.NodeToValue(node, &recipe.vars); err != nil {
		return nil, serror.New("unable to decode recipe manifest vars").
			WithErr(source.From(yamlerrors.From(err), origin))
	}

	// Infer schema & options
	inf := Inferrer{
		Schema:  &recipe.schema,
		Options: &recipe.options,
	}
	if err = inf.Infer(node); err != nil {
		return nil, serror.New("unable to infer recipe manifest vars").
			WithErr(source.From(err, origin))
	}

	return recipe, nil
}
