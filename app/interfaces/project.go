package interfaces

import (
	"io"
	internalTemplate "manala/internal/template"
	internalValidation "manala/internal/validation"
	internalWatcher "manala/internal/watcher"
)

type Project interface {
	Dir() string
	Recipe() Recipe
	Vars() map[string]interface{}
	Template() *internalTemplate.Template
}

type ProjectManifest interface {
	Recipe() string
	Repository() string
	Vars() map[string]interface{}
	ReadFrom(reader io.Reader) error
	internalValidation.Reporter
}

type ProjectManager interface {
	IsProject(dir string) bool
	CreateProject(dir string, rec Recipe, vars map[string]interface{}) (Project, error)
	LoadProject(dir string) (Project, error)
	WatchProject(proj Project, watcher *internalWatcher.Watcher) error
}
