package getter

import (
	"context"
	"log/slog"
	"manala/app"
	"manala/app/repository"
	"manala/internal/cache"
	"manala/internal/schema"
	"time"

	"github.com/hashicorp/go-getter/v2"
)

func NewGitLoaderHandler(log *slog.Logger, cache *cache.Cache) *GitLoaderHandler {
	return &GitLoaderHandler{
		log:   log.With("handler", "getter.git"),
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

type GitLoaderHandler struct {
	log    *slog.Logger
	cache  *cache.Cache
	client *getter.Client
}

func (handler *GitLoaderHandler) Handle(query *repository.LoaderQuery, chain repository.LoaderHandlerChain) (app.Repository, error) {
	handler.log.Debug("handle repository", "url", query.Url)

	// Cache dir
	cacheDir, err := handler.cache.
		WithHashDir(query.Url).
		Dir()
	if err != nil {
		return nil, err
	}

	// Request
	request := &getter.Request{
		Src:     query.Url,
		Dst:     cacheDir,
		GetMode: getter.ModeDir,
	}

	// Force git repo format (ensure backward compatibility)
	if (&schema.GitRepoFormatChecker{}).IsFormat(request.Src) {
		request.Forced = "git"
	}

	response, err := handler.client.Get(context.Background(), request)
	if err != nil {
		if IsNotDetected(err) {
			// Chain
			return chain.Next(query)
		}

		return nil, NewError(err)
	}

	return NewRepository(query.Url, response.Dst), nil
}
