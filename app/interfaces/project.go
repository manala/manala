package interfaces

import (
	"io"
	"manala/internal/template"
	"manala/internal/validation"
	"manala/internal/watcher"
)

type Project interface {
	Dir() string
	Recipe() Recipe
	Vars() map[string]interface{}
	Template() *template.Template
}

type ProjectManifest interface {
	Recipe() string
	Repository() string
	Vars() map[string]interface{}
	ReadFrom(reader io.Reader) error
	ValidationResultErrorDecorator() validation.ResultErrorDecorator
}

type ProjectManager interface {
	IsProject(dir string) bool
	CreateProject(dir string, rec Recipe, vars map[string]interface{}) (Project, error)
	LoadProject(dir string) (Project, error)
	WatchProject(proj Project, watcher *watcher.Watcher) error
}
