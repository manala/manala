package manifest

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/manala/manala/app"
	"github.com/manala/manala/internal/schema"
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

func (recipe *Recipe) Template() string {
	if template := recipe.Manifest.Template(); template != "" {
		return filepath.Join(recipe.Dir(), template)
	}
	return ""
}

func (recipe *Recipe) Partials() []string {
	var partials []string

	helpers := filepath.Join(recipe.Dir(), "_helpers.tmpl")
	if _, err := os.Stat(helpers); err == nil {
		partials = append(partials, helpers)
	}

	return partials
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
