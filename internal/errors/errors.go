package errors

import (
	"github.com/caarlos0/log"
)

func New(message string) *InternalError {
	err := &InternalError{
		Message: message,
		Fields:  make(log.Fields),
		Errs:    make([]*InternalError, 0),
	}

	return err
}

type InternalError struct {
	Message string
	Fields  log.Fields
	Err     error
	Errs    []*InternalError
	Trace   string
}

func (err *InternalError) Error() string {
	return err.Message
}

func (err *InternalError) With(message string) *InternalError {
	err.Message = message

	return err
}

func (err *InternalError) WithField(key string, value interface{}) *InternalError {
	err.Fields[key] = value

	return err
}

func (err *InternalError) WithError(error error) *InternalError {
	err.Err = error

	return err
}

func (err *InternalError) WithErrors(errors []*InternalError) *InternalError {
	err.Errs = errors

	return err
}

func (err *InternalError) WithTrace(trace string) *InternalError {
	err.Trace = trace

	return err
}
