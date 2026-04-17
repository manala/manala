package api

import (
	"github.com/manala/manala/internal/caching"
	"github.com/manala/manala/internal/log"
)

// New creates an api.
func New(log *log.Log, cache *caching.Cache, opts ...Option) *API {
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
	log                  *log.Log
	cache                *caching.Cache
	defaultRepositoryURL string
}

type Option func(api *API)

func WithDefaultRepositoryURL(url string) Option {
	return func(api *API) {
		api.defaultRepositoryURL = url
	}
}
