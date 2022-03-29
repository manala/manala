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
	recName string,
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
			prjManifest, err := app.projectLoader.Find(path, false)
			if err != nil {
				return err
			}

			// Sync
			if prjManifest != nil {
				if err := app.syncProject(
					prjManifest,
					app.config.GetString("repository"),
					recName,
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
		prjManifest, err := app.projectLoader.Find(dir, true)
		if err != nil {
			return err
		}

		if prjManifest == nil {
			return fmt.Errorf("project not found: %s", dir)
		}

		// Sync
		if err = app.syncProject(
			prjManifest,
			app.config.GetString("repository"),
			recName,
		); err != nil {
			return err
		}
	}

	return nil
}

func (app *App) syncProject(
	prjManifest *os.File,
	defaultRepository string,
	recName string,
) error {
	// Load project
	prj, err := app.projectLoader.Load(
		prjManifest,
		defaultRepository,
		recName,
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
