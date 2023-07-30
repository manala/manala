package recipe

import (
	"errors"
	"io/fs"
	"log/slog"
	"manala/app/interfaces"
	"manala/core"
	"manala/internal/errors/serrors"
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
			return nil, &core.NotFoundRecipeManifestError{File: file}
		}
		return nil, serrors.WrapOs("unable to stat recipe manifest", err).
			WithArguments("file", file)
	} else {
		if fileInfo.IsDir() {
			return nil, serrors.New("recipe manifest is a directory").
				WithArguments("dir", file)
		}
	}

	man := NewManifest()

	// Open file
	if reader, err := os.Open(file); err != nil {
		return nil, serrors.WrapOs("unable to open recipe manifest", err).
			WithArguments("file", file)
	} else {
		// Read from file
		if err = man.ReadFrom(reader); err != nil {
			return nil, serrors.Wrap("unable to read recipe manifest", err).
				WithArguments("file", file)
		}
	}

	// Log
	manager.log.Debug("recipe manifest loaded",
		"description", man.Description(),
		"template", man.Template(),
	)

	return man, nil
}

func (manager *DirManager) LoadRecipe(repo interfaces.Repository, name string) (interfaces.Recipe, error) {
	// Log
	manager.log.Debug("load recipe",
		"name", name,
	)

	dir := filepath.Join(repo.Dir(), name)

	// Load manifest
	manFile := filepath.Join(dir, manifestFilename)
	man, err := manager.loadManifest(manFile)
	if err != nil {
		return nil, err
	}

	rec := NewRecipe(
		dir,
		name,
		man,
		repo,
	)

	return rec, nil
}

func (manager *DirManager) WalkRecipes(repo interfaces.Repository, walker func(rec interfaces.Recipe) error) error {
	// Log
	manager.log.Debug("walk recipes",
		"dir", repo.Dir(),
	)

	dir, err := os.Open(repo.Dir())
	if err != nil {
		return serrors.WrapOs("file system error", err).
			WithArguments("dir", repo.Dir())
	}

	//goland:noinspection GoUnhandledErrorResult
	defer dir.Close()

	files, err := dir.ReadDir(0) // 0 to read all files and folders
	if err != nil {
		return serrors.WrapOs("file system error", err).
			WithArguments("dir", repo.Dir())
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

		rec, err := manager.LoadRecipe(repo, file.Name())
		if err != nil {
			var _notFoundRecipeManifestError *core.NotFoundRecipeManifestError
			if errors.As(err, &_notFoundRecipeManifestError) {
				continue
			}
			return err
		}

		empty = false

		if err := walker(rec); err != nil {
			return err
		}
	}

	if empty {
		return serrors.New("empty repository").
			WithArguments("dir", repo.Dir())
	}

	return nil
}

func (manager *DirManager) WatchRecipe(rec interfaces.Recipe, watcher *watcher.Watcher) error {
	var dirs []string

	// Walk on recipe dirs
	if err := filepath.WalkDir(
		rec.Dir(),
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
