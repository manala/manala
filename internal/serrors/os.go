package serrors

import (
	"errors"
	"os"
)

func NewOs(err error) Error {
	message := err.Error()
	arguments := []any{}

	var (
		_pathError    *os.PathError
		_syscallError *os.SyscallError
	)

	switch {
	case errors.As(err, &_pathError):
		message = _pathError.Err.Error()
		arguments = append(arguments,
			"operation", _pathError.Op,
			"path", _pathError.Path,
		)
		if _pathError.Timeout() {
			arguments = append(arguments, "timeout", true)
		}
	case errors.As(err, &_syscallError):
		message = _syscallError.Err.Error()
		arguments = append(arguments, "syscall", _syscallError.Syscall)
		if _syscallError.Timeout() {
			arguments = append(arguments, "timeout", true)
		}
	}

	return New(message).
		WithArguments(arguments...)
}
