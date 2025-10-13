package manifest

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"

	"manala/app"
	"manala/app/recipe"
	"manala/internal/serrors"
)

func NewLoaderHandler(log *slog.Logger) *LoaderHandler {
	return &LoaderHandler{
		log: log.With("handler", "manifest"),
	}
}

type LoaderHandler struct {
	log *slog.Logger
}

func (handler *LoaderHandler) Handle(query *recipe.LoaderQuery, chain recipe.LoaderHandlerChain) (app.Recipe, error) {
	dir := filepath.Join(query.Repository.Dir(), query.Name)
	file := filepath.Join(dir, filename)

	handler.log.Debug("handle recipe manifest", "file", file)

	// Stat file
	if fileInfo, err := os.Stat(file); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Chain
			return chain.Next(query)
		}

		return nil, serrors.New("unable to stat recipe manifest").
			WithArguments("file", file).
			WithErrors(serrors.NewOs(err))
	} else if fileInfo.IsDir() {
		return nil, serrors.New("recipe manifest is a directory").
			WithArguments("dir", file)
	}

	manifest := New()

	// Open file
	reader, err := os.Open(file)
	if err != nil {
		return nil, serrors.New("unable to open recipe manifest").
			WithArguments("file", file).
			WithErrors(serrors.NewOs(err))
	}

	// Read from file
	if _, err = manifest.ReadFrom(reader); err != nil {
		return nil, serrors.New("unable to read recipe manifest").
			WithArguments("file", file).
			WithErrors(err)
	}

	handler.log.Debug("recipe manifest loaded", "file", file)

	return NewRecipe(dir, query.Name, manifest, query.Repository), nil
}
