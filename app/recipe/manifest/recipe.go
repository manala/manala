package manifest

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/sync"
)

type Recipe struct {
	dir        string
	name       string
	config     *Config
	repository app.Repository
	vars       map[string]any
	schema     map[string]any
	options    []app.RecipeOption
}

func (recipe *Recipe) Dir() string {
	return recipe.dir
}

func (recipe *Recipe) Name() string {
	return recipe.name
}

func (recipe *Recipe) Description() string {
	return recipe.config.Description
}

func (recipe *Recipe) Icon() string {
	return recipe.config.Icon
}

func (recipe *Recipe) Template() string {
	if recipe.config.Template != "" {
		return filepath.Join(recipe.Dir(), recipe.config.Template)
	}
	return ""
}

func (recipe *Recipe) Partials() []string {
	var partials []string

	for _, partial := range recipe.config.Partials {
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

func (recipe *Recipe) Sync() []sync.Unit {
	return recipe.config.Sync
}

func (recipe *Recipe) Repository() app.Repository {
	return recipe.repository
}

func (recipe *Recipe) Vars() map[string]any {
	return recipe.vars
}

func (recipe *Recipe) Schema() map[string]any {
	return recipe.schema
}

func (recipe *Recipe) Options() []app.RecipeOption {
	return recipe.options
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
