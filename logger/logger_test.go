package logger

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/suite"
	"manala/config"
	"testing"
)

/*********/
/* Suite */
/*********/

type LoggerTestSuite struct {
	suite.Suite
	conf *config.Config
	out  *bytes.Buffer
	log  *Logger
}

func TestLoggerTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(LoggerTestSuite))
}

func (s *LoggerTestSuite) SetupSuite() {
	s.conf = config.New("foo", "bar")

	s.out = bytes.NewBufferString("")

	s.log = New(s.conf)
	s.log.SetOut(s.out)
}

func (s *LoggerTestSuite) SetupTest() {
	s.conf.SetDebug(false)
	s.out.Reset()
}

/*********/
/* Tests */
/*********/

func (s *LoggerTestSuite) TestDebug() {
	s.conf.SetDebug(false)
	s.log.Debug("foo")
	s.Empty(s.out.String())

	s.conf.SetDebug(true)
	s.log.Debug("foo")
	s.Equal(`   • foo                      
`, s.out.String())
}

func (s *LoggerTestSuite) TestDebugWithField() {
	s.conf.SetDebug(false)
	s.log.DebugWithField("foo", "bar", "baz")
	s.Empty(s.out.String())

	s.conf.SetDebug(true)
	s.log.DebugWithField("foo", "bar", "baz")
	s.Equal(`   • foo                       bar=baz
`, s.out.String())
}

func (s *LoggerTestSuite) TestDebugWithFields() {
	s.conf.SetDebug(false)
	s.log.DebugWithFields("foo", Fields{
		"bar": "baz",
		"qux": "quux",
	})
	s.Empty(s.out.String())

	s.conf.SetDebug(true)
	s.log.DebugWithFields("foo", Fields{
		"bar": "baz",
		"qux": "quux",
	})
	s.Equal(`   • foo                       bar=baz qux=quux
`, s.out.String())
}

func (s *LoggerTestSuite) TestInfo() {
	s.log.Info("foo")
	s.Equal(`   • foo                      
`, s.out.String())
}

func (s *LoggerTestSuite) TestInfoWithField() {
	s.log.InfoWithField("foo", "bar", "baz")
	s.Equal(`   • foo                       bar=baz
`, s.out.String())
}

func (s *LoggerTestSuite) TestInfoWithFields() {
	s.log.InfoWithFields("foo", Fields{
		"bar": "baz",
		"qux": "quux",
	})
	s.Equal(`   • foo                       bar=baz qux=quux
`, s.out.String())
}

func (s *LoggerTestSuite) TestError() {
	s.log.Info("foo")
	s.Equal(`   • foo                      
`, s.out.String())
}

func (s *LoggerTestSuite) TestErrorWithError() {
	s.log.ErrorWithError("foo", errors.New("bar"))
	s.Equal(`   ⨯ foo                       error=bar
`, s.out.String())
}
