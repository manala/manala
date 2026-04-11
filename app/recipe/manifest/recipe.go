package manifest

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/manala/manala/app"
	"github.com/manala/manala/internal/schema"
	"github.com/manala/manala/internal/template"
)

type Recipe struct {
	*Manifest

	dir        string
	name       string
	repository app.Repository
}

func NewRecipe(dir, name string, manifest *Manifest, repository app.Repository) *Recipe {
	return &Recipe{
		dir:        dir,
		name:       name,
		Manifest:   manifest,
		repository: repository,
	}
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

func (recipe *Recipe) ProjectValidator() *schema.Validator {
	return schema.NewValidator(recipe.Schema())
}

func (recipe *Recipe) Watches() ([]string, error) {
	var dirs []string

	if err := filepath.WalkDir(
		recipe.Dir(),
		func(path string, entry fs.DirEntry, _ error) error {
			if entry.IsDir() {
				dirs = append(dirs, path)
			}

			return nil
		},
	); err != nil {
		return nil, err
	}

	return dirs, nil
}
