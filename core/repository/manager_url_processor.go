package repository

import (
	"fmt"
	"github.com/imdario/mergo"
	"golang.org/x/exp/maps"
	"manala/app/interfaces"
	"manala/core"
	internalLog "manala/internal/log"
	netUrl "net/url"
	"sort"
	"strings"
)

func NewUrlProcessorManager(log *internalLog.Logger, cascadingManager interfaces.RepositoryManager) *UrlProcessorManager {
	return &UrlProcessorManager{
		log:              log,
		cascadingManager: cascadingManager,
		urls:             map[int]string{},
	}
}

type UrlProcessorManager struct {
	log              *internalLog.Logger
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
	manager.log.
		WithField("manager", "url_processor").
		WithField("url", url).
		Debug("load repository")
	manager.log.IncreasePadding()
	defer manager.log.DecreasePadding()

	// Process url
	url, err := manager.processUrl(url)
	if err != nil {
		return nil, err
	}

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
	priorities := maps.Keys(urls)
	sort.Sort(sort.Reverse(sort.IntSlice(priorities)))

	var query string
	var queryValues netUrl.Values

	for _, priority := range priorities {
		// Split url and query parts
		url, query, _ = strings.Cut(urls[priority], "?")

		manager.log.
			WithField("url", url).
			WithField("query", query).
			WithField("priority", priority).
			Debug("process url")

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
		return "", core.NewUnprocessableRepositoryUrlError(
			"unable to process empty repository url",
		)
	}

	if queryValues != nil {
		url = fmt.Sprintf("%s?%s", url, queryValues.Encode())
	}

	return url, nil
}
