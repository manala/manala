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
	// Log
	manager.log.
		WithField("manager", "cache").
		WithField("url", url).
		Debug("load repository")
	manager.log.IncreasePadding()
	defer manager.log.DecreasePadding()

	// Check if repository already in cache
	if repo, ok := manager.cache[url]; ok {
		// Log
		manager.log.
			Debug("cache hit")

		return repo, nil
	}

	// Log
	manager.log.
		Debug("cache miss")

	// Cascading manager
	repo, err := manager.cascadingManager.LoadRepository(url)
	if err != nil {
		return nil, err
	}

	// Cache repository
	manager.cache[url] = repo

	return repo, nil
}
