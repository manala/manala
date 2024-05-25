package watch

import (
	"context"
	"github.com/fsnotify/fsnotify"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"manala/app"
	"slices"
)

func NewWatcher(log *slog.Logger, recipe bool) *Watcher {
	return &Watcher{
		log:    log,
		recipe: recipe,
	}
}

type Watcher struct {
	log    *slog.Logger
	recipe bool
}

func (watcher *Watcher) Watch(ctx context.Context, project app.Project, fn func(project app.Project) app.Project) error {
	// Create fs watcher
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	//goland:noinspection GoUnhandledErrorResult
	defer fsWatcher.Close()

	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case err := <-fsWatcher.Errors:
				return err
			case event, ok := <-fsWatcher.Events:
				if !ok {
					return nil
				}

				// Ignore chmod events & empty events
				if event.Has(fsnotify.Chmod) || event.Name == "" {
					watcher.log.Debug("ignore file event",
						"path", event.Name,
						"operation", event.Op,
					)
					break
				}

				watcher.log.Debug("file event",
					"path", event.Name,
					"operation", event.Op,
				)

				// Callback
				if p := fn(project); p != nil {
					project = p
				}

				if err := watcher.syncWatches(fsWatcher, project); err != nil {
					return err
				}
			}
		}
	})

	if err := watcher.syncWatches(fsWatcher, project); err != nil {
		return err
	}

	return group.Wait()
}

func (watcher *Watcher) syncWatches(fsWatcher *fsnotify.Watcher, project app.Project) error {
	// Start with project watches
	watches, err := project.Watches()
	if err != nil {
		return err
	}

	// Eventually add recipe watches
	if watcher.recipe {
		recipeWatches, err := project.Recipe().Watches()
		if err != nil {
			return err
		}

		watches = append(watches, recipeWatches...)
	}

	// Get current watches
	currentWatches := fsWatcher.WatchList()

	// Add only new watches (not presents in current ones)
	for _, path := range watches {
		if !slices.Contains(currentWatches, path) {
			watcher.log.Debug("add watch", "path", path)
			if err := fsWatcher.Add(path); err != nil {
				return err
			}
		}
	}

	// Remove only current watches
	for _, path := range currentWatches {
		if !slices.Contains(watches, path) {
			watcher.log.Debug("remove watch", "path", path)
			if err := fsWatcher.Remove(path); err != nil {
				return err
			}
		}
	}

	return nil
}
