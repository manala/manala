package app

import (
	"github.com/apex/log"
	"manala/fs"
	"manala/internal/config"
	"manala/loaders"
	"manala/models"
	"manala/syncer"
	"manala/template"
)

// New creates an app
func New(conf *config.Config, logger log.Interface) *App {
	// Debug
	logger.WithFields(conf).Debug("create app")

	// App
	app := &App{
		config: conf,
		log:    logger,
	}

	// Managers
	fsManager := fs.NewManager()
	modelFsManager := models.NewFsManager(fsManager)
	templateManager := template.NewManager()
	app.templateManager = models.NewTemplateManager(templateManager, modelFsManager)
	app.watcherManager = models.NewWatcherManager(app.log)

	// Syncer
	app.sync = syncer.New(app.log, modelFsManager, app.templateManager)

	// Loaders
	app.repositoryLoader = loaders.NewRepositoryLoader(app.log, app.config.GetString("cache-dir"))
	app.recipeLoader = loaders.NewRecipeLoader(app.log, modelFsManager)
	app.projectLoader = loaders.NewProjectLoader(app.log, app.repositoryLoader, app.recipeLoader)

	return app
}

type App struct {
	config           *config.Config
	log              log.Interface
	repositoryLoader loaders.RepositoryLoaderInterface
	recipeLoader     loaders.RecipeLoaderInterface
	projectLoader    loaders.ProjectLoaderInterface
	templateManager  models.TemplateManagerInterface
	watcherManager   models.WatcherManagerInterface
	sync             *syncer.Syncer
}
