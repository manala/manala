package models

import (
	"fmt"
	"github.com/apex/log"
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
)

/***********/
/* Manager */
/***********/

// NewWatcherManager creates a model watcher manager
func NewWatcherManager(log log.Interface) *watcherManager {
	return &watcherManager{
		log: log,
	}
}

type WatcherManagerInterface interface {
	NewWatcher() (*watcher, error)
}

type watcherManager struct {
	log log.Interface
}

/***********/
/* Watcher */
/***********/

// NewWatcher creates a watcher
func (manager *watcherManager) NewWatcher() (*watcher, error) {
	// Fsnotify watcher
	fsnotifyWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("error creating watcher: %v", err)
	}

	return &watcher{
		log:     manager.log,
		watcher: fsnotifyWatcher,
	}, nil
}

type WatcherInterface interface {
	SetProject(project ProjectInterface) error
	SetRecipe(recipe RecipeInterface) error
	Watch(callback func(watcher WatcherInterface))
	Close()
}

type watcher struct {
	log        log.Interface
	watcher    *fsnotify.Watcher
	projectDir string
	recipeDir  string
}

func (watcher *watcher) SetProject(project ProjectInterface) error {
	dir := project.getDir()
	if dir == watcher.projectDir {
		return nil
	}
	watcher.projectDir = dir

	return watcher.watcher.Add(dir)
}

func (watcher *watcher) SetRecipe(recipe RecipeInterface) error {
	dir := recipe.getDir()

	// If recipe has changed, first, unwatch old one directories
	if (watcher.recipeDir != "") && (dir != watcher.recipeDir) {
		if err := filepath.Walk(watcher.recipeDir, func(path string, info os.FileInfo, err error) error {
			if info.Mode().IsDir() {
				if err := watcher.watcher.Remove(path); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return err
		}
	}

	// Watch all recipe directories; don't care if they are already watched
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.Mode().IsDir() {
			if err := watcher.watcher.Add(path); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	watcher.recipeDir = dir

	return nil
}

func (watcher *watcher) Watch(callback func(watcher WatcherInterface)) {
	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.watcher.Events:
				if !ok {
					return
				}

				watcher.log.WithField("event", event).Debug("Watch event")

				// Ignore chmod events
				if event.Op != fsnotify.Chmod {
					file := filepath.Clean(event.Name)
					dir := filepath.Dir(file)
					if (dir == watcher.projectDir) && (filepath.Base(file) == ProjectManifestFile) {
						// Project manifest
						watcher.log.WithField("file", file).Info("Project manifest modified")
						callback(watcher)
					} else if dir != watcher.projectDir {
						// Recipe dir
						watcher.log.WithField("path", file).Info("Recipe modified")
						callback(watcher)
					}
				}
			case err, ok := <-watcher.watcher.Errors:
				if !ok {
					return
				}

				watcher.log.WithError(err).Error("Watch error")
			}
		}
	}()
	<-done
}

func (watcher *watcher) Close() {
	_ = watcher.watcher.Close()
}
