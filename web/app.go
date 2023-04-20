// Package web needs an app go file system to serve html files and assets from.
//
// Two of them are available, serving different purposes:
//   - "build_embed": this is the one to choose during release, html files and assets coming from "app/build" directory
//     are embedded during the build.
//     Obviously, an "app/build" directory *MUST* have been built before.
//   - "build": this is the one to choose during web app development, html files and assets coming from "app/build" directory
//     are locally served during the runtime.
//     An "app/build" directory should have been built before, and any changes will be applied during the
//     runtime.
//
// Use build tags to switch from "build_embed" or "build"
//   - "-tags=web_app_build_embed": "build_embed"
//   - "-tags=web_app_build": "build"
//
// If none of them is selected, web app is disabled, only web api remains available as it doesn't depend
// on any external assets.
package web

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"io/fs"
	"net/http"
)

var appFS struct {
	Name string
	FS   http.FileSystem
}

func NewApp(fS http.FileSystem) *App {
	return &App{
		fS: fS,
	}
}

type App struct {
	fS http.FileSystem
}

func (app *App) Handler() http.Handler {
	router := chi.NewRouter()

	server := http.FileServer(app.fS)

	router.HandleFunc("/*", func(response http.ResponseWriter, request *http.Request) {
		_, err := app.fS.Open(request.URL.Path)
		if errors.Is(err, fs.ErrNotExist) {
			http.StripPrefix(request.RequestURI, server).ServeHTTP(response, request)
		}
		server.ServeHTTP(response, request)
	})

	return router
}
