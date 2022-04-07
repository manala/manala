package watcher

import (
	"github.com/fsnotify/fsnotify"
	internalLog "manala/internal/log"
	"path/filepath"
)

type Watcher struct {
	log *internalLog.Logger
	*fsnotify.Watcher
	onStart     func(watcher *Watcher)
	onChange    func(watcher *Watcher)
	onAll       func(watcher *Watcher)
	temporaries []string
}

func (watcher *Watcher) Watch() {
	// On start
	watcher.onStart(watcher)
	watcher.onAll(watcher)

	// Start watching
	watcher.doWatch()
}

func (watcher *Watcher) doWatch() {
	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// Log
				watcher.log.WithField("event", event).Debug("watch event")

				// Ignore chmod events & empty events
				if event.Op == fsnotify.Chmod || event.Name == "" {
					break
				}

				watcher.log.WithField("path", filepath.Clean(event.Name)).Info("file modified")

				watcher.onChange(watcher)
				watcher.onAll(watcher)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				// Log
				watcher.log.WithError(err).Warn("watch error")
			}
		}
	}()
	<-done
}

func (watcher *Watcher) AddTemporary(name string) error {
	watcher.temporaries = append(watcher.temporaries, name)
	return watcher.Add(name)
}

func (watcher *Watcher) RemoveTemporaries() error {
	for _, name := range watcher.temporaries {
		err := watcher.Remove(name)
		if err != nil {
			return err
		}
	}
	watcher.temporaries = []string{}

	return nil
}
