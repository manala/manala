package repository

import (
	"github.com/stretchr/testify/suite"
	"io"
	"manala/core"
	internalLog "manala/internal/log"
	internalTesting "manala/internal/testing"
	"testing"
)

type ChainManagerSuite struct{ suite.Suite }

func TestChainManagerSuite(t *testing.T) {
	suite.Run(t, new(ChainManagerSuite))
}

func (s *ChainManagerSuite) TestLoadRepository() {
	log := internalLog.New(io.Discard)

	s.Run("Default", func() {
		repoUrl := internalTesting.DataPath(s, "repository")

		manager := NewChainManager(
			log,
			[]core.RepositoryManager{
				NewDirManager(log),
			},
		)

		repo, err := manager.LoadRepository(repoUrl)

		s.NoError(err)
		s.Equal(repoUrl, repo.Url())
	})
}
