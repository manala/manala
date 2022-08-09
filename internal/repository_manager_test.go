package internal

import (
	"github.com/stretchr/testify/suite"
	"io"
	internalLog "manala/internal/log"
	internalTesting "manala/internal/testing"
	"testing"
)

type RepositoryManagerSuite struct{ suite.Suite }

func TestRepositoryManagerSuite(t *testing.T) {
	suite.Run(t, new(RepositoryManagerSuite))
}

func (s *RepositoryManagerSuite) Test() {
	log := internalLog.New(io.Discard)

	s.Run("AddRepositoryLoader", func() {
		manager := NewRepositoryManager(
			log,
			"repository",
		)

		s.Empty(manager.RepositoryLoaders)

		manager.AddRepositoryLoader(
			&RepositoryDirLoader{Log: log},
		)

		s.Len(manager.RepositoryLoaders, 1)
	})

	s.Run("LoadRepository Default", func() {
		manager := NewRepositoryManager(
			log,
			internalTesting.DataPath(s, "repository"),
		)
		manager.AddRepositoryLoader(
			&RepositoryDirLoader{Log: log},
		)

		repository, err := manager.LoadRepository([]string{})

		s.NoError(err)
		s.Equal(internalTesting.DataPath(s, "repository"), repository.Path())
	})

	s.Run("LoadRepository", func() {
		manager := NewRepositoryManager(
			log,
			internalTesting.DataPath(s, "repository"),
		)
		manager.AddRepositoryLoader(
			&RepositoryDirLoader{Log: log},
		)

		repository, err := manager.LoadRepository([]string{
			internalTesting.DataPath(s, "repository"),
		})

		s.NoError(err)
		s.Equal(internalTesting.DataPath(s, "repository"), repository.Path())
	})
}
