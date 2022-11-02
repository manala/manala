package repository

import (
	"errors"
	"manala/core"
	internalLog "manala/internal/log"
	internalReport "manala/internal/report"
)

func NewChainManager(log *internalLog.Logger, managers []core.RepositoryManager) *ChainManager {
	return &ChainManager{
		log:      log,
		managers: managers,
	}
}

type ChainManager struct {
	log      *internalLog.Logger
	managers []core.RepositoryManager
}

func (manager *ChainManager) LoadRepository(url string) (core.Repository, error) {
	var repo core.Repository

	// Try managers
	for _, _manager := range manager.managers {
		_repo, err := _manager.LoadRepository(url)
		if err != nil {
			var _unsupportedRepositoryError *core.UnsupportedRepositoryError
			if errors.As(err, &_unsupportedRepositoryError) {
				continue
			}
			return nil, err
		}
		repo = _repo
		break
	}

	if repo == nil {
		return nil, internalReport.NewError(
			core.NewUnsupportedRepositoryError("unsupported repository url"),
		).WithField("url", url)
	}

	return repo, nil
}
