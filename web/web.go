package web

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"golang.org/x/sync/errgroup"
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

func New(log *internalLog.Logger, app *application.Application, dir string) *Web {
	web := &Web{
		log: log,
		app: app,
		dir: dir,
	}

	return web
}

type Web struct {
	log *internalLog.Logger
	app *application.Application
	dir string
}

func (web *Web) ListenAndServe(port int) error {
	// Catch system signals
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Router
	router := chi.NewRouter()
	router.Use(NewLogger(web.log))

	// Server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      router,
	}

	// Api
	webApi := api.New(web.log, web.app, web.dir, stop)
	router.Mount("/api", webApi.Handler(port))

	// App
	webApp := NewApp(web.log)
	router.Mount("/", webApp.Handler())

	// Log
	web.log.
		WithField("port", port).
		WithField("dir", web.dir).
		Info("start web server")

	// Server listener
	serverListener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		return err
	}

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
