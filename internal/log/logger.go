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
	padding int
}

func (logger *Logger) LevelDebug() {
	logger.Level = log.DebugLevel
}

func (logger *Logger) ResetPadding() {
	logger.padding = logger.Logger.Padding
	logger.Logger.ResetPadding()
}

func (logger *Logger) RestorePadding() {
	logger.Logger.Padding = logger.padding
}

func (logger *Logger) Report(report *internalReport.Report) {
	// Reset padding
	logger.ResetPadding()

	logger.report(report)

	// Restore padding
	logger.RestorePadding()
}

func (logger *Logger) report(report *internalReport.Report) {
	// Fields
	fields := log.Fields{}
	for k, v := range report.Fields() {
		fields[k] = v
	}
	_logger := logger.WithFields(fields)

	if report.Message() != "" {
		if report.Err() != nil {
			_logger = _logger.WithError(report.Err())
		}
		_logger.Error(report.Message())
	} else {
		if report.Err() != nil {
			_logger.Error(report.Err().Error())
		} else {
			_logger.Error("")
		}
	}

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
