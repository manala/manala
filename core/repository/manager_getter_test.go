package repository

import (
	"github.com/stretchr/testify/suite"
	"io"
	internalCache "manala/internal/cache"
	internalLog "manala/internal/log"
	internalTesting "manala/internal/testing"
	"os"
	"path/filepath"
	"testing"
)

type GetterManagerSuite struct{ suite.Suite }

func TestGetterManagerSuite(t *testing.T) {
	suite.Run(t, new(GetterManagerSuite))
}

func (s *GetterManagerSuite) TestLoadRepository() {
	cacheDir := internalTesting.DataPath(s, "cache")

	manager := NewGetterManager(
		internalLog.New(io.Discard),
		internalCache.New(cacheDir),
	)

	s.Run("Git Http", func() {
		_ = os.RemoveAll(cacheDir)

		repoUrl := "https://github.com/octocat/Hello-World.git"

		repo, err := manager.LoadRepository(repoUrl)

		s.NoError(err)
		s.Equal(repoUrl, repo.Url())
		s.DirExists(filepath.Join(cacheDir, "repositories", "0abe222eb5c9f101220b454d055bd5cbbb419722eddec6ad296f343f"))
	})

	s.Run("Git Http Forced", func() {
		_ = os.RemoveAll(cacheDir)

		repoUrl := "git::https://github.com/octocat/Hello-World.git"

		repo, err := manager.LoadRepository(repoUrl)

		s.NoError(err)
		s.Equal(repoUrl, repo.Url())
		s.DirExists(filepath.Join(cacheDir, "repositories", "fb19e43b4b91cd3024866d556a7d04f40cc66e18c5f7098b9a4af8d0"))
	})

	s.Run("Http Zip", func() {
		_ = os.RemoveAll(cacheDir)

		repoUrl := "https://github.com/octocat/Hello-World/archive/refs/heads/master.zip"

		repo, err := manager.LoadRepository(repoUrl)

		s.NoError(err)
		s.Equal(repoUrl, repo.Url())
		s.DirExists(filepath.Join(cacheDir, "repositories", "5b3a694d5b3f61795c95c7540461f94ffb950b81deeb3cea08e20e97"))
	})

	s.Run("Http Zip Subdirectory", func() {
		_ = os.RemoveAll(cacheDir)

		repoUrl := "https://github.com/octocat/Hello-World/archive/refs/heads/master.zip//Hello-World-master"

		repo, err := manager.LoadRepository(repoUrl)

		s.NoError(err)
		s.Equal(repoUrl, repo.Url())
		s.DirExists(filepath.Join(cacheDir, "repositories", "e8b20dc6a839cb37051c82b0acee5e6b63fd96894a9ce12dbf548ba8"))
	})

	s.Run("Http Zip Subdirectory Glob", func() {
		_ = os.RemoveAll(cacheDir)

		repoUrl := "https://github.com/octocat/Hello-World/archive/refs/heads/master.zip//Hello-World-*"

		repo, err := manager.LoadRepository(repoUrl)

		s.NoError(err)
		s.Equal(repoUrl, repo.Url())
		s.DirExists(filepath.Join(cacheDir, "repositories", "5a7b4af36fbdca7290e1a95ffbb524e645707cf6695fa38c135771cf"))
	})

	s.Run("File Relative", func() {
		_ = os.RemoveAll(cacheDir)

		repoUrl := internalTesting.DataPath(s, "repository")

		repo, err := manager.LoadRepository(repoUrl)

		s.NoError(err)
		s.Equal(repoUrl, repo.Url())
	})

	s.Run("File Absolute", func() {
		_ = os.RemoveAll(cacheDir)

		repoUrl, _ := filepath.Abs(internalTesting.DataPath(s, "repository"))

		repo, err := manager.LoadRepository(repoUrl)

		s.NoError(err)
		s.Equal(repoUrl, repo.Url())
	})
}
