package internal

import (
	"github.com/stretchr/testify/suite"
	"io"
	internalLog "manala/internal/log"
	"path/filepath"
	"testing"
)

type RepositoryManagerSuite struct{ suite.Suite }

func TestRepositoryManagerSuite(t *testing.T) {
	suite.Run(t, new(RepositoryManagerSuite))
}

var repositoryManagerTestPath = filepath.Join("testdata", "repository_manager")

func (s *RepositoryManagerSuite) Test() {
	log := internalLog.New(io.Discard)

	manager := NewRepositoryManager(
		log,
		filepath.Join(repositoryManagerTestPath, "repository_default"),
	)

	s.Empty(manager.RepositoryLoaders)

	manager.AddRepositoryLoader(
		&RepositoryDirLoader{Log: log},
	)

	s.Len(manager.RepositoryLoaders, 1)

	s.Run("LoadRepository Default", func() {
		repository, err := manager.LoadRepository([]string{})

		s.NoError(err)
		s.Equal(filepath.Join(repositoryManagerTestPath, "repository_default"), repository.Path())
	})

	s.Run("LoadRepository", func() {
		repository, err := manager.LoadRepository([]string{
			filepath.Join(repositoryManagerTestPath, "repository"),
		})

		s.NoError(err)
		s.Equal(filepath.Join(repositoryManagerTestPath, "repository"), repository.Path())
	})
}
