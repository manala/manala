package api

import (
	"context"
	"manala/app"
	"manala/app/recipe"
	"manala/app/recipe/manifest"
	"manala/app/recipe/name"
	"manala/internal/filepath/filter"
)

func (api *Api) NewRecipeLoader(ctx context.Context) *recipe.Loader {
	// Name processor
	nameProcessor := name.NewProcessor(api.log)
	if name, ok := app.RecipeName(ctx); ok {
		nameProcessor.Add(name, 10)
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
