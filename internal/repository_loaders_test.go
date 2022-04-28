package internal

import (
	"github.com/stretchr/testify/suite"
	"io"
	internalLog "manala/internal/log"
	"os"
	"path/filepath"
	"testing"
)

type RepositoryLoadersSuite struct{ suite.Suite }

func TestRepositoryLoadersSuite(t *testing.T) {
	suite.Run(t, new(RepositoryLoadersSuite))
}

var repositoryLoadersTestPath = filepath.Join("testdata", "repository_loaders")

func (s *RepositoryLoadersSuite) TestDir() {
	path := filepath.Join(repositoryLoadersTestPath, "dir")

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
		path := filepath.Join(path, "repository_not_found")

		repository, err := loader.LoadRepository([]string{path})

		var _err *NotFoundRepositoryError
		s.ErrorAs(err, &_err)

		s.ErrorAs(err, &internalError)
		s.Equal("repository not found", internalError.Message)
		s.Nil(repository)
	})

	s.Run("LoadRepository Wrong", func() {
		path := filepath.Join(path, "repository_wrong")

		repository, err := loader.LoadRepository([]string{path})

		s.ErrorAs(err, &internalError)
		s.Equal("wrong repository", internalError.Message)
		s.Nil(repository)
	})

	s.Run("LoadRepository", func() {
		path := filepath.Join(path, "repository")

		repository, err := loader.LoadRepository([]string{path})

		s.NoError(err)
		s.Equal(path, repository.Path())
		s.IsType(&RecipeRepositoryDirLoader{}, repository.recipeLoader)
	})
}

func (s *RepositoryLoadersSuite) TestGit() {
	path := filepath.Join(repositoryLoadersTestPath, "git")
	cacheDir := filepath.Join(path, "cache")

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
		repository, err := loader.LoadRepository([]string{"foo"})

		var _err *UnsupportedRepositoryError
		s.ErrorAs(err, &_err)

		s.ErrorAs(err, &internalError)
		s.Equal("unsupported repository", internalError.Message)
		s.Nil(repository)
	})

	s.Run("LoadRepository Error", func() {
		repository, err := loader.LoadRepository([]string{"https://github.com/octocat/Foo-Bar.git"})

		s.ErrorAs(err, &internalError)
		s.Equal("clone git repository", internalError.Message)
		s.Nil(repository)
	})

	s.Run("LoadRepository", func() {
		url := "https://github.com/octocat/Hello-World.git"

		repository, err := loader.LoadRepository([]string{url})

		s.NoError(err)
		s.Equal(url, repository.Path())
		s.DirExists(filepath.Join(cacheDir, "repositories", "3a1b4df2cfc5d2f2cc3259ef67cda77786d47f84"))
	})
}
