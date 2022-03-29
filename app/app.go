package app

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/discard"
	"github.com/spf13/viper"
	"manala/fs"
	"manala/loaders"
	"manala/models"
	"manala/syncer"
	"manala/template"
)

// New creates an app
func New(opts ...func(app *App)) *App {
	// Default logger
	logger := &log.Logger{
		Handler: discard.Default,
	}

	// App
	app := &App{
		config: viper.New(),
		log:    logger,
	}

	// Options
	for _, opt := range opts {
		opt(app)
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
	app.repositoryLoader = loaders.NewRepositoryLoader(app.log)
	app.recipeLoader = loaders.NewRecipeLoader(app.log, modelFsManager)
	app.projectLoader = loaders.NewProjectLoader(app.log, app.repositoryLoader, app.recipeLoader)

	return app
}

func WithConfig(config *viper.Viper) func(app *App) {
	return func(app *App) {
		app.config = config
	}
}

func WithLogger(logger *log.Logger) func(app *App) {
	return func(app *App) {
		app.log = logger
	}
}

type App struct {
	config           *viper.Viper
	log              *log.Logger
	repositoryLoader loaders.RepositoryLoaderInterface
	recipeLoader     loaders.RecipeLoaderInterface
	projectLoader    loaders.ProjectLoaderInterface
	templateManager  models.TemplateManagerInterface
	watcherManager   models.WatcherManagerInterface
	sync             *syncer.Syncer
}
