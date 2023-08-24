package watcher

import (
	"github.com/stretchr/testify/suite"
	"io"
	"log/slog"
	"testing"
)

type WatcherManagerSuite struct {
	suite.Suite
	watcherManager *Manager
}

func TestWatcherManagerSuite(t *testing.T) {
	suite.Run(t, new(WatcherManagerSuite))
}

func (s *WatcherManagerSuite) SetupTest() {
	s.watcherManager = NewManager(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)
}

func (s *WatcherManagerSuite) TestNewWatcher() {
	watcher, err := s.watcherManager.NewWatcher(
		func(watcher *Watcher) {},
		func(watcher *Watcher) {},
		func(watcher *Watcher) {},
	)
	s.NoError(err)
	s.IsType(&Watcher{}, watcher)
}
