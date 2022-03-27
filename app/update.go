package app

import (
	"errors"
	"fmt"
	"manala/validator"
	"os"
	"path/filepath"
	"strings"
)

func (app *App) Update(
	dir string,
	withRepositorySource string,
	withRecipeName string,
	recursive bool,
) error {
	// Check directory
	if dir != "." {
		if _, err := os.Stat(dir); err != nil {
			return fmt.Errorf("invalid directory: %s", dir)
		}
	}

	if recursive {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			// Only directories
			if err != nil || !info.IsDir() {
				return err
			}

			// Only not dotted directories
			// (except - of course - current directory)
			if strings.HasPrefix(filepath.Base(path), ".") && (path != ".") {
				return filepath.SkipDir
			}

			// Find project manifest
			prjManifest, err := app.ProjectLoader.Find(path, false)
			if err != nil {
				return err
			}

			// Sync
			if prjManifest != nil {
				if err := app.syncProject(
					prjManifest,
					withRepositorySource,
					withRecipeName,
				); err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	} else {
		// Find project manifest
		prjManifest, err := app.ProjectLoader.Find(dir, true)
		if err != nil {
			return err
		}

		if prjManifest == nil {
			return fmt.Errorf("project not found: %s", dir)
		}

		// Sync
		if err = app.syncProject(
			prjManifest,
			withRepositorySource,
			withRecipeName,
		); err != nil {
			return err
		}
	}

	return nil
}

func (app *App) syncProject(
	prjManifest *os.File,
	withRepositorySource string,
	withRecipeName string,
) error {
	// Load project
	prj, err := app.ProjectLoader.Load(prjManifest, withRepositorySource, withRecipeName)
	if err != nil {
		return err
	}

	// Validate project
	if err := validator.ValidateProject(prj); err != nil {
		return err
	}

	app.Log.Info("Project validated")

	// Sync project
	if err := app.Sync.SyncProject(prj); err != nil {
		return err
	}

	app.Log.Info("Project synced")

	return nil
}
