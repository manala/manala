package application

import (
	"errors"
	"github.com/caarlos0/log"
	"github.com/gen2brain/beeep"
	"manala/core"
	"manala/core/project"
	"manala/core/repository"
	internalConfig "manala/internal/config"
	internalFilepath "manala/internal/filepath"
	internalLog "manala/internal/log"
	internalOs "manala/internal/os"
	internalReport "manala/internal/report"
	internalSyncer "manala/internal/syncer"
	internalWatcher "manala/internal/watcher"
	"os"
	"path/filepath"
)

// NewApplication creates an application
func NewApplication(config *internalConfig.Config, log *internalLog.Logger) *Application {
	// Log
	log.WithFields(config).Debug("config")

	// App
	app := &Application{
		config: config,
		log:    log,
	}

	// Syncer
	app.syncer = &internalSyncer.Syncer{
		Log: app.log,
	}

	// Watcher manager
	app.watcherManager = &internalWatcher.Manager{
		Log: app.log,
	}

	// Repository manager
	app.repositoryManager = repository.NewChainManager(
		app.log,
		app.config.GetString("default-repository"),
		[]core.RepositoryManager{
			repository.NewGitManager(
				app.log,
				app.config.GetString("cache-dir"),
			),
			repository.NewDirManager(
				app.log,
			),
		},
	)

	// Project manager
	app.projectManager = project.NewManager(
		app.log,
		app.repositoryManager,
	)

	return app
}

type Application struct {
	config            *internalConfig.Config
	log               *internalLog.Logger
	syncer            *internalSyncer.Syncer
	watcherManager    *internalWatcher.Manager
	repositoryManager core.RepositoryManager
	projectManager    *project.Manager
}

func (app *Application) Repository(path string) (core.Repository, error) {
	// Log
	app.log.WithFields(log.Fields{
		"path": path,
	}).Debug("load repository")
	app.log.IncreasePadding()

	repo, err := app.repositoryManager.LoadRepository(
		[]string{path},
	)

	// Log
	app.log.DecreasePadding()

	return repo, err
}

func (app *Application) ProjectManifest(path string) (core.ProjectManifest, error) {
	return app.projectManager.LoadProjectManifest(path)
}

func (app *Application) CreateProject(
	path string,
	repo core.Repository,
	recSelector func(recipeWalker core.RecipeWalker) (core.Recipe, error),
	optionsSelector func(rec core.Recipe, options []core.RecipeOption) error,
) (core.Project, error) {
	// Select recipe
	rec, err := recSelector(repo)
	if err != nil {
		return nil, err
	}

	// Select options to get init vars
	vars, err := rec.InitVars(func(options []core.RecipeOption) error {
		if err := optionsSelector(rec, options); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return app.projectManager.CreateProject(path, rec, vars)
}

func (app *Application) Project(path string, repoPath string, recName string) (core.Project, error) {
	return app.projectManager.LoadProject(path, repoPath, recName)
}

func (app *Application) ProjectFrom(path string, repoPath string, recName string) (core.Project, error) {
	// Log
	app.log.WithFields(log.Fields{
		"path":       path,
		"repository": repoPath,
		"recipe":     recName,
	}).Debug("load project from")
	app.log.IncreasePadding()

	var proj core.Project

	// Backwalks from path
	err := internalFilepath.Backwalk(
		path,
		func(path string, file os.DirEntry, err error) error {
			if err != nil {
				return internalReport.NewError(internalOs.NewError(err)).
					WithMessage("file system error")
			}

			// Load project
			proj, err = app.Project(path, repoPath, recName)
			if err != nil {
				var _notFoundProjectManifestError *core.NotFoundProjectManifestError
				if errors.As(err, &_notFoundProjectManifestError) {
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

	if proj == nil {
		return nil, internalReport.NewError(
			core.NewNotFoundProjectManifestError("project manifest not found"),
		).WithField("path", path)
	}

	app.log.WithFields(log.Fields{
		"path":       proj.Path(),
		"repository": proj.Recipe().Repository().Path(),
		"recipe":     proj.Recipe().Name(),
	}).Info("project loaded")

	return proj, nil
}

func (app *Application) WalkProjects(path string, repoPath string, recName string, walker func(proj core.Project) error) error {
	// Log
	app.log.WithFields(log.Fields{
		"path": path,
	}).Info("load projects from")
	app.log.IncreasePadding()

	err := filepath.WalkDir(path, func(path string, file os.DirEntry, err error) error {
		if err != nil {
			return internalReport.NewError(internalOs.NewError(err)).
				WithMessage("file system error")
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
			"repository": repoPath,
			"recipe":     recName,
		}).Debug("load project")
		app.log.IncreasePadding()

		// Load project
		proj, err := app.Project(path, repoPath, recName)

		// Log
		app.log.DecreasePadding()

		if err != nil {
			var _notFoundProjectManifestError *core.NotFoundProjectManifestError
			if errors.As(err, &_notFoundProjectManifestError) {
				err = nil
			}
			return err
		}

		app.log.WithFields(log.Fields{
			"path":       proj.Path(),
			"repository": proj.Recipe().Repository().Path(),
			"recipe":     proj.Recipe().Name(),
		}).Info("project loaded")

		// Walk function
		return walker(proj)
	})

	// Log
	app.log.DecreasePadding()

	return err
}

func (app *Application) SyncProject(proj core.Project) error {
	// Log
	app.log.IncreasePadding()
	app.log.WithFields(log.Fields{
		"src": proj.Recipe().Path(),
		"dst": proj.Path(),
	}).Info("sync project")
	app.log.IncreasePadding()

	// Loop over project recipe sync units
	for _, unit := range proj.Recipe().Sync() {
		if err := app.syncer.Sync(
			proj.Recipe().Path(),
			unit.Source(),
			proj.Path(),
			unit.Destination(),
			proj,
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

func (app *Application) WatchProject(proj core.Project, repoPath string, recName string, all bool, notify bool) error {
	// Log
	app.log.IncreasePadding()
	app.log.WithFields(log.Fields{
		"src": proj.Recipe().Path(),
		"dst": proj.Path(),
	}).Info("watch project")
	app.log.IncreasePadding()

	path := proj.Path()

	watcher, err := app.watcherManager.NewWatcher(

		// On start
		func(watcher *internalWatcher.Watcher) {
			// Watch project
			_ = proj.Watch(watcher)
		},

		// On change
		func(watcher *internalWatcher.Watcher) {
			// Log
			app.log.WithFields(log.Fields{
				"path":       path,
				"repository": repoPath,
				"recipe":     recName,
			}).Debug("load project")
			app.log.IncreasePadding()

			// Load project
			var err error
			proj, err = app.Project(path, repoPath, recName)

			// Log
			app.log.DecreasePadding()

			if err != nil {
				report := internalReport.NewErrorReport(err)
				if notify {
					_ = beeep.Alert("Manala", report.String(), "")
				}
				app.log.Report(report)
				return
			}

			// Log
			app.log.WithFields(log.Fields{
				"path":       proj.Path(),
				"repository": proj.Recipe().Repository().Path(),
				"recipe":     proj.Recipe().Name(),
			}).Info("project loaded")

			// Sync project
			err = app.SyncProject(proj)

			if err != nil {
				report := internalReport.NewErrorReport(err)
				if notify {
					_ = beeep.Alert("Manala", report.String(), "")
				}
				app.log.Report(report)
				return
			}

			if notify {
				_ = beeep.Notify("Manala", "Project synced", "")
			}
		},

		// On all
		func(watcher *internalWatcher.Watcher) {
			if all && proj != nil {
				_ = proj.Recipe().Watch(watcher)
			}
		},
	)
	if err != nil {
		return nil
	}

	//goland:noinspection GoUnhandledErrorResult
	defer watcher.Close()

	watcher.Watch()

	app.log.DecreasePadding()
	app.log.DecreasePadding()

	return nil
}
