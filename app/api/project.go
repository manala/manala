package api

import (
	"github.com/manala/manala/app/project"
	"github.com/manala/manala/app/project/manifest"
	"github.com/manala/manala/app/project/sync"
	"github.com/manala/manala/app/recipe"
	"github.com/manala/manala/app/repository"
	"github.com/manala/manala/internal/filepath/filter"
)

/***********/
/* Project */
/***********/

func (api *API) NewProjectLoader(repositoryLoader *repository.Loader, recipeLoader *recipe.Loader, opts ...ProjectLoaderOption) *project.Loader {
	var handlers []project.LoaderHandler

	// Options
	options := &projectLoaderOptions{}
	for _, opt := range opts {
		opt(options)
	}

	if options.from {
		handlers = append(handlers,
			manifest.NewFromLoaderHandler(api.log),
		)
	}

	return project.NewLoader(api.log,
		project.WithLoaderFilter(
			filter.New(
				filter.WithDotfiles(false),
				filter.Without(
					"node_modules", // NodeJS
					"vendor",       // Composer
					"venv",         // Python
				),
			),
		),
		project.WithLoaderHandlers(
			append(handlers,
				manifest.NewLoaderHandler(api.log, repositoryLoader, recipeLoader),
			)...,
		),
	)
}

type projectLoaderOptions struct {
	from bool
}

type ProjectLoaderOption func(options *projectLoaderOptions)

func (api *API) WithProjectLoaderFrom(from bool) ProjectLoaderOption {
	return func(options *projectLoaderOptions) {
		options.from = from
	}
}

func (api *API) NewProjectFinder() *manifest.Finder {
	return manifest.NewFinder()
}

func (api *API) NewProjectSyncer() *sync.Syncer {
	return sync.NewSyncer(api.log)
}

func (api *API) NewProjectCreator() *manifest.Creator {
	return manifest.NewCreator()
}
