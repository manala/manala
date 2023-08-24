package watcher

import (
	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/suite"
	"io"
	"log/slog"
	"path/filepath"
	"testing"
)

type WatcherSuite struct{ suite.Suite }

func TestWatcherSuite(t *testing.T) {
	suite.Run(t, new(WatcherSuite))
}

func (s *WatcherSuite) TestGroups() {
	path := filepath.FromSlash("testdata/WatcherSuite/TestGroups")
	fooPath := filepath.Join(path, "foo")
	barPath := filepath.Join(path, "bar")
	bazPath := filepath.Join(path, "baz")

	fsnotifyWatcher, _ := fsnotify.NewWatcher()
	watcher := &Watcher{
		log:     slog.New(slog.NewTextHandler(io.Discard, nil)),
		Watcher: fsnotifyWatcher,
		groups:  map[string][]string{},
	}

	s.Empty(watcher.groups)

	_ = watcher.AddGroup("group", fooPath)
	_ = watcher.AddGroup("group", barPath)

	s.Equal([]string{fooPath, barPath}, watcher.groups["group"])

	_ = watcher.ReplaceGroup("group", []string{barPath, bazPath})

	s.Equal([]string{barPath, bazPath}, watcher.groups["group"])
}
