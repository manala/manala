package git

import (
	internalErrors "manala/internal/errors"
)

func CloneRepositoryUrlError(dir string, url string, err error) *internalErrors.InternalError {
	return internalErrors.New("clone git repository").
		WithField("dir", dir).
		WithField("url", url).
		WithError(err)
}

func InvalidRepositoryError(dir string, err error) *internalErrors.InternalError {
	return internalErrors.New("invalid git repository").
		WithField("dir", dir).
		WithError(err)
}

func PullRepositoryError(dir string, err error) *internalErrors.InternalError {
	return internalErrors.New("pull git repository").
		WithField("dir", dir).
		WithError(err)
}

func DeleteRepositoryError(dir string, err error) *internalErrors.InternalError {
	return internalErrors.New("delete git repository").
		WithField("dir", dir).
		WithError(err)
}

func OpenRepositoryError(dir string, err error) *internalErrors.InternalError {
	return internalErrors.New("open git repository").
		WithField("dir", dir).
		WithError(err)
}
