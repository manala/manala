package os

import (
	"errors"
	internalErrors "manala/internal/errors"
	"os"
)

func FileSystemError(err error) *internalErrors.InternalError {
	_err := internalErrors.New("file system error").
		WithError(err)

	var pathError *os.PathError
	if errors.As(err, &pathError) {
		_ = _err.
			WithError(pathError.Err).
			WithField("operation", pathError.Op).
			WithField("path", pathError.Path)
	}

	return _err
}
