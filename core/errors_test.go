package core

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type ErrorsSuite struct{ suite.Suite }

func TestErrorsSuite(t *testing.T) {
	suite.Run(t, new(ErrorsSuite))
}

func (s *ErrorsSuite) TestProject() {
	s.Run("NotFoundProjectManifestError", func() {
		err := NewNotFoundProjectManifestError("error")

		var _notFoundProjectManifestError *NotFoundProjectManifestError
		s.ErrorAs(err, &_notFoundProjectManifestError)

		s.EqualError(err, "error")
	})
}

func (s *ErrorsSuite) TestRecipe() {
	s.Run("NotFoundRecipeManifestError", func() {
		err := NewNotFoundRecipeManifestError("error")

		var _notFoundRecipeManifestError *NotFoundRecipeManifestError
		s.ErrorAs(err, &_notFoundRecipeManifestError)

		s.EqualError(err, "error")
	})
}

func (s *ErrorsSuite) TestRepository() {
	s.Run("NotFoundRepositoryError", func() {
		err := NewNotFoundRepositoryError("error")

		var _notFoundRepositoryError *NotFoundRepositoryError
		s.ErrorAs(err, &_notFoundRepositoryError)

		s.EqualError(err, "error")
	})

	s.Run("UnsupportedRepositoryError", func() {
		err := NewUnsupportedRepositoryError("error")

		var _unsupportedRepositoryError *UnsupportedRepositoryError
		s.ErrorAs(err, &_unsupportedRepositoryError)

		s.EqualError(err, "error")
	})
}
