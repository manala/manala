package serrors

import (
	"errors"
	"os"
)

func FromOs(err error) Error {
	var arguments []any
	message := err.Error()

	// Path error
	if err, ok := errors.AsType[*os.PathError](err); ok {
		switch {
		case errors.Is(err.Err, os.ErrNotExist):
			message = os.ErrNotExist.Error()
		default:
			message = err.Err.Error()
		}
		arguments = append(arguments,
			"operation", err.Op,
			"path", err.Path,
		)

		if err.Timeout() {
			arguments = append(arguments, "timeout", true)
		}
	}

	// Syscall error
	if err, ok := errors.AsType[*os.SyscallError](err); ok {
		message = err.Err.Error()
		arguments = append(arguments, "syscall", err.Syscall)

		if err.Timeout() {
			arguments = append(arguments, "timeout", true)
		}
	}

	return New(message).
		WithArguments(arguments...)
}
