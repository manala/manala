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

	var _pathError *os.PathError
	if errors.As(err.error, &_pathError) {
		report.Compose(
			internalReport.WithErr(_pathError.Err),
			internalReport.WithField("operation", _pathError.Op),
		)
		if _pathError.Timeout() {
			report.Compose(
				internalReport.WithField("timeout", true),
			)
		}
	}

	var _syscallError *os.SyscallError
	if errors.As(err.error, &_syscallError) {
		report.Compose(
			internalReport.WithErr(_syscallError.Err),
			internalReport.WithField("syscall", _syscallError.Syscall),
		)
		if _syscallError.Timeout() {
			report.Compose(
				internalReport.WithField("timeout", true),
			)
		}
	}
}
