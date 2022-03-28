package models

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/discard"
	"github.com/stretchr/testify/suite"
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
	logger := &log.Logger{
		Handler: discard.Default,
	}

	s.manager = NewWatcherManager(logger)
}

/*********/
/* Tests */
/*********/

func (s *WatcherTestSuite) Test() {
	watcher, err := s.manager.NewWatcher()
	s.NoError(err)
	s.Implements((*WatcherInterface)(nil), watcher)
}
