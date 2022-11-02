package repository

import (
	"manala/core"
	internalLog "manala/internal/log"
)

func NewUrlProcessorManager(log *internalLog.Logger, cascadingManager core.RepositoryManager) *UrlProcessorManager {
	return &UrlProcessorManager{
		log:              log,
		cascadingManager: cascadingManager,
	}
}

type UrlProcessorManager struct {
	log              *internalLog.Logger
	lowermostUrl     string
	uppermostUrl     string
	cascadingManager core.RepositoryManager
}

func (manager *UrlProcessorManager) WithLowermostUrl(url string) {
	manager.lowermostUrl = url
}

func (manager *UrlProcessorManager) WithUppermostUrl(url string) {
	manager.uppermostUrl = url
}

func (manager *UrlProcessorManager) LoadRepository(url string) (core.Repository, error) {
	url, err := manager.processUrl(url)
	if err != nil {
		return nil, err
	}

	repo, err := manager.cascadingManager.LoadRepository(url)
	if err != nil {
		return nil, err
	}

	return repo, err
}

func (manager *UrlProcessorManager) LoadPrecedingRepository() (core.Repository, error) {
	return manager.LoadRepository("")
}

func (manager *UrlProcessorManager) processUrl(url string) (string, error) {
	var processedUrl string

	for _, _url := range []string{manager.uppermostUrl, url, manager.lowermostUrl} {
		if _url == "" {
			continue
		}

		processedUrl = _url
		break
	}

	if processedUrl == "" {
		return "", core.NewUnprocessableRepositoryUrlError(
			"unable to process empty repository url",
		)
	}

	return processedUrl, nil
}
