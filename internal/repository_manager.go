package internal

import (
	"errors"
	internalLog "manala/internal/log"
	"strings"
)

func NewRepositoryManager(logger *internalLog.Logger, defaultRepository string) *RepositoryManager {
	return &RepositoryManager{
		Log:               logger,
		DefaultRepository: defaultRepository,
		cache:             make(map[string]*Repository),
	}
}

type RepositoryManager struct {
	Log               *internalLog.Logger
	DefaultRepository string
	RepositoryLoaders []RepositoryLoaderInterface
	cache             map[string]*Repository
}

func (manager *RepositoryManager) AddRepositoryLoader(loader RepositoryLoaderInterface) *RepositoryManager {
	manager.RepositoryLoaders = append(manager.RepositoryLoaders, loader)

	return manager
}

func (manager *RepositoryManager) LoadRepository(paths []string) (*Repository, error) {
	// Check if repository already in cache
	if repository, ok := manager.cache[strings.Join(paths, "|")]; ok {
		// Log
		manager.Log.Debug("load from cache")

		return repository, nil
	}

	var repository *Repository

	// Try loaders
	for _, loader := range manager.RepositoryLoaders {
		_repository, err := loader.LoadRepository(append(paths, manager.DefaultRepository))
		if err != nil {
			var _err *UnsupportedRepositoryError
			if errors.As(err, &_err) {
				continue
			}
			return nil, err
		}
		repository = _repository
		break
	}
	if repository == nil {
		return nil, UnsupportedRepositoryPathError()
	}

	// Cache repository
	manager.cache[strings.Join(paths, "|")] = repository

	return repository, nil
}
