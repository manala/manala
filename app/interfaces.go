package app

import (
	"github.com/manala/manala/app/sync"
)

/***********/
/* Project */
/***********/

// Project describe a project interface.
type Project interface {
	Dir() string
	Recipe() Recipe
	Vars() map[string]any
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
	Template() string
	Partials() []string
	Sync() []sync.Unit
	Repository() Repository
	Vars() map[string]any
	Schema() map[string]any
	Options() []RecipeOption
	Watches() ([]string, error)
}

// RecipeOption describe a recipe option interface.
type RecipeOption interface {
	Name() string
	Label() string
	Help() string
}

/**************/
/* Repository */
/**************/

// Repository describe a repository interface.
type Repository interface {
	URL() string
	Dir() string
}
