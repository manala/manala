package models

import (
	"path/filepath"
)

const RecipeManifestFile = ".manala.yaml"

// Create a recipe
func NewRecipe(name string, description string, template string, dir string, repository RepositoryInterface, vars map[string]interface{}, sync []RecipeSyncUnit, schema map[string]interface{}, options []RecipeOption) *recipe {
	return &recipe{
		name:        name,
		description: description,
		template:    template,
		dir:         dir,
		repository:  repository,
		vars:        vars,
		sync:        sync,
		schema:      schema,
		options:     options,
	}
}

type RecipeInterface interface {
	model
	Name() string
	Description() string
	Template() string
	Repository() RepositoryInterface
	Vars() map[string]interface{}
	Sync() []RecipeSyncUnit
	Schema() map[string]interface{}
	Options() []RecipeOption
}

type recipe struct {
	name        string
	description string
	template    string
	dir         string
	repository  RepositoryInterface
	vars        map[string]interface{}
	sync        []RecipeSyncUnit
	schema      map[string]interface{}
	options     []RecipeOption
}

func (rec *recipe) Name() string {
	return rec.name
}

func (rec *recipe) Description() string {
	return rec.description
}

func (rec *recipe) Template() string {
	return rec.template
}

func (rec *recipe) getDir() string {
	return filepath.Join(rec.Repository().getDir(), rec.dir)
}

func (rec *recipe) Repository() RepositoryInterface {
	return rec.repository
}

func (rec *recipe) Vars() map[string]interface{} {
	return rec.vars
}

func (rec *recipe) Sync() []RecipeSyncUnit {
	return rec.sync
}

func (rec *recipe) Schema() map[string]interface{} {
	return rec.schema
}

func (rec *recipe) Options() []RecipeOption {
	return rec.options
}

type RecipeSyncUnit struct {
	Source      string
	Destination string
}

type RecipeOption struct {
	Label  string                 `json:"label" validate:"required"`
	Path   string                 `json:"path"`
	Schema map[string]interface{} `json:"schema"`
}
