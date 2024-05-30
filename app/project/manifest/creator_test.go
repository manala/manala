package manifest

import (
	"manala/app/recipe"
	"manala/app/recipe/manifest"
	"manala/app/repository"
	"manala/app/repository/getter"
	"manala/internal/log"
	"manala/internal/serrors"
	"manala/internal/testing/heredoc"
	"os"
	"path/filepath"
	"testing"

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
		getter.NewFileLoaderHandler(log.Discard),
	))
	repository, _ := repositoryLoader.Load(repositoryURL)

	recipeLoader := recipe.NewLoader(log.Discard, recipe.WithLoaderHandlers(
		manifest.NewLoaderHandler(log.Discard),
	))
	recipe, _ := recipeLoader.Load(repository, recipeName)

	s.Run("File", func() {
		projectDir := filepath.FromSlash("testdata/CreatorSuite/TestCreateErrors/File/project")

		creator := NewCreator()
		project, err := creator.Create(projectDir, recipe, nil)

		s.Nil(project)
		serrors.Equal(s.T(), &serrors.Assertion{
			Type:    serrors.Error{},
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
		getter.NewFileLoaderHandler(log.Discard),
	))
	repository, _ := repositoryLoader.Load(repositoryURL)

	recipeLoader := recipe.NewLoader(log.Discard, recipe.WithLoaderHandlers(
		manifest.NewLoaderHandler(log.Discard),
	))
	recipe, _ := recipeLoader.Load(repository, recipeName)

	vars := recipe.Vars()
	vars["string_float_int"] = "3.0"
	vars["string_asterisk"] = "*"

	creator := NewCreator()
	project, err := creator.Create(projectDir, recipe, vars)

	s.NotNil(project)
	s.NoError(err)

	heredoc.EqualFile(s.T(), `
		manala:
		    recipe: recipe

		string_float_int: '3.0'
		string_float_int_value: '3.0'

		string_asterisk: '*'
		string_asterisk_value: '*'
	`, filepath.Join(projectDir, ".manala.yaml"))
}
