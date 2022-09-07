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
		manager := NewChainManager(
			log,
			internalTesting.DataPath(s, "repository"),
			[]core.RepositoryManager{
				NewDirManager(log),
			},
		)

		repository, err := manager.LoadRepository([]string{})

		s.NoError(err)
		s.Equal(internalTesting.DataPath(s, "repository"), repository.Path())
	})

	s.Run("LoadRepository", func() {
		manager := NewChainManager(
			log,
			internalTesting.DataPath(s, "repository"),
			[]core.RepositoryManager{
				NewDirManager(log),
			},
		)

		repository, err := manager.LoadRepository([]string{
			internalTesting.DataPath(s, "repository"),
		})

		s.NoError(err)
		s.Equal(internalTesting.DataPath(s, "repository"), repository.Path())
	})
}

func (s *ManagerSuite) TestDir() {
	log := internalLog.New(io.Discard)
	loader := NewDirManager(log)

	s.Run("LoadRepository Empty", func() {
		repository, err := loader.LoadRepository([]string{})

		var _unsupportedRepositoryError *core.UnsupportedRepositoryError
		s.ErrorAs(err, &_unsupportedRepositoryError)

		s.EqualError(err, "unsupported repository")
		s.Nil(repository)
	})

	s.Run("LoadRepository Not Found", func() {
		repository, err := loader.LoadRepository([]string{
			internalTesting.DataPath(s, "repository"),
		})

		var _notFoundRepositoryError *core.NotFoundRepositoryError
		s.ErrorAs(err, &_notFoundRepositoryError)

		s.EqualError(err, "repository not found")
		s.Nil(repository)
	})

	s.Run("LoadRepository Wrong", func() {
		repository, err := loader.LoadRepository([]string{
			internalTesting.DataPath(s, "repository"),
		})

		s.EqualError(err, "wrong repository")
		s.Nil(repository)
	})

	s.Run("LoadRepository", func() {
		repository, err := loader.LoadRepository([]string{
			internalTesting.DataPath(s, "repository"),
		})

		s.NoError(err)
		s.Equal(internalTesting.DataPath(s, "repository"), repository.Path())
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
		repository, err := loader.LoadRepository([]string{})

		var _unsupportedRepositoryError *core.UnsupportedRepositoryError
		s.ErrorAs(err, &_unsupportedRepositoryError)

		s.EqualError(err, "unsupported repository")
		s.Nil(repository)
	})

	s.Run("LoadRepository Unsupported", func() {
		repository, err := loader.LoadRepository([]string{
			"foo",
		})

		var _unsupportedRepositoryError *core.UnsupportedRepositoryError
		s.ErrorAs(err, &_unsupportedRepositoryError)

		s.EqualError(err, "unsupported repository")
		s.Nil(repository)
	})

	s.Run("LoadRepository Error", func() {
		repository, err := loader.LoadRepository([]string{
			"https://github.com/octocat/Foo-Bar.git",
		})

		s.EqualError(err, "authentication required")
		s.Nil(repository)
	})

	s.Run("LoadRepository", func() {
		url := "https://github.com/octocat/Hello-World.git"

		repository, err := loader.LoadRepository([]string{
			url,
		})

		s.NoError(err)
		s.Equal(url, repository.Path())
		s.DirExists(filepath.Join(cacheDir, "repositories", "3a1b4df2cfc5d2f2cc3259ef67cda77786d47f84"))
	})
}
