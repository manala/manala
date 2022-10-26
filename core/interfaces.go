package core

import (
	"io"
	internalSyncer "manala/internal/syncer"
	internalTemplate "manala/internal/template"
	internalValidation "manala/internal/validation"
	internalWatcher "manala/internal/watcher"
)

/***********/
/* Project */
/***********/

type Project interface {
	Path() string
	Recipe() Recipe
	Vars() map[string]interface{}
	Template() *internalTemplate.Template
	Watch(watcher *internalWatcher.Watcher) error
}

type ProjectManifest interface {
	Path() string
	Recipe() string
	Repository() string
	Vars() map[string]interface{}
	ReadFrom(reader io.Reader) error
	internalValidation.Reporter
}

/**********/
/* Recipe */
/**********/

type RecipeManager interface {
	LoadRecipe(name string) (Recipe, error)
}

type RecipeWalker interface {
	WalkRecipes(walker func(rec Recipe)) error
}

type Recipe interface {
	Path() string
	Name() string
	Description() string
	Vars() map[string]interface{}
	Sync() []internalSyncer.UnitInterface
	Schema() map[string]interface{}
	InitVars(callback func(options []RecipeOption) error) (map[string]interface{}, error)
	Repository() Repository
	Template() *internalTemplate.Template
	ProjectManifestTemplate() *internalTemplate.Template
	Watch(watcher *internalWatcher.Watcher) error
}

type RecipeManifest interface {
	Path() string
	Description() string
	Template() string
	Vars() map[string]interface{}
	Sync() []internalSyncer.UnitInterface
	Schema() map[string]interface{}
	ReadFrom(reader io.Reader) error
	internalValidation.Reporter
	InitVars(callback func(options []RecipeOption) error) (map[string]interface{}, error)
}

type RecipeOption interface {
	Label() string
	Schema() map[string]interface{}
	Set(value interface{}) error
}

/**************/
/* Repository */
/**************/

type RepositoryManager interface {
	LoadRepository(paths []string) (Repository, error)
}

type Repository interface {
	Path() string
	// Source keep backward compatibility, when "source" was used instead of "path" to define repository origin
	Source() string
	Dir() string
	LoadRecipe(name string) (Recipe, error)
	WalkRecipes(walker func(rec Recipe)) error
}
