package api

import (
	"log/slog"
	"github.com/manala/manala/internal/caching"
)

// New creates an api.
func New(log *slog.Logger, cache *caching.Cache, opts ...Option) *API {
	api := &API{
		log:   log,
		cache: cache,
	}

	// Options
	for _, opt := range opts {
		opt(api)
	}

	return api
}

type API struct {
	log                  *slog.Logger
	cache                *caching.Cache
	defaultRepositoryURL string
}

type Option func(api *API)

func WithDefaultRepositoryURL(url string) Option {
	return func(api *API) {
		api.defaultRepositoryURL = url
	}
}
