package sync_test

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/manala/manala/app/project"
	projectManifest "github.com/manala/manala/app/project/manifest"
	"github.com/manala/manala/app/project/sync"
	"github.com/manala/manala/app/recipe"
	recipeManifest "github.com/manala/manala/app/recipe/manifest"
	"github.com/manala/manala/app/repository"
	repositoryGetter "github.com/manala/manala/app/repository/getter"
	"github.com/manala/manala/app/template"
	"github.com/manala/manala/internal/testing/heredoc"

	"github.com/stretchr/testify/suite"
)

type SyncerSuite struct{ suite.Suite }

func TestSyncerSuite(t *testing.T) {
	suite.Run(t, new(SyncerSuite))
}

func (s *SyncerSuite) TestSync() {
	projectDir := filepath.FromSlash("testdata/SyncerSuite/TestSync/project")

	_ = os.RemoveAll(filepath.Join(projectDir, "file.txt"))

	projectLoader := project.NewLoader(slog.New(slog.DiscardHandler),
		project.WithLoaderHandlers(
			projectManifest.NewLoaderHandler(slog.New(slog.DiscardHandler),
				repository.NewLoader(repository.WithLoaderHandlers(
					repositoryGetter.NewFileLoaderHandler(slog.New(slog.DiscardHandler)),
				)),
				recipe.NewLoader(slog.New(slog.DiscardHandler), recipe.WithLoaderHandlers(
					recipeManifest.NewLoaderHandler(slog.New(slog.DiscardHandler)),
				)),
			),
		),
	)

	project, err := projectLoader.Load(projectDir)
	s.Require().NoError(err)

	syncer := sync.NewSyncer(slog.New(slog.DiscardHandler), template.NewEngine())
	err = syncer.Sync(project)

	s.Require().NoError(err)
	heredoc.EqualFile(s.T(), `
		File
	`, filepath.Join(projectDir, "file.txt"))
}
