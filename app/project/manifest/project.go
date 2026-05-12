package manifest

import (
	"path/filepath"

	"github.com/manala/manala/app"
)

type Project struct {
	dir    string
	recipe app.Recipe
	vars   map[string]any
}

func (project *Project) Dir() string {
	return project.dir
}

func (project *Project) Recipe() app.Recipe {
	return project.recipe
}

func (project *Project) Vars() map[string]any {
	return project.vars
}

func (project *Project) Watches() ([]string, error) {
	return []string{
		filepath.Join(project.Dir(), filename),
	}, nil
}
