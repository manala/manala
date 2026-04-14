package manifest

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/manala/manala/app"
	"github.com/manala/manala/internal/schema"
	"github.com/manala/manala/internal/sync"
)

type Recipe struct {
	manifest   *Manifest
	dir        string
	name       string
	repository app.Repository
}

func NewRecipe(dir, name string, manifest *Manifest, repository app.Repository) *Recipe {
	return &Recipe{
		dir:        dir,
		name:       name,
		manifest:   manifest,
		repository: repository,
	}
}

func (recipe *Recipe) Dir() string {
	return recipe.dir
}

func (recipe *Recipe) Name() string {
	return recipe.name
}

func (recipe *Recipe) Description() string {
	return recipe.manifest.Description
}

func (recipe *Recipe) Icon() string {
	return recipe.manifest.Icon
}

func (recipe *Recipe) Vars() map[string]any {
	return recipe.manifest.Vars()
}

func (recipe *Recipe) Sync() []sync.UnitInterface {
	return recipe.manifest.Sync
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

func (recipe *Recipe) Template() string {
	if recipe.manifest.Template != "" {
		return filepath.Join(recipe.Dir(), recipe.manifest.Template)
	}
	return ""
}

func (recipe *Recipe) Partials() []string {
	var partials []string

	for _, partial := range recipe.manifest.Partials {
		partials = append(partials, filepath.Join(recipe.Dir(), partial))
	}

	// Legacy: if no partials defined, check for _helpers.tmpl
	if len(partials) == 0 {
		helpers := filepath.Join(recipe.Dir(), "_helpers.tmpl")
		if _, err := os.Stat(helpers); err == nil {
			partials = append(partials, helpers)
		}
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
