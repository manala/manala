package repository

import (
	"manala/core"
	internalLog "manala/internal/log"
)

func NewLogManager(log *internalLog.Logger, cascadingManager core.RepositoryManager) *LogManager {
	return &LogManager{
		log:              log,
		cascadingManager: cascadingManager,
	}
}

type LogManager struct {
	log              *internalLog.Logger
	cascadingManager core.RepositoryManager
}

func (manager *LogManager) LoadRepository(url string) (core.Repository, error) {
	// Log
	manager.log.
		WithField("url", url).
		Debug("load repository")
	manager.log.IncreasePadding()

	repo, err := manager.cascadingManager.LoadRepository(url)

	// Log
	manager.log.DecreasePadding()

	if err != nil {
		return nil, err
	}

	return repo, nil
}
