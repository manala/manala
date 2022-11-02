package project

import (
	_ "embed"
	"github.com/imdario/mergo"
	"manala/core"
	internalTemplate "manala/internal/template"
)

func NewProject(dir string, projMan core.ProjectManifest, rec core.Recipe) *Project {
	return &Project{
		dir:      dir,
		manifest: projMan,
		recipe:   rec,
	}
}

type Project struct {
	dir      string
	manifest core.ProjectManifest
	recipe   core.Recipe
}

func (proj *Project) Dir() string {
	return proj.dir
}

func (proj *Project) Recipe() core.Recipe {
	return proj.recipe
}

func (proj *Project) Vars() map[string]interface{} {
	var vars map[string]interface{}

	_ = mergo.Merge(&vars, proj.recipe.Vars())
	_ = mergo.Merge(&vars, proj.manifest.Vars(), mergo.WithOverride)

	return vars
}

func (proj *Project) Template() *internalTemplate.Template {
	return proj.recipe.Template().
		WithData(core.NewProjectView(proj))
}
