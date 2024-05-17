package watch

import (
	"github.com/fsnotify/fsnotify"
	"log/slog"
	"manala/app"
	"os"
	"slices"
)

func watch(log *slog.Logger, project app.Project, recipe bool, fn func(project app.Project) app.Project, done chan os.Signal) error {
	// Create watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	//goland:noinspection GoUnhandledErrorResult
	defer watcher.Close()

	// Start listening to events
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// Ignore chmod events & empty events
				if event.Has(fsnotify.Chmod) || event.Name == "" {
					log.Debug("ignore file event",
						"path", event.Name,
						"operation", event.Op,
					)
					break
				}

				log.Debug("file event",
					"path", event.Name,
					"operation", event.Op,
				)

				// Callback
				if p := fn(project); project != nil {
					project = p
				}

				// Sync watch paths
				paths, _ := project.Watches()
				if recipe {
					recipePaths, _ := project.Recipe().Watches()
					paths = append(paths, recipePaths...)
				}
				_ = syncWatchPaths(watcher, paths)

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Warn("watch error", "error", err)
			}
		}
	}()

	// Sync watch paths
	paths, _ := project.Watches()
	if recipe {
		recipePaths, _ := project.Recipe().Watches()
		paths = append(paths, recipePaths...)
	}
	_ = syncWatchPaths(watcher, paths)

	// Block until done
	<-done

	return nil
}

func syncWatchPaths(watcher *fsnotify.Watcher, new []string) error {
	// Get old paths
	old := watcher.WatchList()

	// Add only new paths (not presents in old ones)
	for _, path := range new {
		if !slices.Contains(old, path) {
			if err := watcher.Add(path); err != nil {
				return err
			}
		}
	}

	// Remove only old paths (not presents in new ones)
	for _, path := range old {
		if !slices.Contains(new, path) {
			if err := watcher.Remove(path); err != nil {
				return err
			}
		}
	}

	return nil
}
