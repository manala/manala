package models

import (
	"bytes"
	"github.com/stretchr/testify/suite"
	"manala/config"
	"manala/logger"
	"testing"
)

/*********/
/* Suite */
/*********/

type WatcherTestSuite struct {
	suite.Suite
	manager WatcherManagerInterface
}

func TestWatcherTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(WatcherTestSuite))
}

func (s *WatcherTestSuite) SetupTest() {
	conf := config.New("test", "foo")

	log := logger.New(conf)
	log.SetOut(bytes.NewBufferString(""))

	s.manager = NewWatcherManager(log)
}

/*********/
/* Tests */
/*********/

func (s *WatcherTestSuite) Test() {
	watcher, err := s.manager.NewWatcher()
	s.NoError(err)
	s.Implements((*WatcherInterface)(nil), watcher)
}
