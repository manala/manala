package repository

import (
	"github.com/stretchr/testify/mock"
	"log/slog"
	"manala/app"
)

func NewLoader(log *slog.Logger, handlers ...LoaderHandler) *Loader {
	return &Loader{
		log:      log,
		handlers: handlers,
	}
}

type Loader struct {
	log      *slog.Logger
	handlers []LoaderHandler
}

//goland:noinspection GoMixedReceiverTypes
func (loader *Loader) Load(url string) (app.Repository, error) {
	loader.log.Info("loading repository…")

	// Prepare query
	query := &LoaderQuery{Url: url}

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
	return nil, &app.NotFoundRepositoryError{Url: query.Url}
}

type LoaderQuery struct {
	Url string
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
