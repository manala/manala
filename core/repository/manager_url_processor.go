package repository

import (
	"golang.org/x/exp/maps"
	"manala/core"
	internalLog "manala/internal/log"
	"sort"
)

func NewUrlProcessorManager(log *internalLog.Logger, cascadingManager core.RepositoryManager) *UrlProcessorManager {
	return &UrlProcessorManager{
		log:              log,
		cascadingManager: cascadingManager,
		urls:             map[int]string{},
	}
}

type UrlProcessorManager struct {
	log              *internalLog.Logger
	cascadingManager core.RepositoryManager
	urls             map[int]string
}

func (manager *UrlProcessorManager) AddUrl(url string, priority int) {
	manager.urls[priority] = url
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
	// Clone manager urls to add current processed one with priority 0 without touching them
	urls := maps.Clone(manager.urls)
	urls[0] = url

	// Reverse order priorities
	priorities := maps.Keys(urls)
	sort.Sort(sort.Reverse(sort.IntSlice(priorities)))

	var processedUrl string

	for _, priority := range priorities {
		_url := urls[priority]

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
