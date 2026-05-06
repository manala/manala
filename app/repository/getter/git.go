package getter

import (
	"context"
	"time"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/repository"
	"github.com/manala/manala/internal/cache"
	"github.com/manala/manala/internal/log"
	"github.com/manala/manala/internal/validation"

	"github.com/hashicorp/go-getter/v2"
)

type GitLoaderHandler struct {
	log    *log.Log
	cache  *cache.Cache
	client *getter.Client
}

func NewGitLoaderHandler(log *log.Log, cache *cache.Cache) *GitLoaderHandler {
	return &GitLoaderHandler{
		log:   log,
		cache: cache.WithDir("repositories"),
		client: &getter.Client{
			// Prevent copying or writing files through symlinks
			DisableSymlinks: true,
			Getters: []getter.Getter{
				&getter.GitGetter{
					Detectors: []getter.Detector{
						&getter.GitHubDetector{},
						&getter.GitDetector{},
						&getter.BitBucketDetector{},
						&getter.GitLabDetector{},
					},
					Timeout: 30 * time.Second,
				},
			},
			Decompressors: getter.Decompressors,
		},
	}
}

func (handler *GitLoaderHandler) Handle(query *repository.LoaderQuery, chain repository.LoaderHandlerChain) (app.Repository, error) {
	handler.log.Debug("handle repository", "handler", "getter.git", "url", query.URL)

	// Cache dir
	cacheDir, err := handler.cache.
		WithHashDir(query.URL).
		Dir()
	if err != nil {
		return nil, err
	}

	// Request
	request := &getter.Request{
		Src:     query.URL,
		Dst:     cacheDir,
		GetMode: getter.ModeDir,
	}

	// Legacy: ensure backward compatibility by forcing git repo format ()
	if err := validation.GitRepoFormat.Validate(request.Src); err == nil {
		request.Forced = "git"
	}

	response, err := handler.client.Get(context.Background(), request)
	if err != nil {
		if IsNotDetected(err) {
			// Chain
			return chain.Next(query)
		}

		return nil, ErrorFrom(err)
	}

	return NewRepository(query.URL, response.Dst), nil
}
