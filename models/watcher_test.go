package models

import (
	"github.com/stretchr/testify/suite"
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
	log := logger.New(logger.WithDiscardment())

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
