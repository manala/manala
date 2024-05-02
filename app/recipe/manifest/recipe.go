package manifest

import (
	"manala/app"
	"manala/app/recipe/option"
	"manala/internal/schema"
	"manala/internal/template"
	"manala/internal/validator"
	"os"
	"path/filepath"
)

func NewRecipe(dir string, name string, manifest *Manifest, repository app.Repository) *Recipe {
	return &Recipe{
		dir:        dir,
		name:       name,
		Manifest:   manifest,
		repository: repository,
	}
}

type Recipe struct {
	dir  string
	name string
	*Manifest
	repository app.Repository
}

func (recipe *Recipe) Dir() string {
	return recipe.dir
}

func (recipe *Recipe) Name() string {
	return recipe.name
}

func (recipe *Recipe) Repository() app.Repository {
	return recipe.repository
}

func (recipe *Recipe) Template() *template.Template {
	tmpl := template.NewTemplate()

	// Include template helpers if any
	helpersFile := filepath.Join(recipe.Dir(), "_helpers.tmpl")
	if _, err := os.Stat(helpersFile); err == nil {
		tmpl.WithDefaultFile(helpersFile)
	}

	return tmpl
}

func (recipe *Recipe) ProjectManifestTemplate() *template.Template {
	tmpl := recipe.Template()

	if recipe.Manifest.Template() != "" {
		tmpl.WithFile(filepath.Join(recipe.Dir(), recipe.Manifest.Template()))
	}

	return tmpl
}

func (recipe *Recipe) ProjectValidator() validator.Validator {
	// Option validators
	var optionValidators []validator.Validator
	for _, _option := range recipe.Options() {
		optionValidators = append(optionValidators, option.NewPathedValidator(_option))
	}

	return validator.Validators(
		schema.NewValidator(recipe.Schema()),
		validator.Validators(optionValidators...),
	)
}
