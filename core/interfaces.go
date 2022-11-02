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

/**********/
/* Recipe */
/**********/

type RecipeManager interface {
	LoadRecipe(repo Repository, name string) (Recipe, error)
	WalkRecipes(repo Repository, walker func(rec Recipe) error) error
	WatchRecipe(rec Recipe, watcher *internalWatcher.Watcher) error
}

type Recipe interface {
	Dir() string
	Name() string
	Description() string
	Vars() map[string]interface{}
	Sync() []internalSyncer.UnitInterface
	Schema() map[string]interface{}
	InitVars(callback func(options []RecipeOption) error) (map[string]interface{}, error)
	Repository() Repository
	Template() *internalTemplate.Template
	ProjectManifestTemplate() *internalTemplate.Template
}

type RecipeManifest interface {
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
	LoadRepository(url string) (Repository, error)
}

type Repository interface {
	Url() string
	Dir() string
}
