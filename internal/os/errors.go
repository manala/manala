package os

import (
	"errors"
	internalReport "manala/internal/report"
	"os"
)

func NewError(err error) *Error {
	return &Error{
		error: err,
	}
}

type Error struct {
	error
}

func (err *Error) Unwrap() error {
	return err.error
}

func (err *Error) Report(report *internalReport.Report) {
	report.Compose(
		internalReport.WithErr(err.error),
	)

	var pathError *os.PathError
	if errors.As(err.error, &pathError) {
		report.Compose(
			internalReport.WithErr(pathError.Err),
			internalReport.WithField("operation", pathError.Op),
		)
		if pathError.Timeout() {
			report.Compose(
				internalReport.WithField("timeout", true),
			)
		}
	}

	var syscallError *os.SyscallError
	if errors.As(err.error, &syscallError) {
		report.Compose(
			internalReport.WithErr(syscallError.Err),
			internalReport.WithField("syscall", syscallError.Syscall),
		)
		if syscallError.Timeout() {
			report.Compose(
				internalReport.WithField("timeout", true),
			)
		}
	}
}
