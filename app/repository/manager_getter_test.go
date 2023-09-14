package repository

import (
	"github.com/stretchr/testify/suite"
	"io"
	"log/slog"
	"manala/internal/cache"
	"os"
	"path/filepath"
	"testing"
)

type GetterManagerSuite struct{ suite.Suite }

func TestGetterManagerSuite(t *testing.T) {
	suite.Run(t, new(GetterManagerSuite))
}

func (s *GetterManagerSuite) TestLoadRepository() {
	cacheDir := filepath.FromSlash("testdata/GetterManagerSuite/TestLoadRepository/cache")

	manager := NewGetterManager(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		cache.New(cacheDir),
	)

	s.Run("GitHttp", func() {
		_ = os.RemoveAll(cacheDir)

		url := "https://github.com/octocat/Hello-World.git"

		repository, err := manager.LoadRepository(url)

		s.NoError(err)

		s.Equal(url, repository.Url())
		s.DirExists(filepath.Join(cacheDir, "repositories", "0abe222eb5c9f101220b454d055bd5cbbb419722eddec6ad296f343f"))
	})

	s.Run("GitHttpForced", func() {
		_ = os.RemoveAll(cacheDir)

		url := "git::https://github.com/octocat/Hello-World.git"

		repository, err := manager.LoadRepository(url)

		s.NoError(err)

		s.Equal(url, repository.Url())
		s.DirExists(filepath.Join(cacheDir, "repositories", "fb19e43b4b91cd3024866d556a7d04f40cc66e18c5f7098b9a4af8d0"))
	})

	s.Run("HttpZip", func() {
		_ = os.RemoveAll(cacheDir)

		url := "https://github.com/octocat/Hello-World/archive/refs/heads/master.zip"

		repository, err := manager.LoadRepository(url)

		s.NoError(err)

		s.Equal(url, repository.Url())
		s.DirExists(filepath.Join(cacheDir, "repositories", "5b3a694d5b3f61795c95c7540461f94ffb950b81deeb3cea08e20e97"))
	})

	s.Run("HttpZipSubdirectory", func() {
		_ = os.RemoveAll(cacheDir)

		url := "https://github.com/octocat/Hello-World/archive/refs/heads/master.zip//Hello-World-master"

		repository, err := manager.LoadRepository(url)

		s.NoError(err)

		s.Equal(url, repository.Url())
		s.DirExists(filepath.Join(cacheDir, "repositories", "e8b20dc6a839cb37051c82b0acee5e6b63fd96894a9ce12dbf548ba8"))
	})

	s.Run("HttpZipSubdirectoryGlob", func() {
		_ = os.RemoveAll(cacheDir)

		url := "https://github.com/octocat/Hello-World/archive/refs/heads/master.zip//Hello-World-*"

		repository, err := manager.LoadRepository(url)

		s.NoError(err)

		s.Equal(url, repository.Url())
		s.DirExists(filepath.Join(cacheDir, "repositories", "5a7b4af36fbdca7290e1a95ffbb524e645707cf6695fa38c135771cf"))
	})

	s.Run("FileRelative", func() {
		_ = os.RemoveAll(cacheDir)

		url := filepath.FromSlash("testdata/GetterManagerSuite/TestLoadRepository/FileRelative/repository")

		repository, err := manager.LoadRepository(url)

		s.NoError(err)

		s.Equal(url, repository.Url())
	})

	s.Run("FileAbsolute", func() {
		_ = os.RemoveAll(cacheDir)

		url := filepath.FromSlash("testdata/GetterManagerSuite/TestLoadRepository/FileAbsolute/repository")
		url, _ = filepath.Abs(url)

		repository, err := manager.LoadRepository(url)

		s.NoError(err)

		s.Equal(url, repository.Url())
	})
}
