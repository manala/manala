package models

import (
	"github.com/imdario/mergo"
)

const ProjectManifestFile = ".manala.yaml"

// NewProject creates a project
func NewProject(dir string, recipe RecipeInterface, vars map[string]interface{}) ProjectInterface {
	project := &project{
		dir:    dir,
		recipe: recipe,
	}

	// Merge vars
	_ = mergo.Merge(&project.vars, recipe.Vars())
	_ = mergo.Merge(&project.vars, vars, mergo.WithOverride)

	return project
}

type ProjectInterface interface {
	model
	Recipe() RecipeInterface
	Vars() map[string]interface{}
}

type project struct {
	dir    string
	recipe RecipeInterface
	vars   map[string]interface{}
}

func (prj *project) getDir() string {
	return prj.dir
}

func (prj *project) Recipe() RecipeInterface {
	return prj.recipe
}

func (prj *project) Vars() map[string]interface{} {
	return prj.vars
}
