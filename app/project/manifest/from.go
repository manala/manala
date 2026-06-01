package manifest

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/project"
	"github.com/manala/manala/internal/errors/serror"
	"github.com/manala/manala/internal/errors/std"
	"github.com/manala/manala/internal/filepath/backwalk"
	"github.com/manala/manala/internal/log"
)

type FromLoaderHandler struct {
	log *log.Log
}

func NewFromLoaderHandler(log *log.Log) *FromLoaderHandler {
	return &FromLoaderHandler{
		log: log,
	}
}

func (handler *FromLoaderHandler) Handle(query *project.LoaderQuery, chain project.LoaderHandlerChain) (app.Project, error) {
	dir := query.Dir

	handler.log.Debug("handle project manifest from", "handler", "manifest.from", "dir", dir)

	var project app.Project

	// Stat dir
	if dirInfo, err := os.Stat(dir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, serror.New("project from dir does not exist").
				With("dir", dir)
		}

		return nil, serror.New("unable to stat project from dir").
			With("dir", dir).
			WithErr(std.From(err))
	} else if !dirInfo.IsDir() {
		return nil, serror.New("project from dir is not a dir").
			With("dir", dir)
	}

	// Backwalk from dir
	if err := backwalk.WalkDir(dir,
		func(path string, _ os.DirEntry, err error) error {
			if err != nil {
				return serror.New("file system error").
					With("path", path).
					WithErr(std.From(err))
			}

			// Update query
			query.Dir = path

			// Load project
			project, err = chain.Next(query)
			if err != nil {
				if _, ok := errors.AsType[*app.NotFoundProjectError](err); ok {
					return nil
				}

				return err
			}

			// Stop backwalk
			return filepath.SkipAll
		},
	); err != nil {
		return nil, err
	}

	if project == nil {
		query.Dir = dir

		return chain.Last(query)
	}

	return project, nil
}
