package application

import (
	"errors"
	"github.com/gen2brain/beeep"
	"log/slog"
	"manala/app/interfaces"
	"manala/core"
	"manala/core/project"
	"manala/core/recipe"
	"manala/core/repository"
	"manala/internal/cache"
	"manala/internal/errors/serrors"
	"manala/internal/filepath/backwalk"
	"manala/internal/syncer"
	"manala/internal/ui/output"
	"manala/internal/watcher"
	"os"
	"path/filepath"
	"slices"
)

// NewApplication creates an application
func NewApplication(conf interfaces.Config, log *slog.Logger, out output.Output, opts ...Option) *Application {
	// Log
	log.Debug("app config",
		conf.Args()...,
	)

	// App
	app := &Application{
		config: conf,
		log:    log,
		out:    out,
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
	cache := cache.New(
		app.config.CacheDir(),
		cache.WithUserDir("manala"),
	)

	// Syncer
	app.syncer = syncer.New(app.log)

	// Watcher manager
	app.watcherManager = watcher.NewManager(app.log)

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
		recipe.NewDirManager(
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
	log               *slog.Logger
	out               output.Output
	syncer            *syncer.Syncer
	watcherManager    *watcher.Manager
	repositoryManager *repository.UrlProcessorManager
	recipeManager     *recipe.NameProcessorManager
	projectManager    interfaces.ProjectManager
	exclusionPaths    []string
}

func (app *Application) WalkRecipes(walker func(rec interfaces.Recipe) error) error {
	var (
		repo interfaces.Repository
		err  error
	)

	// Log
	app.log.Debug("load preceding repository…")

	// Load preceding repository
	repo, err = app.repositoryManager.LoadPrecedingRepository()
	if err != nil {
		return err
	}

	// Log
	app.log.Debug("walk repository recipes…")

	// Walk repository recipes
	return app.recipeManager.WalkRecipes(repo, walker)
}

func (app *Application) CreateProject(
	dir string,
	recSelector func(recipeWalker func(walker func(rec interfaces.Recipe) error) error) (interfaces.Recipe, error),
	optionsSelector func(rec interfaces.Recipe, options []interfaces.RecipeOption) error,
) (interfaces.Project, error) {
	var (
		repo interfaces.Repository
		rec  interfaces.Recipe
		vars map[string]interface{}
		err  error
	)

	// Log
	app.log.Debug("check already existing project…",
		"dir", dir,
	)

	// Check already existing project
	if app.projectManager.IsProject(dir) {
		return nil, &core.AlreadyExistingProjectError{Dir: dir}
	}

	// Log
	app.log.Debug("load preceding repository…")

	// Load preceding repository
	repo, err = app.repositoryManager.LoadPrecedingRepository()
	if err != nil {
		return nil, err
	}

	// Log
	app.log.Debug("try to load preceding recipe…")

	// Try loading preceding recipe
	rec, err = app.recipeManager.LoadPrecedingRecipe(repo)

	if err != nil {
		var _unprocessableRecipeNameError *core.UnprocessableRecipeNameError
		if !errors.As(err, &_unprocessableRecipeNameError) {
			return nil, err
		}

		// Log
		app.log.Debug("unable to load preceding recipe")

		// Log
		app.log.Debug("select recipe…")

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

	// Log
	app.log.Debug("select recipe options…")

	// Select options to get init vars
	vars, err = rec.InitVars(func(options []interfaces.RecipeOption) error {
		if err := optionsSelector(rec, options); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Log
	app.log.Debug("create project…",
		"dir", dir,
	)

	// Create project
	return app.projectManager.CreateProject(dir, rec, vars)
}

func (app *Application) LoadProjectFrom(dir string) (interfaces.Project, error) {
	var (
		proj interfaces.Project
		err  error
	)

	// Log
	app.log.Debug("backwalk projects from…",
		"dir", dir,
	)

	// Backwalk projects from dir
	err = backwalk.Backwalk(
		dir,
		func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return serrors.WrapOs("file system error", err).
					WithArguments("path", path)
			}

			// Log
			app.log.Debug("try to load project…",
				"dir", path,
			)

			// Load project
			proj, err = app.projectManager.LoadProject(path)
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
		return nil, serrors.New("project not found").
			WithArguments("dir", dir)
	}

	return proj, nil
}

func (app *Application) WalkProjects(dir string, walker func(proj interfaces.Project) error) error {
	// Log
	app.log.Info("walk projects from…",
		"dir", dir,
	)

	err := filepath.WalkDir(
		dir,
		func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return serrors.WrapOs("file system error", err).
					WithArguments("path", path)
			}

			// Only directories
			if !entry.IsDir() {
				return nil
			}

			// Exclusions
			if slices.Contains(app.exclusionPaths, filepath.Base(path)) {
				// Log
				app.log.Debug("exclude path",
					"path", path,
				)

				return filepath.SkipDir
			}

			// Log
			app.log.Debug("try to load project…",
				"dir", path,
			)

			// Load project
			proj, err := app.projectManager.LoadProject(path)
			if err != nil {
				var _notFoundProjectManifestError *core.NotFoundProjectManifestError
				if errors.As(err, &_notFoundProjectManifestError) {
					err = nil
				}
				return err
			}

			// Walk function
			return walker(proj)
		},
	)

	return err
}

func (app *Application) SyncProject(proj interfaces.Project) error {
	// Log
	app.log.Info("sync project…",
		"src", proj.Recipe().Dir(),
		"dst", proj.Dir(),
	)

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
	app.log.Info("watch project…",
		"src", proj.Recipe().Dir(),
		"dst", proj.Dir(),
	)

	dir := proj.Dir()

	watcher, err := app.watcherManager.NewWatcher(

		// On start
		func(watcher *watcher.Watcher) {
			// Watch project
			_ = app.projectManager.WatchProject(proj, watcher)
		},

		// On change
		func(watcher *watcher.Watcher) {
			// Load project
			var err error
			proj, err := app.projectManager.LoadProject(dir)
			if err != nil {
				if notify {
					_ = beeep.Alert("Manala", err.Error(), "")
				}
				app.out.Error(err)
				return
			}

			// Sync project
			err = app.SyncProject(proj)

			if err != nil {
				if notify {
					_ = beeep.Alert("Manala", err.Error(), "")
				}
				app.out.Error(err)
				return
			}

			if notify {
				_ = beeep.Notify("Manala", "Project synced", "")
			}
		},

		// On all
		func(watcher *watcher.Watcher) {
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
