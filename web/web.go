package web

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"golang.org/x/sync/errgroup"
	"manala/app/interfaces"
	"manala/core/application"
	internalLog "manala/internal/log"
	"manala/web/api"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func New(log *internalLog.Logger, app *application.Application, conf interfaces.Config, dir string) *Web {
	web := &Web{
		log:    log,
		app:    app,
		config: conf,
		dir:    dir,
	}

	return web
}

type Web struct {
	log    *internalLog.Logger
	app    *application.Application
	config interfaces.Config
	dir    string
}

func (web *Web) ListenAndServe() error {
	// Catch system signals
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer stop()

	// Router
	router := chi.NewRouter()
	router.Use(NewLogger(web.log))

	// Server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", web.config.WebPort()),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      router,
	}

	// Api
	web.log.
		Debug("handle web api")

	webApi := api.New(web.app, web.dir, stop)
	router.Mount("/api", webApi.Handler(web.config.WebPort()))

	// App
	if appFS.FS != nil {
		web.log.
			WithField("fs", appFS.Name).
			Debug("handle web app")

		webApp := NewApp(appFS.FS)
		router.Mount("/", webApp.Handler())
	}

	// Log
	web.log.
		WithField("port", web.config.WebPort()).
		WithField("dir", web.dir).
		Info("start web server")

	// Server listener
	serverListener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		return err
	}

	// Open browser
	//url := fmt.Sprintf("http://localhost:%d", port)

	//web.log.
	//	WithField("url", url).
	//	Info("open web browser")

	//OpenBrowser(url)

	errGroup, errCtx := errgroup.WithContext(ctx)

	// Serve
	errGroup.Go(func() error {
		if err := server.Serve(serverListener); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				return err
			}
		}
		return nil
	})

	// Shutdown
	errGroup.Go(func() error {
		<-errCtx.Done()

		// Log
		web.log.
			Info("stop web server")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		return server.Shutdown(ctx)
	})

	return errGroup.Wait()
}
