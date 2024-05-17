package web

import (
	"context"
	"errors"
	"log/slog"
	"manala/app"
	"manala/app/api"
	"manala/internal/ui"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

func NewServer(log *slog.Logger, api *api.API, out ui.Output, dir string) *Server {
	return &Server{
		log: log,
		api: api,
		out: out,
		dir: dir,
	}
}

type Server struct {
	log *slog.Logger
	api *api.API
	out ui.Output
	dir string
}

func (server *Server) Serve(ctx context.Context, address string) error {
	// Cancel context
	ctx, cancel := context.WithCancel(ctx)

	// Handler
	httpHandler := http.NewServeMux()

	// Home
	httpHandler.Handle("/", server.handler(server.index))

	// Api
	httpHandler.Handle("/api", server.handler(server.apiIndex))
	httpHandler.Handle("/api/openapi.yaml", server.handler(server.apiDocument))
	httpHandler.Handle("/api/project", server.handler(server.apiProject(server.dir)))
	httpHandler.Handle("/api/recipes", server.handler(server.apiRecipes))
	httpHandler.Handle("/api/recipes/{name}", server.handler(server.apiRecipe))
	httpHandler.HandleFunc("/api/server/stop", server.apiServerStop(cancel))

	// Server
	httpServer := &http.Server{
		Addr:    address,
		Handler: httpHandler,
	}

	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		// Serve
		server.log.Debug("start web server")

		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		return nil
	})
	group.Go(func() error {
		<-ctx.Done()

		// Shutdown
		server.log.Debug("stop web server")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		return httpServer.Shutdown(ctx)
	})

	return group.Wait()
}

func (server *Server) handler(handler http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		// Get query
		query := request.URL.Query()

		// Repository url
		if query.Has("repository") {
			request = request.WithContext(
				app.WithRepositoryURL(request.Context(), query.Get("repository")),
			)
		}

		// Repository ref
		if query.Has("ref") {
			request = request.WithContext(
				app.WithRepositoryRef(request.Context(), query.Get("ref")),
			)
		}

		// Recipe
		if query.Has("recipe") {
			request = request.WithContext(
				app.WithRecipeName(request.Context(), query.Get("recipe")),
			)
		}

		// Log
		server.log.Info("serving…",
			"method", request.Method,
			"url", request.URL.Path,
		)

		handler.ServeHTTP(response, request)
	})
}
