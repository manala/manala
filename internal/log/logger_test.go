package log

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

type LoggerSuite struct{ suite.Suite }

func TestLoggerSuite(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	suite.Run(t, new(LoggerSuite))
}

func (s *LoggerSuite) Test() {
	out := &bytes.Buffer{}
	logger := New(out)

	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warn")
	logger.Error("error")

	s.Equal(`  • info
  • warn
  ⨯ error
`, out.String())
}

func (s *LoggerSuite) TestLevelDebug() {
	out := &bytes.Buffer{}
	logger := New(out)

	logger.Debug("debug")
	s.Empty(out.String())

	logger.LevelDebug()

	logger.Debug("debug")
	s.Equal(`  • debug
`, out.String())
}

func (s *LoggerSuite) TestLogError() {
	out := &bytes.Buffer{}
	logger := New(out)

	logger.IncreasePadding()
	logger.LogError(fmt.Errorf("error"))
	s.Equal(`  ⨯ error
`, out.String())
}

func (s *LoggerSuite) TestCaptureError() {
	out := &bytes.Buffer{}
	logger := New(out)

	logger.IncreasePadding()
	err := logger.CaptureError(fmt.Errorf("error"))
	s.Empty(out.String())
	s.Equal(`  ⨯ error
`, string(err))
}

func (s *LoggerSuite) TestPadding() {
	out := &bytes.Buffer{}
	logger := New(out)

	logger.Info("info")
	s.Equal(`  • info
`, out.String())

	logger.IncreasePadding()
	logger.Info("info")
	s.Equal(`  • info
    • info
`, out.String())

	logger.DecreasePadding()
	logger.Info("info")
	s.Equal(`  • info
    • info
  • info
`, out.String())
}
