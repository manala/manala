package recipe

import (
	"errors"
	"os"
	"sort"

	"github.com/manala/manala/app"
	"github.com/manala/manala/internal/filepath/filter"
	"github.com/manala/manala/internal/log"
	"github.com/manala/manala/internal/serrors"

	"github.com/stretchr/testify/mock"
)

type Loader struct {
	log      *log.Log
	filter   *filter.Filter
	handlers []LoaderHandler
}

func NewLoader(log *log.Log, opts ...LoaderOption) *Loader {
	loader := &Loader{
		log: log,
	}

	// Options
	for _, opt := range opts {
		opt(loader)
	}

	return loader
}

func (loader *Loader) Load(repository app.Repository, name string) (app.Recipe, error) {
	// Prepare query
	query := &LoaderQuery{Repository: repository, Name: name}

	// Start chain
	return loader.Next(query)
}

func (loader *Loader) LoadAll(repository app.Repository) ([]app.Recipe, error) {
	dir, err := os.Open(repository.Dir())
	if err != nil {
		return nil, serrors.New("file system error").
			With("dir", repository.Dir()).
			WithErrors(serrors.FromOs(err))
	}

	defer dir.Close()

	files, err := dir.ReadDir(0) // 0 to read all files and folders
	if err != nil {
		return nil, serrors.New("file system error").
			With("dir", repository.Dir()).
			WithErrors(serrors.FromOs(err))
	}

	// Sort alphabetically
	sort.Slice(files, func(a, b int) bool { return files[a].Name() < files[b].Name() })

	recipes := make([]app.Recipe, 0)

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		if loader.filter != nil {
			// Exclusions
			if loader.filter.Excluded(file.Name()) {
				loader.log.Debug("exclude recipe path", "path", file.Name())

				continue
			}
		}

		recipe, err := loader.Load(repository, file.Name())
		if err != nil {
			if _, ok := errors.AsType[*app.NotFoundRecipeError](err); ok {
				continue
			}

			return nil, err
		}

		recipes = append(recipes, recipe)
	}

	if len(recipes) == 0 {
		return nil, &app.EmptyRepositoryError{Repository: repository}
	}

	return recipes, nil
}

func (loader Loader) Next(query *LoaderQuery) (app.Recipe, error) {
	if len(loader.handlers) == 0 {
		return loader.Last(query)
	}

	handler := loader.handlers[0]
	loader.handlers = loader.handlers[1:]

	return handler.Handle(query, loader)
}

func (loader Loader) Last(query *LoaderQuery) (app.Recipe, error) {
	return nil, &app.NotFoundRecipeError{Repository: query.Repository, Name: query.Name}
}

type LoaderOption func(loader *Loader)

func WithLoaderFilter(filter *filter.Filter) LoaderOption {
	return func(loader *Loader) {
		loader.filter = filter
	}
}

func WithLoaderHandlers(handlers ...LoaderHandler) LoaderOption {
	return func(loader *Loader) {
		loader.handlers = append(loader.handlers, handlers...)
	}
}

type LoaderQuery struct {
	Repository app.Repository
	Name       string
}

type LoaderHandler interface {
	Handle(query *LoaderQuery, chain LoaderHandlerChain) (app.Recipe, error)
}

type LoaderHandlerMock struct {
	mock.Mock
}

func (mock *LoaderHandlerMock) Handle(query *LoaderQuery, chain LoaderHandlerChain) (app.Recipe, error) {
	args := mock.Called(query, chain)

	return args.Get(0).(app.Recipe), args.Error(1)
}

type LoaderHandlerChain interface {
	Next(query *LoaderQuery) (app.Recipe, error)
	Last(query *LoaderQuery) (app.Recipe, error)
}

type LoaderHandlerChainMock struct {
	mock.Mock
}

func (mock *LoaderHandlerChainMock) Next(query *LoaderQuery) (app.Recipe, error) {
	args := mock.Called(query)

	return args.Get(0).(app.Recipe), args.Error(1)
}

func (mock *LoaderHandlerChainMock) Last(query *LoaderQuery) (app.Recipe, error) {
	args := mock.Called(query)

	return args.Get(0).(app.Recipe), args.Error(1)
}
