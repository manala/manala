package watcher

import (
	"bytes"
	"github.com/stretchr/testify/suite"
	internalLog "manala/internal/log"
	"testing"
)

type WatcherManagerSuite struct {
	suite.Suite
	stderr         *bytes.Buffer
	watcherManager *WatcherManager
}

func TestWatcherManagerSuite(t *testing.T) {
	suite.Run(t, new(WatcherManagerSuite))
}

func (s *WatcherManagerSuite) SetupTest() {
	s.stderr = bytes.NewBufferString("")
	s.watcherManager = &WatcherManager{
		Log: internalLog.New(s.stderr),
	}
}

func (s *WatcherManagerSuite) TestNewWatcher() {
	watcher, err := s.watcherManager.NewWatcher(
		func(watcher *Watcher) { return },
		func(watcher *Watcher) { return },
		func(watcher *Watcher) { return },
	)
	s.NoError(err)
	s.IsType(&Watcher{}, watcher)
}
