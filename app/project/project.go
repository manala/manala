package project

import (
	"dario.cat/mergo"
	_ "embed"
	"manala/app"
	"manala/internal/template"
)

func New(dir string, manifest app.ProjectManifest, recipe app.Recipe) *Project {
	return &Project{
		dir:      dir,
		manifest: manifest,
		recipe:   recipe,
	}
}

type Project struct {
	dir      string
	manifest app.ProjectManifest
	recipe   app.Recipe
}

func (project *Project) Dir() string {
	return project.dir
}

func (project *Project) Recipe() app.Recipe {
	return project.recipe
}

func (project *Project) Vars() map[string]any {
	var vars map[string]any

	_ = mergo.Merge(&vars, project.recipe.Vars())
	_ = mergo.Merge(&vars, project.manifest.Vars(), mergo.WithOverride)

	return vars
}

func (project *Project) Template() *template.Template {
	return project.recipe.Template().
		WithData(app.NewProjectView(project))
}
