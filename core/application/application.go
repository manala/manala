package application

import (
	"errors"
	"fmt"
	"github.com/gen2brain/beeep"
	"golang.org/x/exp/slices"
	"manala/app/interfaces"
	"manala/core"
	"manala/core/project"
	"manala/core/recipe"
	"manala/core/repository"
	internalCache "manala/internal/cache"
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
func NewApplication(conf interfaces.Config, log *internalLog.Logger, opts ...Option) *Application {
	// Log
	log.
		WithFields(conf).
		Debug("config")

	// App
	app := &Application{
		config: conf,
		log:    log,
		exclusionPaths: []string{
			// Git
			".git", ".github",
			// NodeJS
			"node_modules",
			// Composer
			"vendor",
			// IntelliJ
			".idea",
			// Manala
			".manala",
		},
	}

	// Cache
	cache := internalCache.New(
		app.config.CacheDir(),
		internalCache.WithUserDir("manala"),
	)

	// Syncer
	app.syncer = internalSyncer.New(app.log)

	// Watcher manager
	app.watcherManager = internalWatcher.NewManager(app.log)

	// Repository manager
	app.repositoryManager = repository.NewUrlProcessorManager(
		app.log,
		repository.NewCacheManager(
			app.log,
			repository.NewGetterManager(
				app.log,
				cache,
			),
		),
	)
	app.repositoryManager.AddUrl(
		app.config.Repository(),
		-10,
	)

	// Recipe manager
	app.recipeManager = recipe.NewNameProcessorManager(
		app.log,
		recipe.NewManager(
			app.log,
			recipe.WithExclusionPaths(app.exclusionPaths),
		),
	)

	// Project manager
	app.projectManager = project.NewManager(
		app.log,
		app.repositoryManager,
		app.recipeManager,
	)

	// Options
	for _, opt := range opts {
		opt(app)
	}

	return app
}

type Application struct {
	config            interfaces.Config
	log               *internalLog.Logger
	syncer            *internalSyncer.Syncer
	watcherManager    *internalWatcher.Manager
	repositoryManager *repository.UrlProcessorManager
	recipeManager     *recipe.NameProcessorManager
	projectManager    interfaces.ProjectManager
	exclusionPaths    []string
}

func (app *Application) WalkRecipes(walker func(rec interfaces.Recipe) error) error {
	// Load repository
	repo, err := app.repositoryManager.LoadPrecedingRepository()
	if err != nil {
		return err
	}

	return app.recipeManager.WalkRecipes(repo, walker)
}

func (app *Application) CreateProject(
	dir string,
	recSelector func(recipeWalker func(walker func(rec interfaces.Recipe) error) error) (interfaces.Recipe, error),
	optionsSelector func(rec interfaces.Recipe, options []interfaces.RecipeOption) error,
) (interfaces.Project, error) {
	// Ensure no already existing project
	if app.projectManager.IsProject(dir) {
		return nil, internalReport.NewError(fmt.Errorf("already existing project")).
			WithField("dir", dir)
	}

	// Load repository
	repo, err := app.repositoryManager.LoadPrecedingRepository()
	if err != nil {
		return nil, err
	}

	var rec interfaces.Recipe

	// Try with preceding recipe
	rec, err = app.recipeManager.LoadPrecedingRecipe(repo)

	if err != nil {
		var _unprocessableRecipeNameError *core.UnprocessableRecipeNameError
		if !errors.As(err, &_unprocessableRecipeNameError) {
			return nil, err
		}

		// Or use recipe selector
		rec, err = recSelector(func(walker func(rec interfaces.Recipe) error) error {
			if err := app.recipeManager.WalkRecipes(repo, func(_rec interfaces.Recipe) error {
				if err := walker(_rec); err != nil {
					return err
				}
				return nil
			}); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	// Select options to get init vars
	vars, err := rec.InitVars(func(options []interfaces.RecipeOption) error {
		if err := optionsSelector(rec, options); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return app.projectManager.CreateProject(dir, rec, vars)
}

func (app *Application) LoadProjectFrom(dir string) (interfaces.Project, error) {
	// Log
	app.log.
		WithField("dir", dir).
		Debug("load project from")
	app.log.IncreasePadding()
	defer app.log.DecreasePadding()

	var proj interfaces.Project

	// Backwalks from dir
	err := internalFilepath.Backwalk(
		dir,
		func(_dir string, file os.DirEntry, err error) error {
			if err != nil {
				return internalReport.NewError(internalOs.NewError(err)).
					WithMessage("file system error")
			}

			// Load project
			proj, err = app.projectManager.LoadProject(_dir)
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

	if err != nil {
		return nil, err
	}

	if proj == nil {
		return nil, internalReport.NewError(
			core.NewNotFoundProjectManifestError("project manifest not found"),
		).WithField("dir", dir)
	}

	return proj, nil
}

func (app *Application) WalkProjects(dir string, walker func(proj interfaces.Project) error) error {
	// Log
	app.log.
		WithField("dir", dir).
		Info("load projects from")
	app.log.IncreasePadding()
	defer app.log.DecreasePadding()

	err := filepath.WalkDir(dir, func(_dir string, file os.DirEntry, err error) error {
		if err != nil {
			return internalReport.NewError(internalOs.NewError(err)).
				WithMessage("file system error")
		}

		// Only directories
		if !file.IsDir() {
			return nil
		}

		// Exclusions
		if slices.Contains(app.exclusionPaths, filepath.Base(_dir)) {
			app.log.
				WithField("path", _dir).
				Debug("exclude path")
			return filepath.SkipDir
		}

		// Load project
		proj, err := app.projectManager.LoadProject(_dir)
		if err != nil {
			var _notFoundProjectManifestError *core.NotFoundProjectManifestError
			if errors.As(err, &_notFoundProjectManifestError) {
				err = nil
			}
			return err
		}

		// Walk function
		return walker(proj)
	})

	return err
}

func (app *Application) SyncProject(proj interfaces.Project) error {
	// Log
	app.log.
		WithField("src", proj.Recipe().Dir()).
		WithField("dst", proj.Dir()).
		Info("sync project")
	app.log.IncreasePadding()
	defer app.log.DecreasePadding()

	// Loop over project recipe sync units
	for _, unit := range proj.Recipe().Sync() {
		if err := app.syncer.Sync(
			proj.Recipe().Dir(),
			unit.Source(),
			proj.Dir(),
			unit.Destination(),
			proj,
		); err != nil {

			return err
		}
	}

	return nil
}

func (app *Application) WatchProject(proj interfaces.Project, all bool, notify bool) error {
	// Log
	app.log.
		WithField("src", proj.Recipe().Dir()).
		WithField("dst", proj.Dir()).
		Info("watch project")
	app.log.IncreasePadding()
	defer app.log.DecreasePadding()

	dir := proj.Dir()

	watcher, err := app.watcherManager.NewWatcher(

		// On start
		func(watcher *internalWatcher.Watcher) {
			// Watch project
			_ = app.projectManager.WatchProject(proj, watcher)
		},

		// On change
		func(watcher *internalWatcher.Watcher) {
			// Load project
			var err error
			proj, err := app.projectManager.LoadProject(dir)
			if err != nil {
				report := internalReport.NewErrorReport(err)
				if notify {
					_ = beeep.Alert("Manala", report.String(), "")
				}
				app.log.Report(report)
				return
			}

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
				_ = app.recipeManager.WatchRecipe(proj.Recipe(), watcher)
			}
		},
	)
	if err != nil {
		return nil
	}

	//goland:noinspection GoUnhandledErrorResult
	defer watcher.Close()

	watcher.Watch()

	return nil
}
