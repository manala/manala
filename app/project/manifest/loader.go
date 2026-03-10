package manifest

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/project"
	"github.com/manala/manala/app/recipe"
	"github.com/manala/manala/app/repository"
	"github.com/manala/manala/internal/filepath/backwalk"
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/validator"
)

type LoaderHandler struct {
	log              *slog.Logger
	repositoryLoader *repository.Loader
	recipeLoader     *recipe.Loader
}

func NewLoaderHandler(log *slog.Logger, repositoryLoader *repository.Loader, recipeLoader *recipe.Loader) *LoaderHandler {
	return &LoaderHandler{
		log:              log.With("handler", "manifest"),
		repositoryLoader: repositoryLoader,
		recipeLoader:     recipeLoader,
	}
}

func (handler *LoaderHandler) Handle(query *project.LoaderQuery, chain project.LoaderHandlerChain) (app.Project, error) {
	dir := query.Dir
	file := filepath.Join(dir, filename)

	handler.log.Debug("handle project manifest", "file", file)

	// Stat file
	if fileInfo, err := os.Stat(file); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Chain
			return chain.Next(query)
		}

		return nil, serrors.New("unable to stat project manifest").
			WithArguments("file", file).
			WithErrors(serrors.NewOs(err))
	} else if fileInfo.IsDir() {
		return nil, serrors.New("project manifest is a directory").
			WithArguments("dir", file)
	}

	manifest := New()

	// Open file
	reader, err := os.Open(file)
	if err != nil {
		return nil, serrors.New("unable to open project manifest").
			WithArguments("file", file).
			WithErrors(serrors.NewOs(err))
	}

	// Read from file
	if _, err = manifest.ReadFrom(reader); err != nil {
		return nil, serrors.New("unable to read project manifest").
			WithArguments("file", file).
			WithErrors(err)
	}

	handler.log.Debug("project manifest loaded", "file", file,
		"repository", manifest.Repository(),
		"recipe", manifest.Recipe(),
	)

	// Load repository
	repository, err := handler.repositoryLoader.Load(manifest.Repository())
	if err != nil {
		return nil, err
	}

	// Load recipe
	recipe, err := handler.recipeLoader.Load(repository, manifest.Recipe())
	if err != nil {
		return nil, err
	}

	project := NewProject(dir, manifest, recipe)

	// Validate project vars against recipe
	if violations, err := validator.New(
		validator.WithValidators(project.Recipe().ProjectValidator()),
		validator.WithFormatters(manifest.ValidatorFormatter()),
	).Validate(project.Vars()); err != nil {
		return nil, serrors.New("unable to validate project manifest").
			WithArguments("file", file).
			WithErrors(err)
	} else if len(violations) != 0 {
		return nil, serrors.New("invalid project manifest vars").
			WithArguments("file", file).
			WithErrors(violations.StructuredErrors()...)
	}

	return project, nil
}

type FromLoaderHandler struct {
	log *slog.Logger
}

func NewFromLoaderHandler(log *slog.Logger) *FromLoaderHandler {
	return &FromLoaderHandler{
		log: log.With("handler", "manifest.from"),
	}
}

func (handler *FromLoaderHandler) Handle(query *project.LoaderQuery, chain project.LoaderHandlerChain) (app.Project, error) {
	dir := query.Dir

	handler.log.Debug("handle project manifest from", "dir", dir)

	var project app.Project

	// Backwalk from dir
	if err := backwalk.WalkDir(dir,
		func(path string, _ os.DirEntry, err error) error {
			if err != nil {
				return serrors.New("file system error").
					WithArguments("path", path).
					WithErrors(serrors.NewOs(err))
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
