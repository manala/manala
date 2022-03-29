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
		Config: viper.New(),
		Log:    logger,
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
	app.watcherManager = models.NewWatcherManager(app.Log)

	// Syncer
	app.sync = syncer.New(app.Log, modelFsManager, app.templateManager)

	// Loaders
	app.repositoryLoader = loaders.NewRepositoryLoader(app.Log)
	app.recipeLoader = loaders.NewRecipeLoader(app.Log, modelFsManager)
	app.projectLoader = loaders.NewProjectLoader(app.Log, app.repositoryLoader, app.recipeLoader)

	return app
}

func WithConfig(config *viper.Viper) func(app *App) {
	return func(app *App) {
		app.Config = config
	}
}

func WithLogger(logger *log.Logger) func(app *App) {
	return func(app *App) {
		app.Log = logger
	}
}

type App struct {
	Config           *viper.Viper
	Log              *log.Logger
	repositoryLoader loaders.RepositoryLoaderInterface
	recipeLoader     loaders.RecipeLoaderInterface
	projectLoader    loaders.ProjectLoaderInterface
	templateManager  models.TemplateManagerInterface
	watcherManager   models.WatcherManagerInterface
	sync             *syncer.Syncer
}
