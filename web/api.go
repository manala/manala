package web

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"manala/app"
	"net/http"
	"path/filepath"
	"strings"
)

const swaggerUIVersion = "5.17.12"

//go:embed api.yaml
var apiDocument []byte

func (server *Server) apiIndex(response http.ResponseWriter, request *http.Request) {
	if _, err := io.Copy(response, strings.NewReader(
		fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
	<head>
      	<meta charset="utf-8" />
      	<meta name="viewport" content="width=device-width, initial-scale=1" />
      	<meta name="description" content="SwaggerUI" />
      	<title>SwaggerUI</title>
      	<link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@%[2]s/swagger-ui.css" />
    </head>
    <body>
    	<div id="swagger-ui"></div>
    	<script src="https://unpkg.com/swagger-ui-dist@%[2]s/swagger-ui-bundle.js" crossorigin></script>
    	<script>
      		window.onload = () => {
        		window.ui = SwaggerUIBundle({
          			url: 'http://%[1]s/api/openapi.yaml',
          			dom_id: '#swagger-ui',
        		});
      		};
    	</script>
    </body>
</html>`,
			request.Host,
			swaggerUIVersion,
		))); err != nil {
		server.error(request.Context(), response, err)

		return
	}
}

func (server *Server) apiDocument(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/yaml")

	if _, err := io.Copy(response, bytes.NewReader(apiDocument)); err != nil {
		server.error(request.Context(), response, err)

		return
	}
}

func (server *Server) apiProject(dir string) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		// Get query
		query := request.URL.Query()

		// Dir
		if query.Has("dir") {
			dir = filepath.Clean(query.Get("dir"))
		}

		// Context
		ctx := request.Context()

		// Get repository loader
		repositoryLoader := server.api.NewRepositoryLoader(ctx)

		// Get recipe loader
		recipeLoader := server.api.NewRecipeLoader(ctx)

		// Get project loader
		projectLoader := server.api.NewProjectLoader(repositoryLoader, recipeLoader,
			server.api.WithProjectLoaderFrom(true),
		)

		// Load project
		server.log.Info("loading project…")

		project, err := projectLoader.Load(dir)
		if err != nil {
			server.error(ctx, response, err)

			return
		}

		response.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(response).Encode(app.NewProjectView(project)); err != nil {
			server.error(ctx, response, err)

			return
		}
	}
}

func (server *Server) apiRecipes(response http.ResponseWriter, request *http.Request) {
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

	response.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(response).Encode(app.NewRecipesView(recipes)); err != nil {
		server.error(ctx, response, err)

		return
	}
}

func (server *Server) apiRecipe(response http.ResponseWriter, request *http.Request) {
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

	// Load recipe
	server.log.Info("loading recipe…")

	recipe, err := recipeLoader.Load(repository, request.PathValue("name"))
	if err != nil {
		server.error(ctx, response, err)

		return
	}

	response.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(response).Encode(app.NewRecipeView(recipe)); err != nil {
		server.error(ctx, response, err)

		return
	}
}

func (server *Server) apiServerStop(cancel context.CancelFunc) http.HandlerFunc {
	return func(_ http.ResponseWriter, _ *http.Request) {
		cancel()
	}
}
