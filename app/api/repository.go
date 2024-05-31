package api

import (
	"context"
	"manala/app"
	"manala/app/repository"
	"manala/app/repository/cache"
	"manala/app/repository/getter"
	"manala/app/repository/url"
)

func (api *API) NewRepositoryLoader(ctx context.Context) *repository.Loader {
	// URL processor
	urlProcessor := url.NewProcessor(api.log)
	if api.defaultRepositoryURL != "" {
		urlProcessor.Add(api.defaultRepositoryURL, -10)
	}

	if url, ok := app.RepositoryURL(ctx); ok {
		urlProcessor.Add(url, 10)
	}

	if ref, ok := app.RepositoryRef(ctx); ok {
		urlProcessor.AddQuery("ref", ref, 20)
	}

	return repository.NewLoader(
		repository.WithLoaderHandlers(
			url.NewProcessorLoaderHandler(api.log, urlProcessor),
			cache.NewLoaderHandler(api.log, cache.New()),
			getter.NewGitLoaderHandler(api.log, api.cache),
			getter.NewS3LoaderHandler(api.log, api.cache),
			getter.NewHTTPLoaderHandler(api.log, api.cache),
			getter.NewFileLoaderHandler(api.log),
		),
	)
}
