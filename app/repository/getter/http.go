package getter

import (
	"context"
	"time"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/repository"
	"github.com/manala/manala/internal/caching"
	"github.com/manala/manala/internal/log"

	"github.com/hashicorp/go-getter/v2"
)

type HTTPLoaderHandler struct {
	log    *log.Log
	cache  *caching.Cache
	client *getter.Client
}

func NewHTTPLoaderHandler(log *log.Log, cache *caching.Cache) *HTTPLoaderHandler {
	return &HTTPLoaderHandler{
		log:   log,
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

func (handler *HTTPLoaderHandler) Handle(query *repository.LoaderQuery, chain repository.LoaderHandlerChain) (app.Repository, error) {
	handler.log.Debug("handle repository", "handler", "getter.http", "url", query.URL)

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

		return nil, ErrorFrom(err)
	}

	return NewRepository(query.URL, response.Dst), nil
}
