package std

import (
	"errors"
	"os"

	"github.com/manala/manala/internal/errors/serror"
)

func From(err error) serror.Error {
	var arguments []any
	message := err.Error()

	// OS Path error
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

	// OS Syscall error
	if err, ok := errors.AsType[*os.SyscallError](err); ok {
		message = err.Err.Error()
		arguments = append(arguments, "syscall", err.Syscall)

		if err.Timeout() {
			arguments = append(arguments, "timeout", true)
		}
	}

	return serror.New(message).
		With(arguments...)
}
