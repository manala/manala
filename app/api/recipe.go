package api

import (
	"context"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/recipe"
	"github.com/manala/manala/app/recipe/manifest"
	"github.com/manala/manala/app/recipe/name"
	"github.com/manala/manala/internal/filepath/filter"
)

func (api *API) NewRecipeLoader(ctx context.Context) *recipe.Loader {
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
