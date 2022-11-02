package watcher

import (
	"bytes"
	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/suite"
	internalLog "manala/internal/log"
	internalTesting "manala/internal/testing"
	"path/filepath"
	"testing"
)

type WatcherSuite struct {
	suite.Suite
	stderr *bytes.Buffer
	log    *internalLog.Logger
}

func TestWatcherSuite(t *testing.T) {
	suite.Run(t, new(WatcherSuite))
}

func (s *WatcherSuite) SetupTest() {
	s.stderr = &bytes.Buffer{}
	s.log = internalLog.New(s.stderr)
}

func (s *WatcherSuite) TestGroups() {
	path := internalTesting.DataPath(s)
	fooPath := filepath.Join(path, "foo")
	barPath := filepath.Join(path, "bar")
	bazPath := filepath.Join(path, "baz")

	fsnotifyWatcher, _ := fsnotify.NewWatcher()
	watcher := &Watcher{
		log:     s.log,
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
