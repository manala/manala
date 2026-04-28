package manifest

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/recipe"
	"github.com/manala/manala/internal/log"
	"github.com/manala/manala/internal/parsing"
	"github.com/manala/manala/internal/serrors"
)

type LoaderHandler struct {
	log *log.Log
}

func NewLoaderHandler(log *log.Log) *LoaderHandler {
	return &LoaderHandler{
		log: log,
	}
}

func (handler *LoaderHandler) Handle(query *recipe.LoaderQuery, chain recipe.LoaderHandlerChain) (app.Recipe, error) {
	dir := filepath.Join(query.Repository.Dir(), query.Name)
	file := filepath.Join(dir, filename)

	handler.log.Debug("handle recipe manifest", "handler", "manifest", "file", file)

	// Stat file
	if fileInfo, err := os.Stat(file); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Chain
			return chain.Next(query)
		}

		return nil, serrors.New("unable to stat recipe manifest").
			With("file", file).
			WithErrors(serrors.FromOs(err))
	} else if fileInfo.IsDir() {
		return nil, serrors.New("recipe manifest is a directory").
			With("dir", file)
	}

	// Open file
	reader, err := os.Open(file)
	if err != nil {
		return nil, serrors.New("unable to open recipe manifest").
			With("file", file).
			WithErrors(serrors.FromOs(err))
	}
	defer reader.Close()

	// Read file
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, serrors.New("unable to read recipe manifest").
			With("file", file).
			WithErrors(err)
	}

	// Parse file content
	manifest := New()
	if err := manifest.Unmarshal(content); err != nil {
		e := serrors.New("unable to parse recipe manifest")
		if err, ok := errors.AsType[*parsing.Error](err); ok {
			return nil, e.WithDumper(parsing.ErrorDumper{
				Err:   err.Flatten(),
				File:  file,
				Src:   string(content),
				Lexer: "yaml",
			})
		}
		return nil, e.With("file", file).
			WithErrors(err)
	}

	handler.log.Debug("recipe manifest loaded", "handler", "manifest", "file", file)

	return NewRecipe(dir, query.Name, manifest, query.Repository), nil
}
