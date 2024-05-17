package api

import (
	"manala/app/project"
	"manala/app/project/manifest"
	"manala/app/project/syncer"
	"manala/app/recipe"
	"manala/app/repository"
	"manala/internal/filepath/filter"
)

/***********/
/* Project */
/***********/

func (api *Api) NewProjectLoader(repositoryLoader *repository.Loader, recipeLoader *recipe.Loader, opts ...ProjectLoaderOption) *project.Loader {
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

func (api *Api) WithProjectLoaderFrom(from bool) ProjectLoaderOption {
	return func(options *projectLoaderOptions) {
		options.from = from
	}
}

func (api *Api) NewProjectFinder() *manifest.Finder {
	return manifest.NewFinder()
}

func (api *Api) NewProjectSyncer() *syncer.Syncer {
	return syncer.New(api.log)
}

func (api *Api) NewProjectCreator() *manifest.Creator {
	return manifest.NewCreator()
}
