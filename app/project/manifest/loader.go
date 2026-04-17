package manifest

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/project"
	"github.com/manala/manala/app/recipe"
	"github.com/manala/manala/app/repository"
	"github.com/manala/manala/internal/filepath/backwalk"
	"github.com/manala/manala/internal/log"
	"github.com/manala/manala/internal/parsing"
	"github.com/manala/manala/internal/serrors"
)

type LoaderHandler struct {
	log              *log.Log
	repositoryLoader *repository.Loader
	recipeLoader     *recipe.Loader
}

func NewLoaderHandler(log *log.Log, repositoryLoader *repository.Loader, recipeLoader *recipe.Loader) *LoaderHandler {
	return &LoaderHandler{
		log:              log,
		repositoryLoader: repositoryLoader,
		recipeLoader:     recipeLoader,
	}
}

func (handler *LoaderHandler) Handle(query *project.LoaderQuery, chain project.LoaderHandlerChain) (app.Project, error) {
	dir := query.Dir
	file := filepath.Join(dir, filename)

	handler.log.Debug("handle project manifest", "handler", "manifest", "file", file)

	// Stat file
	if fileInfo, err := os.Stat(file); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Chain
			return chain.Next(query)
		}

		return nil, serrors.New("unable to stat project manifest").
			With("file", file).
			WithErrors(serrors.FromOs(err))
	} else if fileInfo.IsDir() {
		return nil, serrors.New("project manifest is a directory").
			With("dir", file)
	}

	// Open file
	reader, err := os.Open(file)
	if err != nil {
		return nil, serrors.New("unable to open project manifest").
			With("file", file).
			WithErrors(serrors.FromOs(err))
	}
	defer reader.Close()

	// Read file
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, serrors.New("unable to read project manifest").
			With("file", file).
			WithErrors(err)
	}

	// Parse file content
	manifest := New()
	if err := manifest.Unmarshal(content); err != nil {
		e := serrors.New("unable to parse project manifest")
		if err, ok := errors.AsType[*parsing.Error](err); ok {
			return nil, e.WithDumper(parsing.ErrorDumper{
				Err:   err,
				File:  file,
				Src:   string(content),
				Lexer: "yaml",
			})
		}
		return nil, e.With("file", file).
			WithErrors(err)
	}

	handler.log.Debug("project manifest loaded", "handler", "manifest",
		"file", file,
		"repository", manifest.Repository,
		"recipe", manifest.Recipe,
	)

	// Load repository
	repository, err := handler.repositoryLoader.Load(manifest.Repository)
	if err != nil {
		return nil, err
	}

	// Load recipe
	recipe, err := handler.recipeLoader.Load(repository, manifest.Recipe)
	if err != nil {
		return nil, err
	}

	project := NewProject(dir, manifest, recipe)

	// Validate project vars against recipe
	if violations, err := project.Recipe().ProjectValidator().Validate(project.Vars()); err != nil {
		return nil, serrors.New("unable to validate project manifest").
			With("file", file).
			WithErrors(err)
	} else if len(violations) != 0 {
		return nil, serrors.New("invalid project manifest vars").
			With("file", file).
			WithErrors(violations.StructuredErrors()...)
	}

	return project, nil
}

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

	// Backwalk from dir
	if err := backwalk.WalkDir(dir,
		func(path string, _ os.DirEntry, err error) error {
			if err != nil {
				return serrors.New("file system error").
					With("path", path).
					WithErrors(serrors.FromOs(err))
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
