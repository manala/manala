package api

import (
	"manala/app/repository"
	"manala/app/repository/cache"
	"manala/app/repository/getter"
	"manala/app/repository/url"
)

func (api *Api) NewRepositoryLoader(opts ...RepositoryLoaderOption) *repository.Loader {
	// Options
	options := &repositoryLoaderOptions{}
	for _, opt := range opts {
		opt(options)
	}

	// Url processor
	urlProcessor := url.NewProcessor(api.log)
	if api.defaultRepositoryUrl != "" {
		urlProcessor.Add(api.defaultRepositoryUrl, -10)
	}
	if options.url != "" {
		urlProcessor.Add(options.url, 10)
	}
	if options.ref != "" {
		urlProcessor.AddQuery("ref", options.ref, 20)
	}

	return repository.NewLoader(api.log,
		url.NewProcessorLoaderHandler(api.log, urlProcessor),
		cache.NewLoaderHandler(api.log, cache.New()),
		getter.NewGitLoaderHandler(api.log, api.cache),
		getter.NewS3LoaderHandler(api.log, api.cache),
		getter.NewHttpLoaderHandler(api.log, api.cache),
		getter.NewFileLoaderHandler(api.log),
	)
}

type repositoryLoaderOptions struct {
	url string
	ref string
}

type RepositoryLoaderOption func(options *repositoryLoaderOptions)

func (api *Api) WithRepositoryLoaderUrl(url string) RepositoryLoaderOption {
	return func(options *repositoryLoaderOptions) {
		options.url = url
	}
}

func (api *Api) WithRepositoryLoaderRef(ref string) RepositoryLoaderOption {
	return func(options *repositoryLoaderOptions) {
		options.ref = ref
	}
}
