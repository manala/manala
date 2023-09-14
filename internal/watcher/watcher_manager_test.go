package watcher

import (
	"io"
	"log/slog"
)

func (s *Suite) TestManager() {
	watcherManager := NewManager(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)

	watcher, err := watcherManager.NewWatcher(
		func(watcher *Watcher) {},
		func(watcher *Watcher) {},
		func(watcher *Watcher) {},
	)
	s.NoError(err)
	s.IsType(&Watcher{}, watcher)
}
