package watcher

import (
	"github.com/fsnotify/fsnotify"
	internalLog "manala/internal/log"
)

type WatcherManager struct {
	Log *internalLog.Logger
}

func (manager *WatcherManager) NewWatcher(onStart func(watcher *Watcher), onChange func(watcher *Watcher), onAll func(watcher *Watcher)) (*Watcher, error) {
	fsnotifyWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, Error(err)
	}

	return &Watcher{
		log:         manager.Log,
		Watcher:     fsnotifyWatcher,
		onStart:     onStart,
		onChange:    onChange,
		onAll:       onAll,
		temporaries: []string{},
	}, nil
}