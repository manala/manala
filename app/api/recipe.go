package api

import (
	"manala/app/recipe"
	"manala/app/recipe/manifest"
	"manala/app/recipe/name"
	"manala/internal/filepath/filter"
)

func (api *Api) NewRecipeLoader(opts ...RecipeLoaderOption) *recipe.Loader {
	// Options
	options := &recipeLoaderOptions{}
	for _, opt := range opts {
		opt(options)
	}

	// Name processor
	nameProcessor := name.NewProcessor(api.log)
	if options.name != "" {
		nameProcessor.Add(options.name, 10)
	}

	return recipe.NewLoader(api.log,
		recipe.WithLoaderFilter(
			filter.New(
				filter.WithDotfiles(false),
			),
		),
		recipe.WithLoaderHandlers(
			name.NewProcessorLoaderHandler(api.log, nameProcessor),
			manifest.NewLoaderHandler(api.log),
		),
	)
}

type recipeLoaderOptions struct {
	name string
}

type RecipeLoaderOption func(options *recipeLoaderOptions)

func (api *Api) WithRecipeLoaderName(name string) RecipeLoaderOption {
	return func(options *recipeLoaderOptions) {
		options.name = name
	}
}
