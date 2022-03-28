package models

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

/**********************/
/* Repository - Suite */
/**********************/

type RepositoryTestSuite struct {
	suite.Suite
}

func TestRepositoryTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(RepositoryTestSuite))
}

/**********************/
/* Repository - Tests */
/**********************/

func (s *RepositoryTestSuite) TestRepository() {
	source := "foo"
	dir := "bar"

	s.Run("New", func() {
		repo := NewRepository(source, dir)
		s.Implements((*RepositoryInterface)(nil), repo)
		s.Equal(source, repo.Source())
		s.Equal(dir, repo.getDir())
	})
}
