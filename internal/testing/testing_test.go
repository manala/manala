package testing

import (
	"github.com/stretchr/testify/suite"
	"path/filepath"
	"testing"
)

type TestingSuite struct{ suite.Suite }

func TestTestingSuite(t *testing.T) {
	suite.Run(t, new(TestingSuite))
}

func (s *TestingSuite) TestPath() {
	s.Equal(filepath.Join("TestTestingSuite", "TestPath"), Path(s))
	s.Equal(filepath.Join("TestTestingSuite", "TestPath", "foo"), Path(s, "foo"))
	s.Equal(filepath.Join("TestTestingSuite", "TestPath", "foo", "bar"), Path(s, "foo", "bar"))
}

func (s *TestingSuite) TestDataPath() {
	s.Equal(filepath.Join("testdata", "TestTestingSuite", "TestDataPath"), DataPath(s))
	s.Equal(filepath.Join("testdata", "TestTestingSuite", "TestDataPath", "foo"), DataPath(s, "foo"))
	s.Equal(filepath.Join("testdata", "TestTestingSuite", "TestDataPath", "foo", "bar"), DataPath(s, "foo", "bar"))
}
