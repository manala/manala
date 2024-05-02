package cache

import "manala/app"

func New() *Cache {
	return &Cache{
		repositories: make(map[string]app.Repository),
	}
}

type Cache struct {
	repositories map[string]app.Repository
}

func (cache *Cache) Get(url string) (app.Repository, bool) {
	repository, ok := cache.repositories[url]
	return repository, ok
}

func (cache *Cache) Set(url string, repository app.Repository) {
	cache.repositories[url] = repository
}
