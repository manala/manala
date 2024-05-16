package manifest

import (
	"github.com/fsnotify/fsnotify"
	"io/fs"
	"log/slog"
	"manala/app"
	"os"
	"path/filepath"
	"slices"
)

func NewWatcher(log *slog.Logger) *Watcher {
	return &Watcher{
		log: log,
	}
}

type Watcher struct {
	log *slog.Logger
}

func (watcher *Watcher) Watch(project app.Project, all bool, fn func(project app.Project) app.Project, done chan os.Signal) error {
	watcher.log.Info("watching project…")

	// Create watcher
	fsnotifyWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	//goland:noinspection GoUnhandledErrorResult
	defer fsnotifyWatcher.Close()

	//recipeDirs := []string{}
	projectManifestPath := filepath.Join(project.Dir(), filename)

	// Start listening for events
	go func() {
		for {
			select {
			case event, ok := <-fsnotifyWatcher.Events:
				if !ok {
					return
				}

				// Ignore chmod events & empty events
				if event.Has(fsnotify.Chmod) || event.Name == "" {
					watcher.log.Debug("ignore event",
						"path", event.Name,
						"operation", event.Op,
					)
					break
				}

				watcher.log.Info("file event",
					"path", event.Name,
					"operation", event.Op,
				)

				// Callback
				if p := fn(project); project != nil {
					project = p
				}

				// Sync watches paths
				paths := []string{projectManifestPath}
				if all {
					paths = append(paths, watcher.recipeDirs(project.Recipe())...)
				}
				_ = watcher.syncPaths(fsnotifyWatcher, paths)

			case err, ok := <-fsnotifyWatcher.Errors:
				if !ok {
					return
				}
				watcher.log.Warn("watch error", "error", err)
			}
		}
	}()

	// Sync watches paths
	paths := []string{projectManifestPath}
	if all {
		paths = append(paths, watcher.recipeDirs(project.Recipe())...)
	}
	_ = watcher.syncPaths(fsnotifyWatcher, paths)

	// Block until done
	<-done

	return nil
}

func (watcher *Watcher) recipeDirs(recipe app.Recipe) []string {
	var dirs []string

	_ = filepath.WalkDir(
		recipe.Dir(),
		func(path string, entry fs.DirEntry, err error) error {
			if entry.IsDir() {
				dirs = append(dirs, path)
			}
			return nil
		},
	)

	return dirs
}

func (watcher *Watcher) syncPaths(fsnotifyWatcher *fsnotify.Watcher, newPaths []string) error {
	// Get old paths
	oldPaths := fsnotifyWatcher.WatchList()

	// Add only new paths (not presents in old ones)
	for _, path := range newPaths {
		if !slices.Contains(oldPaths, path) {
			if err := fsnotifyWatcher.Add(path); err != nil {
				return err
			}
		}
	}

	// Remove only old paths (not presents in new ones)
	for _, path := range oldPaths {
		if !slices.Contains(newPaths, path) {
			if err := fsnotifyWatcher.Remove(path); err != nil {
				return err
			}
		}
	}

	return nil
}
