package app

import (
	"fmt"
	"github.com/gen2brain/beeep"
	"manala/models"
	"manala/validator"
	"os"
	"strings"
)

func (app *App) Watch(
	dir string,
	withRecipeName string,
	watchAll bool,
	useNotify bool,
) error {
	// Check directory
	if dir != "." {
		if _, err := os.Stat(dir); err != nil {
			return fmt.Errorf("invalid directory: %s", dir)
		}
	}

	// Find project manifest
	prjManifest, err := app.projectLoader.Find(dir, true)
	if err != nil {
		return err
	}

	if prjManifest == nil {
		return fmt.Errorf("project not found: %s", dir)
	}

	// Sync function
	syncFunc := app.getSyncFunc(
		prjManifest,
		app.config.GetString("repository"),
		withRecipeName,
		app.config.GetString("cache-dir"),
		watchAll,
	)

	// Watcher
	watcher, err := app.watcherManager.NewWatcher()
	if err != nil {
		return fmt.Errorf("error creating watcher: %v", err)
	}
	defer watcher.Close()

	// Sync
	if err := syncFunc(watcher); err != nil {
		return err
	}

	app.log.Info("Start watching...")

	// Watch
	watcher.Watch(func(watcher models.WatcherInterface) {
		if err := syncFunc(watcher); err != nil {
			app.log.Error(err.Error())
			if useNotify {
				_ = beeep.Alert("Manala", strings.Replace(err.Error(), `"`, `\"`, -1), "")
			}
		} else {
			if useNotify {
				_ = beeep.Notify("Manala", "Project synced", "")
			}
		}
	})

	return nil
}

func (app *App) getSyncFunc(
	prjManifest *os.File,
	defaultRepository string,
	withRecipeName string,
	cacheDir string,
	watchAll bool,
) func(watcher models.WatcherInterface) error {
	return func(watcher models.WatcherInterface) error {
		// Load project
		prj, err := app.projectLoader.Load(prjManifest, defaultRepository, withRecipeName, cacheDir)
		if err != nil {
			return err
		}

		// Validate project
		if err := validator.ValidateProject(prj); err != nil {
			return err
		}

		app.log.Info("Project validated")

		// Watch project
		if err := watcher.SetProject(prj); err != nil {
			return fmt.Errorf("error setting project watching: %v", err)
		}

		// Watch recipe
		if watchAll {
			if err := watcher.SetRecipe(prj.Recipe()); err != nil {
				return fmt.Errorf("error setting recipe watching: %v", err)
			}
		}

		// Sync project
		if err := app.sync.SyncProject(prj); err != nil {
			return err
		}

		app.log.Info("Project synced")

		return nil
	}
}
