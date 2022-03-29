package app

import (
	_ "embed"
	"errors"
	"fmt"
	"manala/loaders"
	"manala/models"
	"manala/validator"
	"os"
	"path/filepath"
)

//go:embed embed/.manala.yaml.tmpl
var recipeTemplate string

func (app *App) Init(
	recipeListApplication func(recipeLoader loaders.RecipeLoaderInterface, repo models.RepositoryInterface) (models.RecipeInterface, error),
	recipeOptionsFormApplication func(rec models.RecipeInterface, vars map[string]interface{}) error,
	dir string,
	recName string,
) error {
	// Ensure directory exists
	if dir != "." {
		stat, err := os.Stat(dir)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				app.log.WithField("dir", dir).Debug("Creating project directory...")
				if err := os.MkdirAll(dir, 0755); err != nil {
					return fmt.Errorf("error creating project directory: %v", err)
				}
				app.log.WithField("dir", dir).Debug("Project directory created")
			} else {
				return fmt.Errorf("error getting project directory stat: %v", err)
			}
		} else if !stat.IsDir() {
			return fmt.Errorf("project directory invalid: %s", dir)
		}
	}

	// Ensure no project already exists
	prjManifest, _ := app.projectLoader.Find(dir, false)
	if prjManifest != nil {
		return fmt.Errorf("project already exists: %s", dir)
	}

	// Load repository
	repo, err := app.repositoryLoader.Load(
		app.config.GetString("repository"),
		app.config.GetString("cache-dir"),
	)
	if err != nil {
		return err
	}

	// Load recipe...
	var rec models.RecipeInterface
	if recName != "" {
		// ...from name if given
		rec, err = app.recipeLoader.Load(recName, repo)
		if err != nil {
			return err
		}
	} else {
		// ...from recipe list
		rec, err = recipeListApplication(app.recipeLoader, repo)
		if err != nil {
			return err
		}
	}

	// Vars
	vars := rec.Vars()

	// Use recipe options form if any
	if len(rec.Options()) > 0 {
		if err := recipeOptionsFormApplication(rec, vars); err != nil {
			return err
		}
	}

	// Template
	template, err := app.templateManager.NewRecipeTemplate(rec)
	if err != nil {
		return err
	}

	if rec.Template() != "" {
		// Load template from recipe
		if err := template.ParseFile(rec.Template()); err != nil {
			return err
		}
	} else {
		// Load default template
		if err := template.Parse(recipeTemplate); err != nil {
			return err
		}
	}

	// Create project manifest
	prjManifest, err = os.Create(filepath.Join(dir, models.ProjectManifestFile))
	if err != nil {
		return err
	}
	defer prjManifest.Close()

	if err := template.Execute(prjManifest, vars); err != nil {
		return err
	}

	prj, err := app.projectLoader.Load(
		prjManifest,
		app.config.GetString("repository"),
		"",
		app.config.GetString("cache-dir"),
	)
	if err != nil {
		return err
	}

	// Validate project
	if err := validator.ValidateProject(prj); err != nil {
		return err
	}

	app.log.Info("Project validated")

	// Sync project
	if err := app.sync.SyncProject(prj); err != nil {
		return err
	}

	app.log.Info("Project synced")

	return nil
}
