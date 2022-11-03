package recipe

import (
	"errors"
	"fmt"
	"github.com/caarlos0/log"
	"golang.org/x/exp/slices"
	"io/fs"
	"manala/core"
	internalLog "manala/internal/log"
	internalOs "manala/internal/os"
	internalReport "manala/internal/report"
	internalWatcher "manala/internal/watcher"
	"os"
	"path/filepath"
	"sort"
)

const manifestFilename = ".manala.yaml"

func NewManager(log *internalLog.Logger, opts ...ManagerOption) *Manager {
	manager := &Manager{
		log: log,
	}

	// Options
	for _, opt := range opts {
		opt(manager)
	}

	return manager
}

type Manager struct {
	log            *internalLog.Logger
	exclusionPaths []string
}

func (manager *Manager) loadManifest(file string) (*Manifest, error) {
	// Log
	manager.log.WithFields(log.Fields{
		"file": file,
	}).Debug("load recipe manifest")

	// Stat file
	if fileInfo, err := os.Stat(file); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, internalReport.NewError(
				core.NewNotFoundRecipeManifestError("recipe manifest not found"),
			).WithField("file", file)
		}
		return nil, internalReport.NewError(internalOs.NewError(err)).
			WithMessage("unable to stat recipe manifest").
			WithField("file", file)
	} else {
		if fileInfo.IsDir() {
			return nil, internalReport.NewError(fmt.Errorf("recipe manifest is a directory")).
				WithField("dir", file)
		}
	}

	man := NewManifest()

	// Open file
	if reader, err := os.Open(file); err != nil {
		return nil, internalReport.NewError(internalOs.NewError(err)).
			WithMessage("unable to open recipe manifest").
			WithField("file", file)
	} else {
		// Read from file
		if err = man.ReadFrom(reader); err != nil {
			return nil, internalReport.NewError(err).
				WithField("file", file)
		}
	}

	// Log
	manager.log.WithFields(log.Fields{
		"description": man.Description(),
		"template":    man.Template(),
	}).Debug("manifest")

	return man, nil
}

func (manager *Manager) LoadRecipe(repo core.Repository, name string) (core.Recipe, error) {
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

func (manager *Manager) WalkRecipes(repo core.Repository, walker func(rec core.Recipe) error) error {
	// Log
	manager.log.
		WithField("dir", repo.Dir()).
		Debug("walk repository recipes")
	manager.log.IncreasePadding()

	dir, err := os.Open(repo.Dir())
	if err != nil {
		return internalReport.NewError(internalOs.NewError(err)).
			WithMessage("file system error")
	}

	//goland:noinspection GoUnhandledErrorResult
	defer dir.Close()

	files, err := dir.ReadDir(0) // 0 to read all files and folders
	if err != nil {
		return internalReport.NewError(internalOs.NewError(err)).
			WithMessage("file system error")
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
			manager.log.
				WithField("path", file.Name()).
				Debug("exclude path")
			continue
		}

		// Log
		manager.log.WithFields(log.Fields{
			"name": file.Name(),
		}).Debug("load recipe")
		manager.log.IncreasePadding()

		rec, err := manager.LoadRecipe(repo, file.Name())

		// Log
		manager.log.DecreasePadding()

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
		return internalReport.NewError(fmt.Errorf("empty repository")).
			WithField("dir", repo.Dir())
	}

	// Log
	manager.log.DecreasePadding()

	return nil
}

func (manager *Manager) WatchRecipe(rec core.Recipe, watcher *internalWatcher.Watcher) error {
	var dirs []string

	// Walk on recipe dirs
	if err := filepath.WalkDir(rec.Dir(), func(dir string, file fs.DirEntry, err error) error {
		if file.IsDir() {
			dirs = append(dirs, dir)
		}
		return nil
	}); err != nil {
		return err
	}

	return watcher.ReplaceGroup("recipe", dirs)
}

type ManagerOption func(manager *Manager)

func WithExclusionPaths(paths []string) ManagerOption {
	return func(manager *Manager) {
		manager.exclusionPaths = paths
	}
}
