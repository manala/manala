package recipe

import (
	"errors"
	"io/fs"
	"log/slog"
	"manala/app"
	"manala/internal/serrors"
	"manala/internal/watcher"
	"os"
	"path/filepath"
	"slices"
	"sort"
)

const manifestFilename = ".manala.yaml"

func NewDirManager(log *slog.Logger, opts ...ManagerOption) *DirManager {
	manager := &DirManager{
		log: log.With("manager", "dir"),
	}

	// Options
	for _, opt := range opts {
		opt(manager)
	}

	return manager
}

type DirManager struct {
	log            *slog.Logger
	exclusionPaths []string
}

func (manager *DirManager) loadManifest(file string) (*Manifest, error) {
	// Log
	manager.log.Debug("try to load recipe manifest",
		"file", file,
	)

	// Stat file
	if fileInfo, err := os.Stat(file); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, &app.NotFoundRecipeManifestError{File: file}
		}
		return nil, serrors.New("unable to stat recipe manifest").
			WithArguments("file", file).
			WithErrors(serrors.NewOs(err))
	} else {
		if fileInfo.IsDir() {
			return nil, serrors.New("recipe manifest is a directory").
				WithArguments("dir", file)
		}
	}

	manifest := NewManifest()

	// Open file
	if reader, err := os.Open(file); err != nil {
		return nil, serrors.New("unable to open recipe manifest").
			WithArguments("file", file).
			WithErrors(serrors.NewOs(err))
	} else {
		// Read from file
		if _, err = manifest.ReadFrom(reader); err != nil {
			return nil, serrors.New("unable to read recipe manifest").
				WithArguments("file", file).
				WithErrors(err)
		}
	}

	// Log
	manager.log.Debug("recipe manifest loaded",
		"description", manifest.Description(),
		"template", manifest.Template(),
	)

	return manifest, nil
}

func (manager *DirManager) LoadRecipe(repository app.Repository, name string) (app.Recipe, error) {
	// Log
	manager.log.Debug("load recipe",
		"name", name,
	)

	dir := filepath.Join(repository.Dir(), name)

	// Load manifest
	manifestFile := filepath.Join(dir, manifestFilename)
	manifest, err := manager.loadManifest(manifestFile)
	if err != nil {
		return nil, err
	}

	return New(
		dir,
		name,
		manifest,
		repository,
	), nil
}

func (manager *DirManager) WalkRecipes(repository app.Repository, walker func(recipe app.Recipe) error) error {
	// Log
	manager.log.Debug("walk recipes",
		"dir", repository.Dir(),
	)

	dir, err := os.Open(repository.Dir())
	if err != nil {
		return serrors.New("file system error").
			WithArguments("dir", repository.Dir()).
			WithErrors(serrors.NewOs(err))
	}

	//goland:noinspection GoUnhandledErrorResult
	defer dir.Close()

	files, err := dir.ReadDir(0) // 0 to read all files and folders
	if err != nil {
		return serrors.New("file system error").
			WithArguments("dir", repository.Dir()).
			WithErrors(serrors.NewOs(err))
	}

	// Sort alphabetically
	sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })

	empty := true

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		// Exclusions
		if slices.Contains(manager.exclusionPaths, filepath.Base(file.Name())) {
			// Log
			manager.log.Debug("exclude recipe path",
				"path", file.Name(),
			)

			continue
		}

		recipe, err := manager.LoadRecipe(repository, file.Name())
		if err != nil {
			var _notFoundRecipeManifestError *app.NotFoundRecipeManifestError
			if errors.As(err, &_notFoundRecipeManifestError) {
				continue
			}
			return err
		}

		empty = false

		if err := walker(recipe); err != nil {
			return err
		}
	}

	if empty {
		return &app.EmptyRepositoryError{Repository: repository}
	}

	return nil
}

func (manager *DirManager) WatchRecipe(recipe app.Recipe, watcher *watcher.Watcher) error {
	var dirs []string

	// Walk on recipe dirs
	if err := filepath.WalkDir(
		recipe.Dir(),
		func(path string, entry fs.DirEntry, err error) error {
			if entry.IsDir() {
				dirs = append(dirs, path)
			}
			return nil
		},
	); err != nil {
		return err
	}

	return watcher.ReplaceGroup("recipe", dirs)
}

type ManagerOption func(manager *DirManager)

func WithExclusionPaths(paths []string) ManagerOption {
	return func(manager *DirManager) {
		manager.exclusionPaths = paths
	}
}
