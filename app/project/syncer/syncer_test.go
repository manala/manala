package syncer

import (
	"github.com/stretchr/testify/suite"
	"manala/app/project"
	"manala/app/project/manifest"
	"manala/app/recipe"
	recipeManifest "manala/app/recipe/manifest"
	"manala/app/repository"
	"manala/app/repository/getter"
	"manala/internal/log"
	"manala/internal/testing/heredoc"
	"os"
	"path/filepath"
	"testing"
)

type Suite struct{ suite.Suite }

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestSync() {
	projectDir := filepath.FromSlash("testdata/Suite/TestSync/project")

	_ = os.Remove(filepath.Join(projectDir, "file.txt"))

	projectLoader := project.NewLoader(log.Discard,
		project.WithLoaderHandlers(
			manifest.NewLoaderHandler(log.Discard,
				repository.NewLoader(repository.WithLoaderHandlers(
					getter.NewFileLoaderHandler(log.Discard),
				)),
				recipe.NewLoader(log.Discard, recipe.WithLoaderHandlers(
					recipeManifest.NewLoaderHandler(log.Discard),
				)),
			),
		),
	)

	project, _ := projectLoader.Load(projectDir)

	syncer := New(log.Discard)
	err := syncer.Sync(project)

	s.NoError(err)
	heredoc.EqualFile(s.T(), `
		File
	`, filepath.Join(projectDir, "file.txt"))
}
