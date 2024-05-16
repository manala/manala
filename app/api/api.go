package api

import (
	"log/slog"
	"manala/internal/cache"
)

// New creates an api
func New(log *slog.Logger, cache *cache.Cache, opts ...Option) *Api {
	api := &Api{
		log:   log,
		cache: cache,
	}

	// Options
	for _, opt := range opts {
		opt(api)
	}

	return api
}

type Api struct {
	log                  *slog.Logger
	cache                *cache.Cache
	defaultRepositoryUrl string
}

type Option func(api *Api)

func WithDefaultRepositoryUrl(url string) Option {
	return func(api *Api) {
		api.defaultRepositoryUrl = url
	}
}
