package caching

import (
	"github.com/manala/manala/app"
	"github.com/manala/manala/app/repository"
	"github.com/manala/manala/internal/log"
)

type LoaderHandler struct {
	log   *log.Log
	cache *Cache
}

func NewLoaderHandler(log *log.Log, cache *Cache) *LoaderHandler {
	return &LoaderHandler{
		log:   log,
		cache: cache,
	}
}

func (handler *LoaderHandler) Handle(query *repository.LoaderQuery, chain repository.LoaderHandlerChain) (app.Repository, error) {
	handler.log.Debug("handle repository cache", "handler", "cache", "url", query.URL)

	// Check if repository already in cache
	if repository, ok := handler.cache.Get(query.URL); ok {
		handler.log.Debug("hit repository cache", "handler", "cache", "url", query.URL)

		return repository, nil
	}

	// Chain
	repository, err := chain.Next(query)

	// Cache repository
	if repository != nil && err == nil {
		handler.cache.Set(query.URL, repository)
	}

	return repository, err
}
