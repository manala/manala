package models

import (
	"github.com/imdario/mergo"
)

// Create a project
func NewProject(dir string, recipe RecipeInterface) ProjectInterface {
	return &project{
		dir:    dir,
		recipe: recipe,
		vars:   recipe.Vars(),
	}
}

type ProjectInterface interface {
	Dir() string
	Recipe() RecipeInterface
	Vars() map[string]interface{}
	MergeVars(vars *map[string]interface{})
}

type project struct {
	dir    string
	recipe RecipeInterface
	vars   map[string]interface{}
}

func (prj *project) Dir() string {
	return prj.dir
}

func (prj *project) Recipe() RecipeInterface {
	return prj.recipe
}

func (prj *project) Vars() map[string]interface{} {
	return prj.vars
}

func (prj *project) MergeVars(vars *map[string]interface{}) {
	_ = mergo.Merge(&prj.vars, vars, mergo.WithOverride)
}
