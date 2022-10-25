package project

import (
	_ "embed"
	"github.com/imdario/mergo"
	"manala/core"
	internalTemplate "manala/internal/template"
	"path/filepath"
)

func NewProject(manifest core.ProjectManifest, rec core.Recipe) *Project {
	return &Project{
		manifest: manifest,
		recipe:   rec,
	}
}

type Project struct {
	manifest core.ProjectManifest
	recipe   core.Recipe
}

func (proj *Project) Path() string {
	return filepath.Dir(proj.manifest.Path())
}

func (proj *Project) Manifest() core.ProjectManifest {
	return proj.manifest
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
