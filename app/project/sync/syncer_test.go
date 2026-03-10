package sync_test

import (
	"github.com/manala/manala/app/project"
	projectManifest "github.com/manala/manala/app/project/manifest"
	"github.com/manala/manala/app/project/sync"
	"github.com/manala/manala/app/recipe"
	recipeManifest "github.com/manala/manala/app/recipe/manifest"
	"github.com/manala/manala/app/repository"
	repositoryGetter "github.com/manala/manala/app/repository/getter"
	"github.com/manala/manala/internal/log"
	"github.com/manala/manala/internal/testing/heredoc"
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
