package manifest

import (
	"path/filepath"

	"github.com/manala/manala/app"

	"dario.cat/mergo"
)

type Project struct {
	*Manifest

	dir    string
	recipe app.Recipe
}

func NewProject(dir string, manifest *Manifest, recipe app.Recipe) *Project {
	return &Project{
		dir:      dir,
		Manifest: manifest,
		recipe:   recipe,
	}
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
	_ = mergo.Merge(&vars, project.Manifest.Vars(), mergo.WithOverride)

	return vars
}

func (project *Project) Watches() ([]string, error) {
	return []string{
		filepath.Join(project.Dir(), filename),
	}, nil
}
