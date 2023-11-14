package api

import (
	"log/slog"
	"manala/app"
	"manala/app/config"
	"manala/app/project"
	"manala/app/recipe"
	"manala/app/repository"
	"manala/internal/cache"
	"manala/internal/syncer"
	"manala/internal/ui"
	"manala/internal/watcher"
)

// New creates an api
func New(config config.Config, log *slog.Logger, out ui.Output, opts ...Option) *Api {
	// Log
	log.Debug("config",
		config.Args()...,
	)

	// Api
	api := &Api{
		config: config,
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
		api.config.CacheDir(),
		cache.WithUserDir("manala"),
	)

	// Syncer
	api.syncer = syncer.New(api.log)

	// Watcher manager
	api.watcherManager = watcher.NewManager(api.log)

	// Repository manager
	api.repositoryManager = repository.NewUrlProcessorManager(
		api.log,
		repository.NewCacheManager(
			api.log,
			repository.NewGetterManager(
				api.log,
				cache,
			),
		),
	)
	api.repositoryManager.AddUrl(
		api.config.Repository(),
		-10,
	)

	// Recipe manager
	api.recipeManager = recipe.NewNameProcessorManager(
		api.log,
		recipe.NewDirManager(
			api.log,
			recipe.WithExclusionPaths(api.exclusionPaths),
		),
	)

	// Project manager
	api.projectManager = project.NewManager(
		api.log,
		api.repositoryManager,
		api.recipeManager,
	)

	// Options
	for _, opt := range opts {
		opt(api)
	}

	return api
}

type Api struct {
	config            config.Config
	log               *slog.Logger
	out               ui.Output
	syncer            *syncer.Syncer
	watcherManager    *watcher.Manager
	repositoryManager *repository.UrlProcessorManager
	recipeManager     *recipe.NameProcessorManager
	projectManager    app.ProjectManager
	exclusionPaths    []string
}

type Option func(api *Api)

func WithRepositoryUrl(url string) Option {
	return func(api *Api) {
		priority := 10

		// Log
		api.log.Debug("repository option",
			"url", url,
			"priority", priority,
		)

		api.repositoryManager.AddUrl(url, priority)
	}
}

func WithRepositoryRef(ref string) Option {
	return func(api *Api) {
		priority := 20

		// Log
		api.log.Debug("repository option",
			"ref", ref,
			"priority", priority,
		)

		api.repositoryManager.AddUrlQuery("ref", ref, priority)
	}
}

func WithRecipeName(name string) Option {
	return func(api *Api) {
		priority := 10

		// Log
		api.log.Debug("recipe option",
			"name", name,
			"priority", priority,
		)

		api.recipeManager.AddName(name, priority)
	}
}
