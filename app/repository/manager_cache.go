package repository

import (
	"log/slog"
	"manala/app"
)

func NewCacheManager(log *slog.Logger, cascadingManager app.RepositoryManager) *CacheManager {
	return &CacheManager{
		log:              log.With("manager", "cache"),
		cache:            make(map[string]app.Repository),
		cascadingManager: cascadingManager,
	}
}

type CacheManager struct {
	log              *slog.Logger
	cache            map[string]app.Repository
	cascadingManager app.RepositoryManager
}

func (manager *CacheManager) LoadRepository(url string) (app.Repository, error) {
	// Log
	manager.log.Debug("load repository",
		"url", url,
	)

	// Check if repository already in cache
	if repository, ok := manager.cache[url]; ok {
		// Log
		manager.log.Debug("repository cache hit",
			"url", url,
		)

		return repository, nil
	}

	// Log
	manager.log.Debug("cascading load repository…",
		"url", url,
	)

	// Cascading manager
	repository, err := manager.cascadingManager.LoadRepository(url)
	if err != nil {
		return nil, err
	}

	// Cache repository
	manager.cache[url] = repository

	return repository, nil
}
