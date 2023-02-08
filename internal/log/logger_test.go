package log

import (
	"bytes"
	"fmt"
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
	suite.Run(t, new(LoggerSuite))
}

func (s *LoggerSuite) SetupTest() {
	s.goldie = goldie.New(s.T())
}

func (s *LoggerSuite) Test() {
	out := &bytes.Buffer{}
	logger := New(out)

	logger.WithField("foo", "bar").Debug("debug")
	logger.WithField("foo", "bar").Info("info")
	logger.WithField("foo", "bar").Warn("warn")
	logger.WithField("foo", "bar").Error("error")

	s.goldie.Assert(s.T(), internalTesting.Path(s, "out"), out.Bytes())
}

func (s *LoggerSuite) TestLevelDebug() {
	out := &bytes.Buffer{}
	logger := New(out)

	logger.WithField("foo", "bar").Debug("debug")

	s.Empty(out.Bytes())

	logger.LevelDebug()
	logger.WithField("foo", "bar").Debug("debug")

	s.goldie.Assert(s.T(), internalTesting.Path(s, "out"), out.Bytes())
}

func (s *LoggerSuite) TestReport() {
	s.Run("Padding", func() {
		out := &bytes.Buffer{}
		logger := New(out)

		logger.IncreasePadding()

		report := internalReport.NewReport("report")

		logger.Report(report)

		s.goldie.Assert(s.T(), internalTesting.Path(s, "out"), out.Bytes())
	})
	s.Run("Empty Message No Error", func() {
		out := &bytes.Buffer{}
		logger := New(out)

		report := internalReport.NewReport("")

		logger.Report(report)

		s.goldie.Assert(s.T(), internalTesting.Path(s, "out"), out.Bytes())
	})
	s.Run("Empty Message Error", func() {
		out := &bytes.Buffer{}
		logger := New(out)

		report := internalReport.NewReport("")
		report.Compose(
			internalReport.WithErr(fmt.Errorf("error")),
		)

		logger.Report(report)

		s.goldie.Assert(s.T(), internalTesting.Path(s, "out"), out.Bytes())
	})
	s.Run("Message No Error", func() {
		out := &bytes.Buffer{}
		logger := New(out)

		report := internalReport.NewReport("report")

		logger.Report(report)

		s.goldie.Assert(s.T(), internalTesting.Path(s, "out"), out.Bytes())
	})
	s.Run("Message Error", func() {
		out := &bytes.Buffer{}
		logger := New(out)

		report := internalReport.NewReport("report")
		report.Compose(
			internalReport.WithErr(fmt.Errorf("error")),
		)

		logger.Report(report)

		s.goldie.Assert(s.T(), internalTesting.Path(s, "out"), out.Bytes())
	})
	s.Run("Fields", func() {
		out := &bytes.Buffer{}
		logger := New(out)

		report := internalReport.NewReport("report")
		report.Compose(
			internalReport.WithField("foo", "foo"),
			internalReport.WithField("bar", "bar"),
		)

		logger.Report(report)

		s.goldie.Assert(s.T(), internalTesting.Path(s, "out"), out.Bytes())
	})
	s.Run("Trace", func() {
		out := &bytes.Buffer{}
		logger := New(out)

		report := internalReport.NewReport("report")
		report.Compose(
			internalReport.WithTrace("trace"),
		)

		logger.Report(report)

		s.goldie.Assert(s.T(), internalTesting.Path(s, "out"), out.Bytes())
	})
}

func (s *LoggerSuite) TestPadding() {
	out := &bytes.Buffer{}
	logger := New(out)

	logger.WithField("foo", "bar").Info("default")

	logger.IncreasePadding()
	logger.WithField("foo", "bar").Info("increase")

	logger.DecreasePadding()
	logger.WithField("foo", "bar").Info("decrease")

	logger.IncreasePadding()
	logger.IncreasePadding()
	logger.WithField("foo", "bar").Info("before without padding")
	logger.WithoutPadding().WithField("foo", "bar").Info("without padding")
	logger.WithField("foo", "bar").Info("after without padding")

	s.goldie.Assert(s.T(), internalTesting.Path(s, "out"), out.Bytes())
}
