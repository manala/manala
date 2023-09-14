package recipe

import (
	"manala/app"
	option "manala/app/recipe/option"
	"manala/internal/schema"
	"manala/internal/syncer"
	"manala/internal/template"
	"manala/internal/validator"
	"os"
	"path/filepath"
)

func New(dir string, name string, manifest app.RecipeManifest, repository app.Repository) *Recipe {
	return &Recipe{
		dir:        dir,
		name:       name,
		manifest:   manifest,
		repository: repository,
	}
}

type Recipe struct {
	dir        string
	name       string
	manifest   app.RecipeManifest
	repository app.Repository
}

func (recipe *Recipe) Dir() string {
	return recipe.dir
}

func (recipe *Recipe) Name() string {
	return recipe.name
}

func (recipe *Recipe) Description() string {
	return recipe.manifest.Description()
}

func (recipe *Recipe) Vars() map[string]any {
	return recipe.manifest.Vars()
}

func (recipe *Recipe) Sync() []syncer.UnitInterface {
	return recipe.manifest.Sync()
}

func (recipe *Recipe) Schema() schema.Schema {
	return recipe.manifest.Schema()
}

func (recipe *Recipe) Options() []app.RecipeOption {
	return recipe.manifest.Options()
}

func (recipe *Recipe) Repository() app.Repository {
	return recipe.repository
}

func (recipe *Recipe) Template() *template.Template {
	template := template.NewTemplate()

	// Include template helpers if any
	helpersFile := filepath.Join(recipe.Dir(), "_helpers.tmpl")
	if _, err := os.Stat(helpersFile); err == nil {
		template.WithDefaultFile(helpersFile)
	}

	return template
}

func (recipe *Recipe) ProjectManifestTemplate() *template.Template {
	template := recipe.Template()

	if recipe.manifest.Template() != "" {
		template.WithFile(filepath.Join(recipe.Dir(), recipe.manifest.Template()))
	}

	return template
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
