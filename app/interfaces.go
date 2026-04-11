package app

import (
	"github.com/manala/manala/internal/path"
	"github.com/manala/manala/internal/schema"
	"github.com/manala/manala/internal/sync"
	"github.com/manala/manala/internal/template"
)

/***********/
/* Project */
/***********/

// Project describe a project interface.
type Project interface {
	Dir() string
	Recipe() Recipe
	Vars() map[string]any
	Template() *template.Template
	Watches() ([]string, error)
}

/**********/
/* Recipe */
/**********/

// Recipe describe a recipe interface.
type Recipe interface {
	Dir() string
	Name() string
	Description() string
	Icon() string
	Vars() map[string]any
	Sync() []sync.UnitInterface
	Schema() schema.Schema
	Options() []RecipeOption
	Repository() Repository
	Template() *template.Template
	ProjectManifestTemplate() *template.Template
	ProjectValidator() *schema.Validator
	Watches() ([]string, error)
}

// RecipeOption describe a recipe option interface.
type RecipeOption interface {
	Name() string
	Label() string
	Help() string
	Path() path.Path
	Schema() schema.Schema
}

/**************/
/* Repository */
/**************/

// Repository describe a repository interface.
type Repository interface {
	URL() string
	Dir() string
}
