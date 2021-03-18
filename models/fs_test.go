package models

import (
	"github.com/stretchr/testify/suite"
	"manala/fs"
	"testing"
)

/*********/
/* Suite */
/*********/

type FsTestSuite struct {
	suite.Suite
	manager FsManagerInterface
}

func TestFsTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(FsTestSuite))
}

func (s *FsTestSuite) SetupTest() {
	s.manager = NewFsManager(fs.NewManager())
}

/*********/
/* Tests */
/*********/

func (s *FsTestSuite) Test() {
	model := &mock{"foo"}
	fsys := s.manager.NewModelFs(model)
	s.Implements((*fs.ReadWriteInterface)(nil), fsys)
}
