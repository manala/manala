package cache

import (
	"log/slog"
	"manala/app"
	"manala/app/repository"
)

func NewLoaderHandler(log *slog.Logger, cache *Cache) *LoaderHandler {
	return &LoaderHandler{
		log:   log.With("handler", "cache"),
		cache: cache,
	}
}

type LoaderHandler struct {
	log   *slog.Logger
	cache *Cache
}

func (handler *LoaderHandler) Handle(query *repository.LoaderQuery, chain repository.LoaderHandlerChain) (app.Repository, error) {
	handler.log.Debug("handle repository cache", "url", query.Url)

	// Check if repository already in cache
	if repository, ok := handler.cache.Get(query.Url); ok {
		handler.log.Debug("hit repository cache", "url", query.Url)

		return repository, nil
	}

	// Chain
	repository, err := chain.Next(query)

	// Cache repository
	if repository != nil && err == nil {
		handler.cache.Set(query.Url, repository)
	}

	return repository, err
}
