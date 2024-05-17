package app

import (
	"manala/internal/path"
	"manala/internal/schema"
	"manala/internal/syncer"
	"manala/internal/template"
	"manala/internal/validator"
)

/***********/
/* Project */
/***********/

// Project describe a project interface
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

// Recipe describe a recipe interface
type Recipe interface {
	Dir() string
	Name() string
	Description() string
	Icon() string
	Vars() map[string]any
	Sync() []syncer.UnitInterface
	Schema() schema.Schema
	Options() []RecipeOption
	Repository() Repository
	Template() *template.Template
	ProjectManifestTemplate() *template.Template
	ProjectValidator() validator.Validator
	Watches() ([]string, error)
}

// RecipeOption describe a recipe option interface
type RecipeOption interface {
	Name() string
	Label() string
	Help() string
	Path() path.Path
	Schema() schema.Schema
	Validate(value any) (validator.Violations, error)
}

/**************/
/* Repository */
/**************/

// Repository describe a repository interface
type Repository interface {
	Url() string
	Dir() string
}
