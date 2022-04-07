package watcher

import (
	"bytes"
	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/suite"
	internalLog "manala/internal/log"
	"testing"
)

type WatcherSuite struct {
	suite.Suite
	stdErr *bytes.Buffer
	logger *internalLog.Logger
}

func TestWatcherSuite(t *testing.T) {
	suite.Run(t, new(WatcherSuite))
}

func (s *WatcherSuite) SetupTest() {
	s.stdErr = bytes.NewBufferString("")
	s.logger = internalLog.New(s.stdErr)
}

func (s *WatcherSuite) TestTemporaries() {
	fsnotifyWatcher, _ := fsnotify.NewWatcher()
	watcher := &Watcher{
		log:         s.logger,
		Watcher:     fsnotifyWatcher,
		temporaries: []string{},
	}

	s.Empty(watcher.temporaries)

	_ = watcher.AddTemporary("foo")
	_ = watcher.AddTemporary("bar")

	s.Equal([]string{"foo", "bar"}, watcher.temporaries)

	_ = watcher.RemoveTemporaries()

	s.Equal([]string{"foo", "bar"}, watcher.temporaries)
}
