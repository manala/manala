package log

import (
	"bytes"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/suite"
	internalReport "manala/internal/report"
	internalTesting "manala/internal/testing"
	"testing"
)

type LoggerSuite struct {
	suite.Suite
	goldie *goldie.Goldie
}

func TestLoggerSuite(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	suite.Run(t, new(LoggerSuite))
}

func (s *LoggerSuite) SetupTest() {
	s.goldie = goldie.New(s.T())
}

func (s *LoggerSuite) Test() {
	out := &bytes.Buffer{}
	logger := New(out)

	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warn")
	logger.Error("error")

	s.goldie.Assert(s.T(), internalTesting.Path(s, "out"), out.Bytes())
}

func (s *LoggerSuite) TestLevelDebug() {
	out := &bytes.Buffer{}
	logger := New(out)

	logger.Debug("debug")

	s.Empty(out.Bytes())

	logger.LevelDebug()
	logger.Debug("debug")

	s.goldie.Assert(s.T(), internalTesting.Path(s, "out"), out.Bytes())
}

func (s *LoggerSuite) TestLogError() {
	out := &bytes.Buffer{}
	logger := New(out)

	logger.IncreasePadding()
	report := internalReport.NewReport("error")
	logger.Report(report)

	s.goldie.Assert(s.T(), internalTesting.Path(s, "out"), out.Bytes())
}

func (s *LoggerSuite) TestPadding() {
	out := &bytes.Buffer{}
	logger := New(out)

	logger.Info("info")

	s.goldie.Assert(s.T(), internalTesting.Path(s, "out"), out.Bytes())

	logger.IncreasePadding()
	logger.Info("info")

	s.goldie.Assert(s.T(), internalTesting.Path(s, "out_increase"), out.Bytes())

	logger.DecreasePadding()
	logger.Info("info")

	s.goldie.Assert(s.T(), internalTesting.Path(s, "out_decrease"), out.Bytes())
}
