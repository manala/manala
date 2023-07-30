package interfaces

import (
	"io"
	"manala/internal/syncer"
	"manala/internal/template"
	"manala/internal/validation"
	"manala/internal/watcher"
)

type Recipe interface {
	Dir() string
	Name() string
	Description() string
	Vars() map[string]interface{}
	Sync() []syncer.UnitInterface
	Schema() map[string]interface{}
	InitVars(callback func(options []RecipeOption) error) (map[string]interface{}, error)
	Repository() Repository
	Template() *template.Template
	ProjectManifestTemplate() *template.Template
}

type RecipeManifest interface {
	Description() string
	Template() string
	Vars() map[string]interface{}
	Sync() []syncer.UnitInterface
	Schema() map[string]interface{}
	ReadFrom(reader io.Reader) error
	ValidationResultErrorDecorator() validation.ResultErrorDecorator
	InitVars(callback func(options []RecipeOption) error) (map[string]interface{}, error)
}

type RecipeOption interface {
	Label() string
	Schema() map[string]interface{}
	Set(value interface{}) error
}

type RecipeManager interface {
	LoadRecipe(repo Repository, name string) (Recipe, error)
	WalkRecipes(repo Repository, walker func(rec Recipe) error) error
	WatchRecipe(rec Recipe, watcher *watcher.Watcher) error
}
