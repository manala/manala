package api

import (
	"context"
	"embed"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"html/template"
	"manala/app/interfaces"
	"manala/core"
	"manala/core/application"
	"manala/web/api/views"
	"net/http"
)

//go:embed openapi.yaml
var openapi string

//go:embed templates
var templates embed.FS

type contextKey string

func New(app *application.Application, dir string, stop context.CancelFunc) *Api {
	return &Api{
		app:  app,
		dir:  dir,
		stop: stop,
	}
}

type Api struct {
	app  *application.Application
	dir  string
	stop context.CancelFunc
}

func (api *Api) Handler(port int) http.Handler {
	router := chi.NewRouter()

	// Cors - Allow all
	router.Use(cors.AllowAll().Handler)

	// OpenAPI
	router.Get("/openapi.yaml", func(response http.ResponseWriter, request *http.Request) {
		render.PlainText(response, request, openapi)
	})

	// Swagger UI
	router.Get("/", func(response http.ResponseWriter, request *http.Request) {
		tpl, _ := template.ParseFS(templates, "templates/swagger-ui.gohtml")
		_ = tpl.Execute(response, map[string]interface{}{
			"Port": port,
		})
	})

	// Project
	router.Route("/project", func(router chi.Router) {
		router.Get("/", api.GetProject)
	})

	// Recipes
	router.Route("/recipes", func(router chi.Router) {
		router.Get("/", api.ListRecipes)
		router.Route("/{name}", func(router chi.Router) {
			router.Use(api.RecipeContext)
			router.Get("/", api.GetRecipe)
			router.Get("/options", api.GetRecipeOptions)
		})
	})

	// Server
	router.Route("/server", func(router chi.Router) {
		router.Get("/stop", api.StopServer)
	})

	return router
}

func (api *Api) GetProject(response http.ResponseWriter, request *http.Request) {
	proj, err := api.app.LoadProjectFrom(api.dir)
	if err != nil {
		var _notFoundProjectManifestError *core.NotFoundProjectManifestError
		if errors.As(err, &_notFoundProjectManifestError) {
			http.Error(response, http.StatusText(404), 404)
		} else {
			http.Error(response, http.StatusText(500), 500)
		}
		return
	}

	projView := views.NormalizeProject(proj)

	render.JSON(response, request, projView)
}

func (api *Api) ListRecipes(response http.ResponseWriter, request *http.Request) {
	var recViews []*views.RecipeView

	// Walk into recipes
	if err := api.app.WalkRecipes(func(rec interfaces.Recipe) error {
		recViews = append(recViews, views.NormalizeRecipe(rec))
		return nil
	}); err != nil {
		http.Error(response, http.StatusText(500), 500)
		return
	}

	render.JSON(response, request, recViews)
}

func (api *Api) RecipeContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		name := chi.URLParam(request, "name")
		rec, err := api.app.LoadRecipe(name)
		if err != nil {
			var _notFoundRecipeManifestError *core.NotFoundRecipeManifestError
			if errors.As(err, &_notFoundRecipeManifestError) {
				http.Error(response, http.StatusText(404), 404)
			} else {
				http.Error(response, http.StatusText(500), 500)
			}
			return
		}
		ctx := context.WithValue(request.Context(), contextKey("recipe"), rec)
		next.ServeHTTP(response, request.WithContext(ctx))
	})
}

func (api *Api) GetRecipe(response http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	rec, ok := ctx.Value(contextKey("recipe")).(interfaces.Recipe)
	if !ok {
		http.Error(response, http.StatusText(422), 422)
		return
	}

	recView := views.NormalizeRecipe(rec)

	render.JSON(response, request, recView)
}

func (api *Api) GetRecipeOptions(response http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	rec, ok := ctx.Value(contextKey("recipe")).(interfaces.Recipe)
	if !ok {
		http.Error(response, http.StatusText(422), 422)
		return
	}

	var recOptionsViews []*views.RecipeOptionView

	for _, option := range rec.Options() {
		recOptionsViews = append(recOptionsViews, views.NormalizeRecipeOption(option))
	}

	render.JSON(response, request, recOptionsViews)
}

func (api *Api) StopServer(response http.ResponseWriter, request *http.Request) {
	api.stop()
}
