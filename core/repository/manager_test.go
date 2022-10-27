package repository

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io"
	"manala/core"
	internalLog "manala/internal/log"
	internalTesting "manala/internal/testing"
	"os"
	"path/filepath"
	"testing"
)

/***********/
/* Default */
/***********/

type DefaultManagerSuite struct{ suite.Suite }

func TestDefaultManagerSuite(t *testing.T) {
	suite.Run(t, new(DefaultManagerSuite))
}

func (s *DefaultManagerSuite) TestLoadRepository() {
	log := internalLog.New(io.Discard)

	s.Run("Default", func() {
		repoMock := core.NewRepositoryMock()

		managerMock := core.NewRepositoryManagerMock()
		managerMock.
			On("LoadRepository", mock.Anything).Return(repoMock, nil)

		manager := NewDefaultManager(
			log,
			"default",
			managerMock,
		)

		repo, err := manager.LoadRepository("foo")

		s.NoError(err)
		s.Equal(repoMock, repo)

		managerMock.AssertCalled(s.T(), "LoadRepository", "foo")
	})

	s.Run("Empty", func() {
		repoMock := core.NewRepositoryMock()

		managerMock := core.NewRepositoryManagerMock()
		managerMock.
			On("LoadRepository", mock.Anything).Return(core.NewRepositoryMock(), nil)

		manager := NewDefaultManager(
			log,
			"default",
			managerMock,
		)

		repo, err := manager.LoadRepository("")

		s.NoError(err)
		s.Equal(repoMock, repo)

		managerMock.AssertCalled(s.T(), "LoadRepository", "default")
	})
}

/*********/
/* Cache */
/*********/

type CacheManagerSuite struct{ suite.Suite }

func TestCacheManagerSuite(t *testing.T) {
	suite.Run(t, new(CacheManagerSuite))
}

func (s *CacheManagerSuite) TestLoadRepository() {
	log := internalLog.New(io.Discard)

	repoMock := core.NewRepositoryMock()

	managerMock := core.NewRepositoryManagerMock()
	managerMock.
		On("LoadRepository", mock.Anything).Return(repoMock, nil)

	manager := NewCacheManager(
		log,
		managerMock,
	)

	// First call should pass to embedded manager
	repo, err := manager.LoadRepository("foo")

	s.NoError(err)
	s.Equal(repoMock, repo)

	managerMock.AssertNumberOfCalls(s.T(), "LoadRepository", 1)

	// Second call with same name should extract from cache, and not pass to embedded manager
	repo, err = manager.LoadRepository("foo")

	s.NoError(err)
	s.Equal(repoMock, repo)

	managerMock.AssertNumberOfCalls(s.T(), "LoadRepository", 1)

	// Third call with different name should pass to embedded manager
	repo, err = manager.LoadRepository("bar")

	s.NoError(err)
	s.Equal(repoMock, repo)

	managerMock.AssertNumberOfCalls(s.T(), "LoadRepository", 2)
}

/*********/
/* Cache */
/*********/

type ChainManagerSuite struct{ suite.Suite }

func TestChainManagerSuite(t *testing.T) {
	suite.Run(t, new(ChainManagerSuite))
}

func (s *ChainManagerSuite) TestLoadRepository() {
	log := internalLog.New(io.Discard)

	s.Run("Default", func() {
		path := internalTesting.DataPath(s, "repository")

		manager := NewChainManager(
			log,
			[]core.RepositoryManager{
				NewDirManager(log),
			},
		)

		repo, err := manager.LoadRepository(
			path,
		)

		s.NoError(err)
		s.Equal(path, repo.Path())
	})
}

/*******/
/* Dir */
/*******/

type DirManagerSuite struct{ suite.Suite }

func TestDirManagerSuite(t *testing.T) {
	suite.Run(t, new(DirManagerSuite))
}

func (s *DirManagerSuite) TestLoadRepository() {
	log := internalLog.New(io.Discard)

	manager := NewDirManager(log)

	s.Run("Default", func() {
		path := internalTesting.DataPath(s, "repository")

		repo, err := manager.LoadRepository(
			path,
		)

		s.NoError(err)
		s.Equal(path, repo.Path())
	})

	s.Run("Empty", func() {
		repo, err := manager.LoadRepository("")

		var _unsupportedRepositoryError *core.UnsupportedRepositoryError
		s.ErrorAs(err, &_unsupportedRepositoryError)

		s.EqualError(err, "unsupported repository")
		s.Nil(repo)
	})

	s.Run("Not Found", func() {
		path := internalTesting.DataPath(s, "repository")

		repo, err := manager.LoadRepository(
			path,
		)

		var _notFoundRepositoryError *core.NotFoundRepositoryError
		s.ErrorAs(err, &_notFoundRepositoryError)

		s.EqualError(err, "repository not found")
		s.Nil(repo)
	})

	s.Run("Wrong", func() {
		path := internalTesting.DataPath(s, "repository")

		repo, err := manager.LoadRepository(
			path,
		)

		s.EqualError(err, "wrong repository")
		s.Nil(repo)
	})
}

/*******/
/* Git */
/*******/

type GitManagerSuite struct{ suite.Suite }

func TestGitManagerSuite(t *testing.T) {
	suite.Run(t, new(GitManagerSuite))
}

func (s *GitManagerSuite) TestLoadRepository() {
	cacheDir := internalTesting.DataPath(s, "cache")

	_ = os.RemoveAll(cacheDir)

	log := internalLog.New(io.Discard)

	manager := NewGitManager(
		log,
		cacheDir,
	)

	s.Run("Default", func() {
		url := "https://github.com/octocat/Hello-World.git"

		repoPath := filepath.Join(cacheDir, "repositories", "3a1b4df2cfc5d2f2cc3259ef67cda77786d47f84")

		repo, err := manager.LoadRepository(
			url,
		)

		s.NoError(err)
		s.Equal(url, repo.Path())
		s.DirExists(repoPath)
	})

	s.Run("Empty", func() {
		repo, err := manager.LoadRepository("")

		var _unsupportedRepositoryError *core.UnsupportedRepositoryError
		s.ErrorAs(err, &_unsupportedRepositoryError)

		s.EqualError(err, "unsupported repository")
		s.Nil(repo)
	})

	s.Run("Unsupported", func() {
		repo, err := manager.LoadRepository(
			"foo",
		)

		var _unsupportedRepositoryError *core.UnsupportedRepositoryError
		s.ErrorAs(err, &_unsupportedRepositoryError)

		s.EqualError(err, "unsupported repository")
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
