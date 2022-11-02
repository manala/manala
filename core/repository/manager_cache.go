package repository

import (
	"manala/core"
	internalLog "manala/internal/log"
)

func NewCacheManager(log *internalLog.Logger, cascadingManager core.RepositoryManager) *CacheManager {
	return &CacheManager{
		log:              log,
		cache:            make(map[string]core.Repository),
		cascadingManager: cascadingManager,
	}
}

type CacheManager struct {
	log              *internalLog.Logger
	cache            map[string]core.Repository
	cascadingManager core.RepositoryManager
}

func (manager *CacheManager) LoadRepository(url string) (core.Repository, error) {
	// Check if repository already in cache
	if repo, ok := manager.cache[url]; ok {
		manager.log.
			WithField("manager", "cache").
			Debug("load from cache")
		return repo, nil
	}

	manager.log.
		WithField("manager", "cache").
		Debug("load from cascading manager")

	repo, err := manager.cascadingManager.LoadRepository(url)
	if err != nil {
		return nil, err
	}

	// Cache repository
	manager.cache[url] = repo

	return repo, nil
}
