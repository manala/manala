package caching

import "github.com/manala/manala/app"

type Cache struct {
	store map[string]app.Repository
}

func NewCache() *Cache {
	return &Cache{
		store: make(map[string]app.Repository),
	}
}

func (cache *Cache) Get(url string) (app.Repository, bool) {
	repository, ok := cache.store[url]

	return repository, ok
}

func (cache *Cache) Set(url string, repository app.Repository) {
	cache.store[url] = repository
}
