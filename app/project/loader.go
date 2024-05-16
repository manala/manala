package project

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"manala/app"
	"manala/internal/filepath/filter"
	"manala/internal/serrors"
	"os"
	"path/filepath"
)

func NewLoader(log *slog.Logger, filter *filter.Filter, handlers ...LoaderHandler) *Loader {
	return &Loader{
		log:      log,
		filter:   filter,
		handlers: handlers,
	}
}

type Loader struct {
	log      *slog.Logger
	filter   *filter.Filter
	handlers []LoaderHandler
}

//goland:noinspection GoMixedReceiverTypes
func (loader *Loader) Load(dir string) (app.Project, error) {
	loader.log.Info("loading project…")

	// Prepare query
	query := &LoaderQuery{Dir: dir}

	// Start chain
	return loader.Next(query)
}

//goland:noinspection GoMixedReceiverTypes
func (loader *Loader) LoadRecursive(dir string, fn func(project app.Project) error) error {
	loader.log.Info("loading projects recursive…")

	err := filepath.WalkDir(dir,
		func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					return serrors.New("dir not found").
						WithArguments("dir", path)
				}
				return serrors.New("file system error").
					WithArguments("path", path).
					WithErrors(serrors.NewOs(err))
			}

			// Only directories
			if !entry.IsDir() {
				return nil
			}

			// Exclusions
			if loader.filter.Excluded(entry.Name()) {
				loader.log.Debug("exclude project path", "path", path)
				return filepath.SkipDir
			}

			// Load project
			project, err := loader.Load(path)
			if err != nil {
				var _notFoundProjectError *app.NotFoundProjectError
				if errors.As(err, &_notFoundProjectError) {
					err = nil
				}
				return err
			}

			// Walk function
			return fn(project)
		},
	)

	return err
}

//goland:noinspection GoMixedReceiverTypes
func (loader Loader) Next(query *LoaderQuery) (app.Project, error) {
	if len(loader.handlers) == 0 {
		return loader.Last(query)
	}

	handler := loader.handlers[0]
	loader.handlers = loader.handlers[1:]

	return handler.Handle(query, loader)
}

//goland:noinspection GoMixedReceiverTypes
func (loader Loader) Last(query *LoaderQuery) (app.Project, error) {
	return nil, &app.NotFoundProjectError{Dir: query.Dir}
}

type LoaderQuery struct {
	Dir string
}

type LoaderHandler interface {
	Handle(query *LoaderQuery, chain LoaderHandlerChain) (app.Project, error)
}

type LoaderHandlerMock struct {
	mock.Mock
}

func (mock *LoaderHandlerMock) Handle(query *LoaderQuery, chain LoaderHandlerChain) (app.Project, error) {
	args := mock.Called(query, chain)
	return args.Get(0).(app.Project), args.Error(1)
}

type LoaderHandlerChain interface {
	Next(query *LoaderQuery) (app.Project, error)
	Last(query *LoaderQuery) (app.Project, error)
}

type LoaderHandlerChainMock struct {
	mock.Mock
}

func (mock *LoaderHandlerChainMock) Next(query *LoaderQuery) (app.Project, error) {
	args := mock.Called(query)
	return args.Get(0).(app.Project), args.Error(1)
}

func (mock *LoaderHandlerChainMock) Last(query *LoaderQuery) (app.Project, error) {
	args := mock.Called(query)
	return args.Get(0).(app.Project), args.Error(1)
}
