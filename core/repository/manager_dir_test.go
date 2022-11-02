package repository

import (
	"github.com/stretchr/testify/suite"
	"io"
	"manala/core"
	internalLog "manala/internal/log"
	internalTesting "manala/internal/testing"
	"testing"
)

type DirManagerSuite struct{ suite.Suite }

func TestDirManagerSuite(t *testing.T) {
	suite.Run(t, new(DirManagerSuite))
}

func (s *DirManagerSuite) TestLoadRepository() {
	manager := NewDirManager(
		internalLog.New(io.Discard),
	)

	s.Run("Default", func() {
		repoUrl := internalTesting.DataPath(s, "repository")

		repo, err := manager.LoadRepository(repoUrl)

		s.NoError(err)
		s.Equal(repoUrl, repo.Url())
	})

	s.Run("Empty", func() {
		repo, err := manager.LoadRepository("")

		var _unsupportedRepositoryError *core.UnsupportedRepositoryError
		s.ErrorAs(err, &_unsupportedRepositoryError)

		s.EqualError(err, "unsupported empty repository url")
		s.Nil(repo)
	})

	s.Run("Not Found", func() {
		repoUrl := internalTesting.DataPath(s, "repository")

		repo, err := manager.LoadRepository(repoUrl)

		var _notFoundRepositoryError *core.NotFoundRepositoryError
		s.ErrorAs(err, &_notFoundRepositoryError)

		s.EqualError(err, "repository not found")
		s.Nil(repo)
	})

	s.Run("Wrong", func() {
		repoUrl := internalTesting.DataPath(s, "repository")

		repo, err := manager.LoadRepository(repoUrl)

		s.EqualError(err, "wrong repository")
		s.Nil(repo)
	})
}
