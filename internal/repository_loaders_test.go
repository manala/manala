package internal

import (
	"github.com/stretchr/testify/suite"
	"io"
	internalLog "manala/internal/log"
	internalTesting "manala/internal/testing"
	"os"
	"path/filepath"
	"testing"
)

type RepositoryLoadersSuite struct{ suite.Suite }

func TestRepositoryLoadersSuite(t *testing.T) {
	suite.Run(t, new(RepositoryLoadersSuite))
}

func (s *RepositoryLoadersSuite) TestDir() {
	log := internalLog.New(io.Discard)
	loader := &RepositoryDirLoader{Log: log}

	s.Run("LoadRepository Empty", func() {
		repository, err := loader.LoadRepository([]string{})

		var _err *UnsupportedRepositoryError
		s.ErrorAs(err, &_err)

		s.ErrorAs(err, &internalError)
		s.Equal("unsupported repository", internalError.Message)
		s.Nil(repository)
	})

	s.Run("LoadRepository Not Found", func() {
		repository, err := loader.LoadRepository([]string{
			internalTesting.DataPath(s, "repository"),
		})

		var _err *NotFoundRepositoryError
		s.ErrorAs(err, &_err)

		s.ErrorAs(err, &internalError)
		s.Equal("repository not found", internalError.Message)
		s.Nil(repository)
	})

	s.Run("LoadRepository Wrong", func() {
		repository, err := loader.LoadRepository([]string{
			internalTesting.DataPath(s, "repository"),
		})

		s.ErrorAs(err, &internalError)
		s.Equal("wrong repository", internalError.Message)
		s.Nil(repository)
	})

	s.Run("LoadRepository", func() {
		repository, err := loader.LoadRepository([]string{
			internalTesting.DataPath(s, "repository"),
		})

		s.NoError(err)
		s.Equal(internalTesting.DataPath(s, "repository"), repository.Path())
		s.IsType(&RecipeRepositoryDirLoader{}, repository.recipeLoader)
	})
}

func (s *RepositoryLoadersSuite) TestGit() {
	cacheDir := internalTesting.DataPath(s, "cache")

	_ = os.RemoveAll(cacheDir)

	log := internalLog.New(io.Discard)
	loader := &RepositoryGitLoader{
		Log:      log,
		CacheDir: cacheDir,
	}

	s.Run("LoadRepository Empty", func() {
		repository, err := loader.LoadRepository([]string{})

		var _err *UnsupportedRepositoryError
		s.ErrorAs(err, &_err)

		s.ErrorAs(err, &internalError)
		s.Equal("unsupported repository", internalError.Message)
		s.Nil(repository)
	})

	s.Run("LoadRepository Unsupported", func() {
		repository, err := loader.LoadRepository([]string{
			"foo",
		})

		var _err *UnsupportedRepositoryError
		s.ErrorAs(err, &_err)

		s.ErrorAs(err, &internalError)
		s.Equal("unsupported repository", internalError.Message)
		s.Nil(repository)
	})

	s.Run("LoadRepository Error", func() {
		repository, err := loader.LoadRepository([]string{
			"https://github.com/octocat/Foo-Bar.git",
		})

		s.ErrorAs(err, &internalError)
		s.Equal("clone git repository", internalError.Message)
		s.Nil(repository)
	})

	s.Run("LoadRepository", func() {
		url := "https://github.com/octocat/Hello-World.git"

		repository, err := loader.LoadRepository([]string{
			url,
		})

		s.NoError(err)
		s.Equal(url, repository.Path())
		s.DirExists(filepath.Join(cacheDir, "repositories", "3a1b4df2cfc5d2f2cc3259ef67cda77786d47f84"))
	})
}
