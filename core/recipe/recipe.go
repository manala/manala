package recipe

import (
	"manala/app/interfaces"
	"manala/internal/syncer"
	"manala/internal/template"
	"os"
	"path/filepath"
)

func NewRecipe(dir string, name string, recMan interfaces.RecipeManifest, repo interfaces.Repository) *Recipe {
	return &Recipe{
		dir:        dir,
		name:       name,
		manifest:   recMan,
		repository: repo,
	}
}

type Recipe struct {
	dir        string
	name       string
	manifest   interfaces.RecipeManifest
	repository interfaces.Repository
}

func (rec *Recipe) Dir() string {
	return rec.dir
}

func (rec *Recipe) Name() string {
	return rec.name
}

func (rec *Recipe) Description() string {
	return rec.manifest.Description()
}

func (rec *Recipe) Vars() map[string]interface{} {
	return rec.manifest.Vars()
}

func (rec *Recipe) Sync() []syncer.UnitInterface {
	return rec.manifest.Sync()
}

func (rec *Recipe) Schema() map[string]interface{} {
	return rec.manifest.Schema()
}

func (rec *Recipe) InitVars(callback func(options []interfaces.RecipeOption) error) (map[string]interface{}, error) {
	return rec.manifest.InitVars(callback)
}

func (rec *Recipe) Repository() interfaces.Repository {
	return rec.repository
}

func (rec *Recipe) Template() *template.Template {
	template := template.NewTemplate()

	// Include template helpers if any
	helpersFile := filepath.Join(rec.Dir(), "_helpers.tmpl")
	if _, err := os.Stat(helpersFile); err == nil {
		template.WithDefaultFile(helpersFile)
	}

	return template
}

func (rec *Recipe) ProjectManifestTemplate() *template.Template {
	template := rec.Template()

	if rec.manifest.Template() != "" {
		template.WithFile(filepath.Join(rec.Dir(), rec.manifest.Template()))
	}

	return template
}
