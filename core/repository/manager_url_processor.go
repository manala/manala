package repository

import (
	"dario.cat/mergo"
	"fmt"
	"log/slog"
	"manala/app/interfaces"
	"manala/core"
	"maps"
	netUrl "net/url"
	"sort"
	"strings"
)

func NewUrlProcessorManager(log *slog.Logger, cascadingManager interfaces.RepositoryManager) *UrlProcessorManager {
	return &UrlProcessorManager{
		log:              log.With("manager", "url_processor"),
		cascadingManager: cascadingManager,
		urls:             map[int]string{},
	}
}

type UrlProcessorManager struct {
	log              *slog.Logger
	cascadingManager interfaces.RepositoryManager
	urls             map[int]string
}

func (manager *UrlProcessorManager) AddUrl(url string, priority int) {
	manager.urls[priority] = url
}

func (manager *UrlProcessorManager) AddUrlQuery(key string, value string, priority int) {
	if value == "" {
		return
	}

	query := netUrl.Values{}
	query.Add(key, value)

	manager.urls[priority] = "?" + query.Encode()
}

func (manager *UrlProcessorManager) LoadRepository(url string) (interfaces.Repository, error) {
	// Log
	manager.log.Debug("load repository",
		"url", url,
	)

	// Process url
	url, err := manager.processUrl(url)
	if err != nil {
		return nil, err
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

	return repo, err
}

func (manager *UrlProcessorManager) LoadPrecedingRepository() (interfaces.Repository, error) {
	return manager.LoadRepository("")
}

func (manager *UrlProcessorManager) processUrl(url string) (string, error) {
	// Clone manager urls to add current processed one with priority 0 without touching them
	urls := maps.Clone(manager.urls)
	urls[0] = url

	// Reverse order priorities
	priorities := make([]int, 0, len(urls))
	for priority := range urls {
		priorities = append(priorities, priority)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(priorities)))

	var query string
	var queryValues netUrl.Values

	for _, priority := range priorities {
		// Split url and query parts
		url, query, _ = strings.Cut(urls[priority], "?")

		// Log
		manager.log.Debug("process repository url",
			"url", url,
			"query", query,
			"priority", priority,
		)

		if query != "" {
			_queryValues, err := netUrl.ParseQuery(query)
			if err != nil {
				return "", err
			}
			_ = mergo.Merge(&queryValues, _queryValues)
		}

		if url != "" {
			break
		}
	}

	if url == "" {
		return "", &core.UnprocessableRepositoryUrlError{}
	}

	if queryValues != nil {
		url = fmt.Sprintf("%s?%s", url, queryValues.Encode())
	}

	return url, nil
}
