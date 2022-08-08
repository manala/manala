package app

import (
	"errors"
	"github.com/caarlos0/log"
	"github.com/gen2brain/beeep"
	"io/fs"
	"manala/internal"
	internalConfig "manala/internal/config"
	internalFilepath "manala/internal/filepath"
	internalLog "manala/internal/log"
	internalOs "manala/internal/os"
	internalSyncer "manala/internal/syncer"
	internalWatcher "manala/internal/watcher"
	"os"
	"path/filepath"
	"strings"
)

// New creates an app
func New(config *internalConfig.Config, logger *internalLog.Logger) *App {
	// Log
	logger.WithFields(config).Debug("config")

	// App
	app := &App{
		config: config,
		log:    logger,
	}

	// Syncer
	app.syncer = &internalSyncer.Syncer{
		Log: app.log,
	}

	// Watcher manager
	app.watcherManager = &internalWatcher.WatcherManager{
		Log: app.log,
	}

	// Repository manager
	app.repositoryManager = internal.NewRepositoryManager(
		app.log,
		app.config.GetString("default-repository"),
	).AddRepositoryLoader(
		&internal.RepositoryGitLoader{
			Log:      app.log,
			CacheDir: app.config.GetString("cache-dir"),
		},
	).AddRepositoryLoader(
		&internal.RepositoryDirLoader{
			Log: app.log,
		},
	)

	// Project manager
	app.projectManager = internal.NewProjectManager(
		app.log,
		app.repositoryManager,
		&internal.ProjectDirLoader{
			Log:              app.log,
			RepositoryLoader: app.repositoryManager,
		},
	)

	return app
}

type App struct {
	config            *internalConfig.Config
	log               *internalLog.Logger
	syncer            *internalSyncer.Syncer
	watcherManager    *internalWatcher.WatcherManager
	repositoryManager *internal.RepositoryManager
	projectManager    *internal.ProjectManager
}

func (app *App) Repository(path string) (*internal.Repository, error) {
	// Log
	app.log.WithFields(log.Fields{
		"path": path,
	}).Debug("load repository")
	app.log.IncreasePadding()

	repository, err := app.repositoryManager.LoadRepository(
		[]string{path},
	)

	// Log
	app.log.DecreasePadding()

	return repository, err
}

func (app *App) ProjectManifest(path string) (*internal.ProjectManifest, error) {
	return app.projectManager.LoadProjectManifest(path)
}

func (app *App) Project(path string, repositoryPath string, recipeName string) (*internal.Project, error) {
	return app.projectManager.LoadProject(path, repositoryPath, recipeName)
}

func (app *App) ProjectFrom(path string, repositoryPath string, recipeName string) (*internal.Project, error) {
	// Log
	app.log.WithFields(log.Fields{
		"path":       path,
		"repository": repositoryPath,
		"recipe":     recipeName,
	}).Debug("load project from")
	app.log.IncreasePadding()

	var project *internal.Project

	// Backwalks from path
	err := internalFilepath.Backwalk(
		path,
		func(path string, file os.DirEntry, err error) error {
			if err != nil {
				return internalOs.FileSystemError(err)
			}

			// Load project
			project, err = app.Project(path, repositoryPath, recipeName)
			if err != nil {
				var _err *internal.NotFoundProjectManifestError
				if errors.As(err, &_err) {
					err = nil
				}
				return err
			}

			// Stop backwalk
			return filepath.SkipDir
		})

	// Log
	app.log.DecreasePadding()

	if err != nil {
		return nil, err
	}

	if project == nil {
		return nil, internal.NotFoundProjectManifestPathError(path)
	}

	app.log.WithFields(log.Fields{
		"path":       project.Path(),
		"repository": project.Recipe().Repository().Path(),
		"recipe":     project.Recipe().Name(),
	}).Info("project loaded")

	return project, nil
}

func (app *App) WalkProjects(path string, repositoryPath string, recipeName string, projectFunc func(project *internal.Project) error) error {
	// Log
	app.log.WithFields(log.Fields{
		"path": path,
	}).Info("load projects from")
	app.log.IncreasePadding()

	err := filepath.WalkDir(path, func(path string, file os.DirEntry, err error) error {
		if err != nil {
			return internalOs.FileSystemError(err)
		}

		// Only directories
		if !file.IsDir() {
			return nil
		}

		// Exclusions
		if internalFilepath.Exclude(path) {
			app.log.WithFields(log.Fields{
				"path": path,
			}).Debug("exclude project")
			return filepath.SkipDir
		}

		// Log
		app.log.WithFields(log.Fields{
			"path":       path,
			"repository": repositoryPath,
			"recipe":     recipeName,
		}).Debug("load project")
		app.log.IncreasePadding()

		// Load project
		project, err := app.Project(path, repositoryPath, recipeName)

		// Log
		app.log.DecreasePadding()

		if err != nil {
			var _err *internal.NotFoundProjectManifestError
			if errors.As(err, &_err) {
				err = nil
			}
			return err
		}

		app.log.WithFields(log.Fields{
			"path":       project.Path(),
			"repository": project.Recipe().Repository().Path(),
			"recipe":     project.Recipe().Name(),
		}).Info("project loaded")

		// Walk function
		return projectFunc(project)
	})

	// Log
	app.log.DecreasePadding()

	return err
}

func (app *App) SyncProject(project *internal.Project) error {
	// Log
	app.log.IncreasePadding()
	app.log.WithFields(log.Fields{
		"src": project.Recipe().Path(),
		"dst": project.Path(),
	}).Info("sync project")
	app.log.IncreasePadding()

	// Loop over project recipe sync nodes
	for _, node := range project.Recipe().Sync() {
		if err := app.syncer.Sync(
			project.Recipe().Path(),
			node.Source,
			project.Path(),
			node.Destination,
			project,
		); err != nil {
			app.log.DecreasePadding()
			app.log.DecreasePadding()

			return err
		}
	}

	app.log.DecreasePadding()
	app.log.DecreasePadding()

	return nil
}

func (app *App) WatchProject(project *internal.Project, repositoryPath string, recipeName string, all bool, notify bool) error {
	// Log
	app.log.IncreasePadding()
	app.log.WithFields(log.Fields{
		"src": project.Recipe().Path(),
		"dst": project.Path(),
	}).Info("watch project")
	app.log.IncreasePadding()

	path := project.Path()

	watcher, err := app.watcherManager.NewWatcher(

		// On start
		func(watcher *internalWatcher.Watcher) {
			// Watch project manifest
			_ = watcher.Add(project.Manifest().Path())
		},

		// On change
		func(watcher *internalWatcher.Watcher) {
			// Log
			app.log.WithFields(log.Fields{
				"path":       path,
				"repository": repositoryPath,
				"recipe":     recipeName,
			}).Debug("load project")
			app.log.IncreasePadding()

			// Load project
			var err error
			project, err = app.Project(path, repositoryPath, recipeName)

			// Log
			app.log.DecreasePadding()

			if err != nil {
				if notify {
					_ = beeep.Alert("Manala", string(app.log.CaptureError(err)), "")
				}
				app.log.LogError(err)
				return
			}

			// Log
			app.log.WithFields(log.Fields{
				"path":       project.Path(),
				"repository": project.Recipe().Repository().Path(),
				"recipe":     project.Recipe().Name(),
			}).Info("project loaded")

			// Sync project
			err = app.SyncProject(project)

			if err != nil {
				if notify {
					_ = beeep.Alert("Manala", strings.Replace(string(app.log.CaptureError(err)), `"`, `\"`, -1), "")
				}
				app.log.LogError(err)
				return
			}

			if all {
				_ = watcher.RemoveTemporaries()
			}

			if notify {
				_ = beeep.Notify("Manala", "Project synced", "")
			}
		},

		// On all
		func(watcher *internalWatcher.Watcher) {
			if all && project != nil {
				// Watch recipe directories
				_ = filepath.WalkDir(project.Recipe().Path(), func(path string, file fs.DirEntry, err error) error {
					if file.IsDir() {
						_ = watcher.AddTemporary(path)
					}
					return nil
				})
			}
		},
	)
	if err != nil {
		return nil
	}

	defer watcher.Close()
	watcher.Watch()

	app.log.DecreasePadding()
	app.log.DecreasePadding()

	return nil
}
