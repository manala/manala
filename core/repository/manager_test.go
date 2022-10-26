package repository

import (
	"github.com/stretchr/testify/suite"
	"io"
	"manala/core"
	internalLog "manala/internal/log"
	internalTesting "manala/internal/testing"
	"os"
	"path/filepath"
	"testing"
)

type ManagerSuite struct{ suite.Suite }

func TestManagerSuite(t *testing.T) {
	suite.Run(t, new(ManagerSuite))
}

func (s *ManagerSuite) TestChain() {
	log := internalLog.New(io.Discard)

	s.Run("LoadRepository Default", func() {
		path := internalTesting.DataPath(s, "repository")

		manager := NewChainManager(
			log,
			path,
			[]core.RepositoryManager{
				NewDirManager(log),
			},
		)

		repo, err := manager.LoadRepository([]string{})

		s.NoError(err)
		s.Equal(path, repo.Path())
	})

	s.Run("LoadRepository", func() {
		path := internalTesting.DataPath(s, "repository")

		manager := NewChainManager(
			log,
			path,
			[]core.RepositoryManager{
				NewDirManager(log),
			},
		)

		repo, err := manager.LoadRepository([]string{
			path,
		})

		s.NoError(err)
		s.Equal(path, repo.Path())
	})
}

func (s *ManagerSuite) TestDir() {
	log := internalLog.New(io.Discard)
	loader := NewDirManager(log)

	s.Run("LoadRepository Empty", func() {
		repo, err := loader.LoadRepository([]string{})

		var _unsupportedRepositoryError *core.UnsupportedRepositoryError
		s.ErrorAs(err, &_unsupportedRepositoryError)

		s.EqualError(err, "unsupported repository")
		s.Nil(repo)
	})

	s.Run("LoadRepository Not Found", func() {
		path := internalTesting.DataPath(s, "repository")

		repo, err := loader.LoadRepository([]string{
			path,
		})

		var _notFoundRepositoryError *core.NotFoundRepositoryError
		s.ErrorAs(err, &_notFoundRepositoryError)

		s.EqualError(err, "repository not found")
		s.Nil(repo)
	})

	s.Run("LoadRepository Wrong", func() {
		path := internalTesting.DataPath(s, "repository")

		repo, err := loader.LoadRepository([]string{
			path,
		})

		s.EqualError(err, "wrong repository")
		s.Nil(repo)
	})

	s.Run("LoadRepository", func() {
		path := internalTesting.DataPath(s, "repository")

		repo, err := loader.LoadRepository([]string{
			path,
		})

		s.NoError(err)
		s.Equal(path, repo.Path())
	})
}

func (s *ManagerSuite) TestGit() {
	cacheDir := internalTesting.DataPath(s, "cache")

	_ = os.RemoveAll(cacheDir)

	log := internalLog.New(io.Discard)
	loader := NewGitManager(
		log,
		cacheDir,
	)

	s.Run("LoadRepository Empty", func() {
		repo, err := loader.LoadRepository([]string{})

		var _unsupportedRepositoryError *core.UnsupportedRepositoryError
		s.ErrorAs(err, &_unsupportedRepositoryError)

		s.EqualError(err, "unsupported repository")
		s.Nil(repo)
	})

	s.Run("LoadRepository Unsupported", func() {
		repo, err := loader.LoadRepository([]string{
			"foo",
		})

		var _unsupportedRepositoryError *core.UnsupportedRepositoryError
		s.ErrorAs(err, &_unsupportedRepositoryError)

		s.EqualError(err, "unsupported repository")
		s.Nil(repo)
	})

	s.Run("LoadRepository Error", func() {
		repo, err := loader.LoadRepository([]string{
			"https://github.com/octocat/Foo-Bar.git",
		})

		s.EqualError(err, "authentication required")
		s.Nil(repo)
	})

	s.Run("LoadRepository", func() {
		url := "https://github.com/octocat/Hello-World.git"

		repoPath := filepath.Join(cacheDir, "repositories", "3a1b4df2cfc5d2f2cc3259ef67cda77786d47f84")

		repo, err := loader.LoadRepository([]string{
			url,
		})

		s.NoError(err)
		s.Equal(url, repo.Path())
		s.DirExists(repoPath)
	})
}
