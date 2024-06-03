package sync_test

import (
	"manala/app/project"
	projectManifest "manala/app/project/manifest"
	"manala/app/project/sync"
	"manala/app/recipe"
	recipeManifest "manala/app/recipe/manifest"
	"manala/app/repository"
	repositoryGetter "manala/app/repository/getter"
	"manala/internal/log"
	"manala/internal/testing/heredoc"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type SyncerSuite struct{ suite.Suite }

func TestSyncerSuite(t *testing.T) {
	suite.Run(t, new(SyncerSuite))
}

func (s *SyncerSuite) TestSync() {
	projectDir := filepath.FromSlash("testdata/SyncerSuite/TestSync/project")

	_ = os.RemoveAll(filepath.Join(projectDir, "file.txt"))

	projectLoader := project.NewLoader(log.Discard,
		project.WithLoaderHandlers(
			projectManifest.NewLoaderHandler(log.Discard,
				repository.NewLoader(repository.WithLoaderHandlers(
					repositoryGetter.NewFileLoaderHandler(log.Discard),
				)),
				recipe.NewLoader(log.Discard, recipe.WithLoaderHandlers(
					recipeManifest.NewLoaderHandler(log.Discard),
				)),
			),
		),
	)

	project, err := projectLoader.Load(projectDir)
	s.Require().NoError(err)

	syncer := sync.NewSyncer(log.Discard)
	err = syncer.Sync(project)

	s.Require().NoError(err)
	heredoc.EqualFile(s.T(), `
		File
	`, filepath.Join(projectDir, "file.txt"))
}
