package watcher

import (
	"github.com/fsnotify/fsnotify"
	internalLog "manala/internal/log"
	internalReport "manala/internal/report"
)

func NewManager(log *internalLog.Logger) *Manager {
	return &Manager{
		log: log,
	}
}

type Manager struct {
	log *internalLog.Logger
}

func (manager *Manager) NewWatcher(onStart func(watcher *Watcher), onChange func(watcher *Watcher), onAll func(watcher *Watcher)) (*Watcher, error) {
	fsnotifyWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, internalReport.NewError(err).
			WithMessage("watcher error")
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
