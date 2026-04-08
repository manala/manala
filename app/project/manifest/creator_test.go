package manifest_test

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/manala/manala/app/project/manifest"
	"github.com/manala/manala/app/recipe"
	recipeManifest "github.com/manala/manala/app/recipe/manifest"
	"github.com/manala/manala/app/repository"
	"github.com/manala/manala/app/repository/getter"
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/testing/errors"
	"github.com/manala/manala/internal/testing/heredoc"

	"github.com/stretchr/testify/suite"
)

type CreatorSuite struct{ suite.Suite }

func TestCreatorSuite(t *testing.T) {
	suite.Run(t, new(CreatorSuite))
}

func (s *CreatorSuite) TestCreateErrors() {
	repositoryURL := filepath.FromSlash("testdata/CreatorSuite/TestCreateErrors/repository")
	recipeName := "recipe"

	repositoryLoader := repository.NewLoader(repository.WithLoaderHandlers(
		getter.NewFileLoaderHandler(slog.New(slog.DiscardHandler)),
	))
	repository, _ := repositoryLoader.Load(repositoryURL)

	recipeLoader := recipe.NewLoader(slog.New(slog.DiscardHandler), recipe.WithLoaderHandlers(
		recipeManifest.NewLoaderHandler(slog.New(slog.DiscardHandler)),
	))
	recipe, _ := recipeLoader.Load(repository, recipeName)

	s.Run("File", func() {
		projectDir := filepath.FromSlash("testdata/CreatorSuite/TestCreateErrors/File/project")

		creator := manifest.NewCreator()
		project, err := creator.Create(projectDir, recipe, nil)

		s.Nil(project)
		errors.Equal(s.T(), &serrors.Assertion{
			Message: "project is not a directory",
			Arguments: []any{
				"path", projectDir,
			},
		}, err)
	})
}

func (s *CreatorSuite) TestCreate() {
	repositoryURL := filepath.FromSlash("testdata/CreatorSuite/TestCreate/repository")
	recipeName := "recipe"

	projectDir := filepath.FromSlash("testdata/CreatorSuite/TestCreate/project")

	_ = os.RemoveAll(projectDir)

	repositoryLoader := repository.NewLoader(repository.WithLoaderHandlers(
		getter.NewFileLoaderHandler(slog.New(slog.DiscardHandler)),
	))
	repository, _ := repositoryLoader.Load(repositoryURL)

	recipeLoader := recipe.NewLoader(slog.New(slog.DiscardHandler), recipe.WithLoaderHandlers(
		recipeManifest.NewLoaderHandler(slog.New(slog.DiscardHandler)),
	))
	recipe, _ := recipeLoader.Load(repository, recipeName)

	vars := recipe.Vars()
	vars["string_float_int"] = "3.0"
	vars["string_asterisk"] = "*"

	creator := manifest.NewCreator()
	project, err := creator.Create(projectDir, recipe, vars)

	s.Require().NoError(err)
	s.NotNil(project)

	heredoc.EqualFile(s.T(), `
		manala:
		    recipe: recipe

		string_float_int: '3.0'
		string_float_int_value: '3.0'

		string_asterisk: '*'
		string_asterisk_value: '*'
	`, filepath.Join(projectDir, ".manala.yaml"))
}
