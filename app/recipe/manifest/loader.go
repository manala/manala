package manifest

import (
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/recipe"
	"github.com/manala/manala/internal/parsing"
	"github.com/manala/manala/internal/serrors"
)

type LoaderHandler struct {
	log *slog.Logger
}

func NewLoaderHandler(log *slog.Logger) *LoaderHandler {
	return &LoaderHandler{
		log: log.With("handler", "manifest"),
	}
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
			WithErrors(serrors.FromOs(err))
	} else if fileInfo.IsDir() {
		return nil, serrors.New("recipe manifest is a directory").
			WithArguments("dir", file)
	}

	// Open file
	reader, err := os.Open(file)
	if err != nil {
		return nil, serrors.New("unable to open recipe manifest").
			WithArguments("file", file).
			WithErrors(serrors.FromOs(err))
	}
	defer reader.Close()

	// Read file
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, serrors.New("unable to read recipe manifest").
			WithArguments("file", file).
			WithErrors(err)
	}

	// Parse file content
	// Bypass yaml.Unmarshal as goccy's decoder discards comments
	manifest := New()
	if err := manifest.UnmarshalYAML(content); err != nil {
		e := serrors.New("unable to parse recipe manifest").WithArguments("file", file)
		if err, ok := errors.AsType[*parsing.Error](err); ok {
			return nil, parsing.ErrorTo(e, err, parsing.Options{
				Src:   string(content),
				Lexer: "yaml",
			})
		}
		return nil, e.WithErrors(err)
	}

	handler.log.Debug("recipe manifest loaded", "file", file)

	return NewRecipe(dir, query.Name, manifest, query.Repository), nil
}
