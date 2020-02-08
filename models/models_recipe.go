package models

import (
	"github.com/imdario/mergo"
)

// Create a recipe
func NewRecipe(name string, description string, dir string, repository RepositoryInterface) RecipeInterface {
	return &recipe{
		name:        name,
		description: description,
		dir:         dir,
		repository:  repository,
		vars:        map[string]interface{}{},
		syncUnits:   []RecipeSyncUnit{},
		schema:      map[string]interface{}{},
	}
}

type RecipeInterface interface {
	Name() string
	Description() string
	Dir() string
	Repository() RepositoryInterface
	Vars() map[string]interface{}
	MergeVars(vars *map[string]interface{})
	SyncUnits() []RecipeSyncUnit
	AddSyncUnits(units []RecipeSyncUnit)
	Schema() map[string]interface{}
	MergeSchema(schema *map[string]interface{})
}

type recipe struct {
	name        string
	description string
	dir         string
	repository  RepositoryInterface
	vars        map[string]interface{}
	syncUnits   []RecipeSyncUnit
	schema      map[string]interface{}
}

func (rec *recipe) Name() string {
	return rec.name
}

func (rec *recipe) Description() string {
	return rec.description
}

func (rec *recipe) Dir() string {
	return rec.dir
}

func (rec *recipe) Repository() RepositoryInterface {
	return rec.repository
}

func (rec *recipe) Vars() map[string]interface{} {
	return rec.vars
}

func (rec *recipe) MergeVars(vars *map[string]interface{}) {
	_ = mergo.Merge(&rec.vars, vars, mergo.WithOverride)
}

func (rec *recipe) SyncUnits() []RecipeSyncUnit {
	return rec.syncUnits
}

func (rec *recipe) AddSyncUnits(units []RecipeSyncUnit) {
	rec.syncUnits = append(rec.syncUnits, units...)
}

func (rec *recipe) Schema() map[string]interface{} {
	return rec.schema
}

func (rec *recipe) MergeSchema(schema *map[string]interface{}) {
	_ = mergo.Merge(&rec.schema, schema, mergo.WithOverride)
}

type RecipeSyncUnit struct {
	Source      string
	Destination string
}
