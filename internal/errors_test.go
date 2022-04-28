package internal

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	internalErrors "manala/internal/errors"
	"testing"
)

var internalError *internalErrors.InternalError

type ErrorsSuite struct{ suite.Suite }

func TestErrorsSuite(t *testing.T) {
	suite.Run(t, new(ErrorsSuite))
}

func (s *ErrorsSuite) TestProject() {

	s.Run("NotFoundProjectManifestPathError", func() {
		err := NotFoundProjectManifestPathError("path")

		var _err *NotFoundProjectManifestError
		s.ErrorAs(err, &_err)

		s.ErrorAs(err, &internalError)
		s.Equal("project manifest not found", internalError.Message)
		s.Equal("path", internalError.Fields["path"])
	})

	s.Run("WrongProjectManifestPathError", func() {
		err := WrongProjectManifestPathError("path")

		s.ErrorAs(err, &internalError)
		s.Equal("wrong project manifest", internalError.Message)
		s.Equal("path", internalError.Fields["path"])
	})

	s.Run("EmptyProjectManifestPathError", func() {
		err := EmptyProjectManifestPathError("path")

		s.ErrorAs(err, &internalError)
		s.Equal("empty project manifest", internalError.Message)
		s.Equal("path", internalError.Fields["path"])
	})

	s.Run("ValidationProjectManifestPathError", func() {
		_err := fmt.Errorf("error")
		var _errs []*internalErrors.InternalError
		err := ValidationProjectManifestPathError("path", _err, _errs)

		s.ErrorAs(err, &internalError)
		s.Equal("project validation error", internalError.Message)
		s.Equal("path", internalError.Fields["path"])
		s.Equal(_err, err.Err)
		s.Equal(_errs, err.Errs)
	})

}

func (s *ErrorsSuite) TestRepository() {

	s.Run("UnsupportedRepositoryPathError", func() {
		err := UnsupportedRepositoryPathError()

		var _err *UnsupportedRepositoryError
		s.ErrorAs(err, &_err)

		s.ErrorAs(err, &internalError)
		s.Equal("unsupported repository", internalError.Message)
	})

	s.Run("NotFoundRepositoryDirError", func() {
		err := NotFoundRepositoryDirError("path")

		var _err *NotFoundRepositoryError
		s.ErrorAs(err, &_err)

		s.ErrorAs(err, &internalError)
		s.Equal("repository not found", internalError.Message)
		s.Equal("path", internalError.Fields["path"])
	})

	s.Run("NotFoundRepositoryGitError", func() {
		err := NotFoundRepositoryGitError("path")

		var _err *NotFoundRepositoryError
		s.ErrorAs(err, &_err)

		s.ErrorAs(err, &internalError)
		s.Equal("repository not found", internalError.Message)
		s.Equal("path", internalError.Fields["path"])
	})

	s.Run("WrongRepositoryDirError", func() {
		err := WrongRepositoryDirError("dir")

		s.ErrorAs(err, &internalError)
		s.Equal("wrong repository", internalError.Message)
		s.Equal("dir", internalError.Fields["dir"])
	})

	s.Run("EmptyRepositoryDirError", func() {
		err := EmptyRepositoryDirError("dir")

		s.ErrorAs(err, &internalError)
		s.Equal("empty repository", internalError.Message)
		s.Equal("dir", internalError.Fields["dir"])
	})

}

func (s *ErrorsSuite) TestRecipe() {

	s.Run("NotFoundRecipeManifestPathError", func() {
		err := NotFoundRecipeManifestPathError("path")

		var _err *NotFoundRecipeManifestError
		s.ErrorAs(err, &_err)

		s.ErrorAs(err, &internalError)
		s.Equal("recipe manifest not found", internalError.Message)
		s.Equal("path", internalError.Fields["path"])
	})

	s.Run("WrongRecipeManifestPathError", func() {
		err := WrongRecipeManifestPathError("path")

		s.ErrorAs(err, &internalError)
		s.Equal("wrong recipe manifest", internalError.Message)
		s.Equal("path", internalError.Fields["path"])
	})

	s.Run("EmptyRecipeManifestPathError", func() {
		err := EmptyRecipeManifestPathError("path")

		s.ErrorAs(err, &internalError)
		s.Equal("empty recipe manifest", internalError.Message)
		s.Equal("path", internalError.Fields["path"])
	})

	s.Run("ValidationRecipeManifestPathError", func() {
		_err := fmt.Errorf("error")
		var _errs []*internalErrors.InternalError
		err := ValidationRecipeManifestPathError("path", _err, _errs)

		s.ErrorAs(err, &internalError)
		s.Equal("recipe validation error", internalError.Message)
		s.Equal("path", internalError.Fields["path"])
		s.Equal(_err, err.Err)
		s.Equal(_errs, err.Errs)
	})
}
