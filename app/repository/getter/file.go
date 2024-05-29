package getter

import (
	"context"
	"errors"
	"log/slog"
	"manala/app"
	"manala/app/repository"
	"manala/internal/serrors"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-getter/v2"
)

func NewFileLoaderHandler(log *slog.Logger) *FileLoaderHandler {
	return &FileLoaderHandler{
		log: log.With("handler", "getter.file"),
		client: &getter.Client{
			// Prevent copying or writing files through symlinks
			DisableSymlinks: true,
			Getters: []getter.Getter{
				&getter.FileGetter{},
			},
			Decompressors: getter.Decompressors,
		},
	}
}

type FileLoaderHandler struct {
	log    *slog.Logger
	client *getter.Client
}

func (handler *FileLoaderHandler) Handle(query *repository.LoaderQuery, chain repository.LoaderHandlerChain) (app.Repository, error) {
	handler.log.Debug("handle repository", "url", query.Url)

	// Request
	request := &getter.Request{
		Src:     query.Url,
		GetMode: getter.ModeDir,
		// In local file mode, the returned operation will simply contain the source file path
		Inplace: true,
	}

	// Stat
	stat, err := os.Stat(request.Src)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Chain
			return chain.Next(query)
		}

		return nil, serrors.New("file system error").
			WithArguments("path", request.Src).
			WithErrors(serrors.NewOs(err))
	} else if !stat.IsDir() {
		// Chain
		return chain.Next(query)
	}

	// Set pwd if relative
	if !filepath.IsAbs(request.Src) {
		var err error
		request.Pwd, err = os.Getwd()
		if err != nil {
			return nil, serrors.New("unable to get current directory")
		}
	}

	response, err := handler.client.Get(context.Background(), request)
	if err != nil {
		if IsNotDetected(err) {
			// Chain
			return chain.Next(query)
		}

		return nil, NewError(err)
	}

	// Switch back to relative dst
	if request.Pwd != "" {
		var err error
		response.Dst, err = filepath.Rel(request.Pwd, response.Dst)
		if err != nil {
			return nil, serrors.New("unable to get relative path")
		}
	}

	return NewRepository(query.Url, response.Dst), nil
}
