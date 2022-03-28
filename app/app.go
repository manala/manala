package app

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/discard"
	"github.com/spf13/viper"
	"io"
	"manala/fs"
	"manala/loaders"
	"manala/models"
	"manala/syncer"
	"manala/template"
)

// New creates an app
func New(opts ...func(app *App)) *App {

	// Logger
	logger := &log.Logger{
		Handler: discard.Default,
	}

	// Managers
	fsManager := fs.NewManager()
	modelFsManager := models.NewFsManager(fsManager)
	templateManager := template.NewManager()
	modelTemplateManager := models.NewTemplateManager(templateManager, modelFsManager)
	modelWatcherManager := models.NewWatcherManager(logger)

	// Syncer
	sync := syncer.New(logger, modelFsManager, modelTemplateManager)

	// Loaders
	repositoryLoader := loaders.NewRepositoryLoader(logger)
	recipeLoader := loaders.NewRecipeLoader(logger, modelFsManager)
	projectLoader := loaders.NewProjectLoader(logger, repositoryLoader, recipeLoader)

	// App
	app := &App{
		Config:           viper.New(),
		Log:              logger,
		repositoryLoader: repositoryLoader,
		recipeLoader:     recipeLoader,
		projectLoader:    projectLoader,
		templateManager:  modelTemplateManager,
		watcherManager:   modelWatcherManager,
		sync:             sync,
	}

	// Config
	app.Config.SetDefault("debug", false)

	// Options
	for _, opt := range opts {
		opt(app)
	}

	// Apply config
	app.ApplyConfig()

	return app
}

func WithVersion(version string) func(app *App) {
	return func(app *App) {
		app.Config.Set("version", version)
	}
}

func WithDefaultRepository(repository string) func(app *App) {
	return func(app *App) {
		app.Config.SetDefault("repository", repository)
	}
}

func WithLogWriter(writer io.Writer) func(app *App) {
	return func(app *App) {
		app.Log.Handler = cli.New(writer)
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

func (app *App) ApplyConfig() {
	// Log level
	if app.Config.GetBool("debug") {
		app.Log.Level = log.DebugLevel
	} else {
		app.Log.Level = log.InfoLevel
	}
}
