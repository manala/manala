package repository

import (
	"log/slog"
	"manala/app/interfaces"
)

func NewCacheManager(log *slog.Logger, cascadingManager interfaces.RepositoryManager) *CacheManager {
	return &CacheManager{
		log:              log.With("manager", "cache"),
		cache:            make(map[string]interfaces.Repository),
		cascadingManager: cascadingManager,
	}
}

type CacheManager struct {
	log              *slog.Logger
	cache            map[string]interfaces.Repository
	cascadingManager interfaces.RepositoryManager
}

func (manager *CacheManager) LoadRepository(url string) (interfaces.Repository, error) {
	// Log
	manager.log.Debug("load repository",
		"url", url,
	)

	// Check if repository already in cache
	if repo, ok := manager.cache[url]; ok {
		// Log
		manager.log.Debug("repository cache hit",
			"url", url,
		)

		return repo, nil
	}

	// Log
	manager.log.Debug("cascade repository loading…",
		"url", url,
	)

	// Cascading manager
	repo, err := manager.cascadingManager.LoadRepository(url)
	if err != nil {
		return nil, err
	}

	// Cache repository
	manager.cache[url] = repo

	return repo, nil
}
