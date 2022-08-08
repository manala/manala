package log

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/caarlos0/log"
	"io"
	internalErrors "manala/internal/errors"
	"regexp"
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

func (logger *Logger) LogError(err error) {
	// Reset padding
	_padding := logger.Padding
	logger.ResetPadding()

	logger.logError(err)

	// Restore padding
	logger.Padding = _padding
}

var ansiCodesRegex = regexp.MustCompile("[\u001B\u009B][[\\]()#;?]*(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007|(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~])")

func (logger *Logger) CaptureError(err error) []byte {
	// Capture writer
	_writer := logger.Writer
	buffer := &bytes.Buffer{}
	logger.Writer = buffer

	logger.LogError(err)

	// Restore writer
	logger.Writer = _writer

	return ansiCodesRegex.ReplaceAll(buffer.Bytes(), []byte{})
}

func (logger *Logger) logError(err error) {
	var _err *internalErrors.InternalError

	if !errors.As(err, &_err) {
		// Not internal error
		logger.Error(err.Error())
		return
	}

	logger.
		WithFields(_err.Fields).
		Error(_err.Error())

	// Error
	if _err.Err != nil {
		logger.WithError(_err.Err)
	}

	// Errors
	if len(_err.Errs) != 0 {
		logger.IncreasePadding()
		for _, err := range _err.Errs {
			logger.logError(err)
		}
		logger.DecreasePadding()
	}

	// Trace
	if _err.Trace != "" {
		_, _ = fmt.Fprint(logger.Writer, "\n", _err.Trace, "\n")
	}
}
