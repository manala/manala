package watcher

import (
	"github.com/fsnotify/fsnotify"
	internalLog "manala/internal/log"
	"path/filepath"
)

type Watcher struct {
	log *internalLog.Logger
	*fsnotify.Watcher
	onStart  func(watcher *Watcher)
	onChange func(watcher *Watcher)
	onAll    func(watcher *Watcher)
	groups   map[string][]string
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
				watcher.log.
					WithField("operation", event.Op).
					WithField("path", event.Name).
					Debug("watch event")

				// Ignore chmod events & empty events
				if event.Has(fsnotify.Chmod) || event.Name == "" {
					watcher.log.IncreasePadding()
					watcher.log.
						Debug("ignore event")
					watcher.log.DecreasePadding()
					break
				}

				watcher.log.
					WithField("path", filepath.Clean(event.Name)).
					Info("file modified")

				watcher.onChange(watcher)
				watcher.onAll(watcher)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				// Log
				watcher.log.
					WithError(err).
					Warn("watch error")
			}
		}
	}()
	<-done
}

func (watcher *Watcher) AddGroup(group string, name string) error {
	if err := watcher.Add(name); err != nil {
		return err
	}

	watcher.groups[group] = append(watcher.groups[group], name)

	return nil
}

func (watcher *Watcher) ReplaceGroup(group string, names []string) error {
	if _, ok := watcher.groups[group]; ok {
		for i, oldName := range watcher.groups[group] {
			found := false
			for j, newName := range names {
				if oldName == newName {
					// Remove new names already present in old names
					names = append(names[:j], names[j+1:]...)
					found = true
					break
				}
			}
			if !found {
				// Remove old names not present in new names
				if err := watcher.Remove(oldName); err != nil {
					return err
				}
				watcher.groups[group] = append(watcher.groups[group][:i], watcher.groups[group][i+1:]...)
			}
		}
	}

	for _, name := range names {
		// Add new names not present in old names
		if err := watcher.Add(name); err != nil {
			return err
		}
		watcher.groups[group] = append(watcher.groups[group], name)
	}

	return nil
}
