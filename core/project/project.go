package project

import (
	"dario.cat/mergo"
	_ "embed"
	"manala/app/interfaces"
	"manala/app/views"
	"manala/internal/template"
)

func NewProject(dir string, projMan interfaces.ProjectManifest, rec interfaces.Recipe) *Project {
	return &Project{
		dir:      dir,
		manifest: projMan,
		recipe:   rec,
	}
}

type Project struct {
	dir      string
	manifest interfaces.ProjectManifest
	recipe   interfaces.Recipe
}

func (proj *Project) Dir() string {
	return proj.dir
}

func (proj *Project) Recipe() interfaces.Recipe {
	return proj.recipe
}

func (proj *Project) Vars() map[string]interface{} {
	var vars map[string]interface{}

	_ = mergo.Merge(&vars, proj.recipe.Vars())
	_ = mergo.Merge(&vars, proj.manifest.Vars(), mergo.WithOverride)

	return vars
}

func (proj *Project) Template() *template.Template {
	return proj.recipe.Template().
		WithData(views.NormalizeProject(proj))
}
