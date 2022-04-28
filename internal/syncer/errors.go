package syncer

import internalErrors "manala/internal/errors"

func SourceNotExistError(path string) *internalErrors.InternalError {
	return internalErrors.New("no source file or directory").
		WithField("path", path)
}
