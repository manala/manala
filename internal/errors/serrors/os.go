package serrors

import (
	"errors"
	"os"
)

func NewOs(err error) *OsError {
	_err := &OsError{
		message:   err.Error(),
		Arguments: NewArguments(),
	}

	var (
		_pathError    *os.PathError
		_syscallError *os.SyscallError
	)

	switch {
	case errors.As(err, &_pathError):
		_err.message = _pathError.Err.Error()
		_err.AppendArguments(
			"operation", _pathError.Op,
			"path", _pathError.Path,
		)
		if _pathError.Timeout() {
			_err.AppendArguments("timeout", true)
		}
	case errors.As(err, &_syscallError):
		_err.message = _syscallError.Err.Error()
		_err.AppendArguments("syscall", _syscallError.Syscall)
		if _syscallError.Timeout() {
			_err.AppendArguments("timeout", true)
		}
	}

	return _err
}

type OsError struct {
	message string
	*Arguments
}

func (err *OsError) Error() string {
	return err.message
}

func WrapOs(message string, err error) *WrapOsError {
	return &WrapOsError{
		WrapError: Wrap(
			message,
			NewOs(err),
		),
	}
}

type WrapOsError struct {
	*WrapError
}

func (err *WrapOsError) Unwrap() error {
	return err.err
}

func (err *WrapOsError) WithArguments(arguments ...any) *WrapOsError {
	_ = err.WrapError.WithArguments(arguments...)
	return err
}

func (err *WrapOsError) WithDetails(details string) *WrapOsError {
	_ = err.WrapError.WithDetails(details)
	return err
}
