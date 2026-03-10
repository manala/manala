package repository

import (
	"github.com/manala/manala/app"

	"github.com/stretchr/testify/mock"
)

func NewLoader(opts ...LoaderOption) *Loader {
	loader := &Loader{}

	// Options
	for _, opt := range opts {
		opt(loader)
	}

	return loader
}

type Loader struct {
	handlers []LoaderHandler
}

//goland:noinspection GoMixedReceiverTypes
func (loader *Loader) Load(url string) (app.Repository, error) {
	// Prepare query
	query := &LoaderQuery{URL: url}

	// Start chain
	return loader.Next(query)
}

//goland:noinspection GoMixedReceiverTypes
func (loader Loader) Next(query *LoaderQuery) (app.Repository, error) {
	if len(loader.handlers) == 0 {
		return loader.Last(query)
	}

	handler := loader.handlers[0]
	loader.handlers = loader.handlers[1:]

	return handler.Handle(query, loader)
}

//goland:noinspection GoMixedReceiverTypes
func (loader Loader) Last(query *LoaderQuery) (app.Repository, error) {
	return nil, &app.NotFoundRepositoryError{URL: query.URL}
}

type LoaderOption func(loader *Loader)

func WithLoaderHandlers(handlers ...LoaderHandler) LoaderOption {
	return func(loader *Loader) {
		loader.handlers = append(loader.handlers, handlers...)
	}
}

type LoaderQuery struct {
	URL string
}

type LoaderHandler interface {
	Handle(query *LoaderQuery, chain LoaderHandlerChain) (app.Repository, error)
}

type LoaderHandlerMock struct {
	mock.Mock
}

func (mock *LoaderHandlerMock) Handle(query *LoaderQuery, chain LoaderHandlerChain) (app.Repository, error) {
	args := mock.Called(query, chain)

	return args.Get(0).(app.Repository), args.Error(1)
}

type LoaderHandlerChain interface {
	Next(query *LoaderQuery) (app.Repository, error)
	Last(query *LoaderQuery) (app.Repository, error)
}

type LoaderHandlerChainMock struct {
	mock.Mock
}

func (mock *LoaderHandlerChainMock) Next(query *LoaderQuery) (app.Repository, error) {
	args := mock.Called(query)

	return args.Get(0).(app.Repository), args.Error(1)
}

func (mock *LoaderHandlerChainMock) Last(query *LoaderQuery) (app.Repository, error) {
	args := mock.Called(query)

	return args.Get(0).(app.Repository), args.Error(1)
}
