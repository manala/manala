package log

import (
	"fmt"
	"github.com/caarlos0/log"
	"io"
	internalReport "manala/internal/report"
)

func New(out io.Writer) *Logger {
	return &Logger{
		Logger: log.New(out),
	}
}

type Logger struct {
	*log.Logger
}

func (logger *Logger) LevelDebug() {
	logger.Level = log.DebugLevel
}

func (logger *Logger) Report(report *internalReport.Report) {
	// Reset padding
	_padding := logger.Padding
	logger.ResetPadding()

	logger.report(report)

	// Restore padding
	logger.Padding = _padding
}

func (logger *Logger) report(report *internalReport.Report) {
	if report.Message() != "" {
		logger.Error(report.Message())
		if report.Err() != nil {
			logger.WithError(report.Err())
		}
	} else {
		if report.Err() != nil {
			logger.Error(report.Err().Error())
		} else {
			logger.Error("")
		}
	}

	// Fields
	fields := log.Fields{}
	for k, v := range report.Fields() {
		fields[k] = v
	}
	logger.WithFields(fields)

	// Errors
	if len(report.Reports()) != 0 {
		logger.IncreasePadding()
		for _, rep := range report.Reports() {
			logger.report(rep)
		}
		logger.DecreasePadding()
	}

	// Trace
	if report.Trace() != "" {
		_, _ = fmt.Fprint(logger.Writer, "\n", report.Trace(), "\n")
	}
}
