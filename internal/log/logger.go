package log

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"io"
	internalErrors "manala/internal/errors"
	"regexp"
)

func New(out io.Writer) *Logger {
	logger := &Logger{
		handler: cli.New(out),
	}
	logger.padding = logger.handler.Padding
	logger.Logger = &log.Logger{
		Handler: logger.handler,
		Level:   log.InfoLevel,
	}

	return logger
}

type Logger struct {
	handler *cli.Handler
	padding int
	*log.Logger
}

func (logger *Logger) LevelDebug() {
	logger.Level = log.DebugLevel
}

func (logger *Logger) LogError(err error) {
	// Reset padding
	padding := logger.handler.Padding
	logger.handler.Padding = logger.padding

	logger.logError(err)

	// Restore padding
	logger.handler.Padding = padding
}

var ansiCodesRegex = regexp.MustCompile("[\u001B\u009B][[\\]()#;?]*(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007|(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~])")

func (logger *Logger) CaptureError(err error) string {
	// Capture writer
	writer := logger.handler.Writer
	buffer := bytes.NewBufferString("")
	logger.handler.Writer = buffer

	logger.LogError(err)

	// Restore writer
	logger.handler.Writer = writer

	return ansiCodesRegex.ReplaceAllString(buffer.String(), "")
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
		logger.PaddingUp()
		for _, err := range _err.Errs {
			logger.logError(err)
		}
		logger.PaddingDown()
	}

	// Trace
	if _err.Trace != "" {
		_, _ = fmt.Fprint(logger.handler.Writer, "\n", _err.Trace, "\n")
	}
}

func (logger *Logger) PaddingUp() {
	logger.handler.Padding = logger.handler.Padding + 3
}

func (logger *Logger) PaddingDown() {
	logger.handler.Padding = logger.handler.Padding - 3
}
