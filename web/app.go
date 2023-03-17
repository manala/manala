// Package web needs an app go file system to serve html files and assets from.
//
// Three of them are available, serving different purposes:
//   - "public_embed": this is the default one, html files and assets coming from "app/public" directory are embedded
//     during the build.
//     This ensure builds to never fail, as "app/public" directory is always present.
//     On the other side, it only serve a minimal, non-dynamic, and non-functional web ui.
//   - "build_embed": this is the one to choose during release, html files and assets coming from "app/build" directory
//     are embedded during the build.
//     Obviously, an "app/build" directory *MUST* have been built before.
//   - "build": this is the one to choose during web app development, html files and assets coming from "app/build" directory
//     are locally served during the runtime.
//     An "app/build" directory should have been built before, and any changes will be applied during the
//     runtime.
//
// Use build tags to switch from "public_embed" to "build_embed" or "build"
//   - "-tags=web_app_build_embed": "build_embed"
//   - "-tags=web_app_build": "build"
package web

import (
	"github.com/go-chi/chi/v5"
	internalLog "manala/internal/log"
	"net/http"
)

func NewApp(log *internalLog.Logger) *App {
	return &App{
		log: log,
	}
}

type App struct {
	log *internalLog.Logger
}

func (app *App) Handler() http.Handler {
	appFS, appFSName := AppFS()

	app.log.
		WithField("fs", appFSName).
		Debug("handle web app")

	router := chi.NewRouter()

	router.Handle("/*", http.FileServer(appFS))

	return router
}
