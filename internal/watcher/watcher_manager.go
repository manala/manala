package watcher

import (
	"github.com/fsnotify/fsnotify"
	"log/slog"
)

func NewManager(log *slog.Logger) *Manager {
	return &Manager{
		log: log,
	}
}

type Manager struct {
	log *slog.Logger
}

func (manager *Manager) NewWatcher(onStart func(watcher *Watcher), onChange func(watcher *Watcher), onAll func(watcher *Watcher)) (*Watcher, error) {
	fsnotifyWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &Watcher{
		log:      manager.log,
		Watcher:  fsnotifyWatcher,
		onStart:  onStart,
		onChange: onChange,
		onAll:    onAll,
		groups:   map[string][]string{},
	}, nil
}
