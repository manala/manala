package watcher

import internalErrors "manala/internal/errors"

func Error(err error) *internalErrors.InternalError {
	return internalErrors.New("watcher error").
		WithError(err)
}
