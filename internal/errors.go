package internal

import (
	internalErrors "manala/internal/errors"
)

/***********/
/* Project */
/***********/

type NotFoundProjectManifestError struct{ *internalErrors.InternalError }

func (err *NotFoundProjectManifestError) Unwrap() error {
	return err.InternalError
}

func NotFoundProjectManifestPathError(path string) *NotFoundProjectManifestError {
	return &NotFoundProjectManifestError{
		internalErrors.New("project manifest not found").
			WithField("path", path),
	}
}

func WrongProjectManifestPathError(path string) *internalErrors.InternalError {
	return internalErrors.New("wrong project manifest").
		WithField("path", path)
}

func EmptyProjectManifestPathError(path string) *internalErrors.InternalError {
	return internalErrors.New("empty project manifest").
		WithField("path", path)
}

func ValidationProjectManifestPathError(path string, err error, errs []*internalErrors.InternalError) *internalErrors.InternalError {
	return internalErrors.New("project validation error").
		WithField("path", path).
		WithError(err).
		WithErrors(errs)
}

/**************/
/* Repository */
/**************/

type UnsupportedRepositoryError struct{ *internalErrors.InternalError }

func (err *UnsupportedRepositoryError) Unwrap() error {
	return err.InternalError
}

func UnsupportedRepositoryPathError() *UnsupportedRepositoryError {
	return &UnsupportedRepositoryError{
		internalErrors.New("unsupported repository"),
	}
}

type NotFoundRepositoryError struct{ *internalErrors.InternalError }

func (err *NotFoundRepositoryError) Unwrap() error {
	return err.InternalError
}

func NotFoundRepositoryDirError(path string) *NotFoundRepositoryError {
	return &NotFoundRepositoryError{
		internalErrors.New("repository not found").
			WithField("path", path),
	}
}

func NotFoundRepositoryGitError(path string) *NotFoundRepositoryError {
	return &NotFoundRepositoryError{
		internalErrors.New("repository not found").
			WithField("path", path),
	}
}

func WrongRepositoryDirError(dir string) *internalErrors.InternalError {
	return internalErrors.New("wrong repository").
		WithField("dir", dir)
}

func EmptyRepositoryDirError(dir string) *internalErrors.InternalError {
	return internalErrors.New("empty repository").
		WithField("dir", dir)
}

/**********/
/* Recipe */
/**********/

type NotFoundRecipeManifestError struct{ *internalErrors.InternalError }

func (err *NotFoundRecipeManifestError) Unwrap() error {
	return err.InternalError
}

func NotFoundRecipeManifestPathError(path string) *NotFoundRecipeManifestError {
	return &NotFoundRecipeManifestError{
		internalErrors.New("recipe manifest not found").
			WithField("path", path),
	}
}

func WrongRecipeManifestPathError(path string) *internalErrors.InternalError {
	return internalErrors.New("wrong recipe manifest").
		WithField("path", path)
}

func EmptyRecipeManifestPathError(path string) *internalErrors.InternalError {
	return internalErrors.New("empty recipe manifest").
		WithField("path", path)
}

func ValidationRecipeManifestPathError(path string, err error, errs []*internalErrors.InternalError) *internalErrors.InternalError {
	return internalErrors.New("recipe validation error").
		WithField("path", path).
		WithError(err).
		WithErrors(errs)
}

func ValidationRecipeManifestOptionError(err error, errs []*internalErrors.InternalError) *internalErrors.InternalError {
	return internalErrors.New("recipe option validation error").
		WithError(err).
		WithErrors(errs)
}
