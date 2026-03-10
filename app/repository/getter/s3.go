package getter

import (
	"context"
	"log/slog"
	"time"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/repository"
	"github.com/manala/manala/internal/caching"

	"github.com/hashicorp/go-getter/s3/v2"
	"github.com/hashicorp/go-getter/v2"
)

type S3LoaderHandler struct {
	log    *slog.Logger
	cache  *caching.Cache
	client *getter.Client
}

func NewS3LoaderHandler(log *slog.Logger, cache *caching.Cache) *S3LoaderHandler {
	return &S3LoaderHandler{
		log:   log.With("handler", "getter.s3"),
		cache: cache.WithDir("repositories"),
		client: &getter.Client{
			// Prevent copying or writing files through symlinks
			DisableSymlinks: true,
			Getters: []getter.Getter{
				&s3.Getter{
					Timeout: 30 * time.Second,
				},
			},
			Decompressors: getter.Decompressors,
		},
	}
}

func (handler *S3LoaderHandler) Handle(query *repository.LoaderQuery, chain repository.LoaderHandlerChain) (app.Repository, error) {
	handler.log.Debug("handle repository", "url", query.URL)

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

	response, err := handler.client.Get(context.Background(), request)
	if err != nil {
		if IsNotDetected(err) {
			// Chain
			return chain.Next(query)
		}

		return nil, NewError(err)
	}

	return NewRepository(query.URL, response.Dst), nil
}
