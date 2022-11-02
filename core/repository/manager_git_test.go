package repository

import (
	"github.com/stretchr/testify/suite"
	"io"
	"manala/core"
	internalCache "manala/internal/cache"
	internalLog "manala/internal/log"
	internalTesting "manala/internal/testing"
	"os"
	"path/filepath"
	"testing"
)

type GitManagerSuite struct{ suite.Suite }

func TestGitManagerSuite(t *testing.T) {
	suite.Run(t, new(GitManagerSuite))
}

func (s *GitManagerSuite) TestLoadRepository() {
	cacheDir := internalTesting.DataPath(s, "cache")

	_ = os.RemoveAll(cacheDir)

	manager := NewGitManager(
		internalLog.New(io.Discard),
		internalCache.New(cacheDir),
	)

	s.Run("Default", func() {
		repoUrl := "https://github.com/octocat/Hello-World.git"

		repo, err := manager.LoadRepository(repoUrl)

		s.NoError(err)
		s.Equal(repoUrl, repo.Url())
		s.DirExists(filepath.Join(cacheDir, "repositories", "3a1b4df2cfc5d2f2cc3259ef67cda77786d47f84"))
	})

	s.Run("Empty", func() {
		repo, err := manager.LoadRepository("")

		var _unsupportedRepositoryError *core.UnsupportedRepositoryError
		s.ErrorAs(err, &_unsupportedRepositoryError)

		s.EqualError(err, "unsupported repository url")
		s.Nil(repo)
	})

	s.Run("Unsupported", func() {
		repo, err := manager.LoadRepository(
			"foo",
		)

		var _unsupportedRepositoryError *core.UnsupportedRepositoryError
		s.ErrorAs(err, &_unsupportedRepositoryError)

		s.EqualError(err, "unsupported repository url")
		s.Nil(repo)
	})

	s.Run("Error", func() {
		repo, err := manager.LoadRepository(
			"https://github.com/octocat/Foo-Bar.git",
		)

		s.EqualError(err, "authentication required")
		s.Nil(repo)
	})
}
