package getter

import (
	"context"
	"log/slog"
	"manala/app"
	"manala/app/repository"
	"manala/internal/cache"
	"time"

	"github.com/hashicorp/go-getter/v2"
)

func NewHttpLoaderHandler(log *slog.Logger, cache *cache.Cache) *HttpLoaderHandler {
	return &HttpLoaderHandler{
		log:   log.With("handler", "getter.http"),
		cache: cache.WithDir("repositories"),
		client: &getter.Client{
			// Prevent copying or writing files through symlinks
			DisableSymlinks: true,
			Getters: []getter.Getter{
				&getter.HttpGetter{
					// Will look up and use auth information found in the user's netrc file if available
					Netrc: true,
					// Disables the client's usage of the "X-Terraform-Get" header value
					XTerraformGetDisabled: true,
					// Enforce a timeout when the server supports HEAD requests
					HeadFirstTimeout: 10 * time.Second,
					// Enforce a timeout when making a request to an HTTP server and reading its response body
					ReadTimeout: 30 * time.Second,
				},
			},
			Decompressors: getter.Decompressors,
		},
	}
}

type HttpLoaderHandler struct {
	log    *slog.Logger
	cache  *cache.Cache
	client *getter.Client
}

func (handler *HttpLoaderHandler) Handle(query *repository.LoaderQuery, chain repository.LoaderHandlerChain) (app.Repository, error) {
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
