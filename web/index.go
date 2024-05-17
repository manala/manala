package web

import (
	"net/http"
)

func (server *Server) index(response http.ResponseWriter, request *http.Request) {
	// Context
	ctx := request.Context()

	// Get repository loader
	repositoryLoader := server.api.NewRepositoryLoader(ctx)

	// Load repository
	server.log.Info("loading repository…")

	repository, err := repositoryLoader.Load("")
	if err != nil {
		server.error(ctx, response, err)

		return
	}

	// Get recipe loader
	recipeLoader := server.api.NewRecipeLoader(ctx)

	// Load recipes
	server.log.Info("loading recipes…")

	recipes, err := recipeLoader.LoadAll(repository)
	if err != nil {
		server.error(ctx, response, err)

		return
	}

	server.mustTemplate(ctx, response, "index.gohtml", map[string]any{
		"Recipes": recipes,
	})
}
