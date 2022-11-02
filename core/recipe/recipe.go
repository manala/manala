package recipe

import (
	"manala/core"
	internalSyncer "manala/internal/syncer"
	internalTemplate "manala/internal/template"
	"os"
	"path/filepath"
)

func NewRecipe(dir string, name string, recMan core.RecipeManifest, repo core.Repository) *Recipe {
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
	manifest   core.RecipeManifest
	repository core.Repository
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

func (rec *Recipe) Sync() []internalSyncer.UnitInterface {
	return rec.manifest.Sync()
}

func (rec *Recipe) Schema() map[string]interface{} {
	return rec.manifest.Schema()
}

func (rec *Recipe) InitVars(callback func(options []core.RecipeOption) error) (map[string]interface{}, error) {
	return rec.manifest.InitVars(callback)
}

func (rec *Recipe) Repository() core.Repository {
	return rec.repository
}

func (rec *Recipe) Template() *internalTemplate.Template {
	template := internalTemplate.NewTemplate()

	// Include template helpers if any
	helpersFile := filepath.Join(rec.Dir(), "_helpers.tmpl")
	if _, err := os.Stat(helpersFile); err == nil {
		template.WithDefaultFile(helpersFile)
	}

	return template
}

func (rec *Recipe) ProjectManifestTemplate() *internalTemplate.Template {
	template := rec.Template()

	if rec.manifest.Template() != "" {
		template.WithFile(filepath.Join(rec.Dir(), rec.manifest.Template()))
	}

	return template
}
