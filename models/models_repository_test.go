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
	src string
	dir string
}

func TestRepositoryTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(RepositoryTestSuite))
}

func (s *RepositoryTestSuite) SetupTest() {
	s.src = "foo"
	s.dir = "bar"
}

/**********************/
/* Repository - Tests */
/**********************/

func (s *RepositoryTestSuite) TestRepository() {
	repo := NewRepository(s.src, s.dir)
	s.Implements((*RepositoryInterface)(nil), repo)
	s.Equal(s.src, repo.Src())
	s.Equal(s.dir, repo.Dir())
}
