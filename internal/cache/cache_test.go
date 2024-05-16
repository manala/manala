package cache

import (
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
	"testing"
)

type CacheSuite struct{ suite.Suite }

func TestCacheSuite(t *testing.T) {
	suite.Run(t, new(CacheSuite))
}

func (s *CacheSuite) TestDir() {
	s.Run("Default", func() {
		cache := New("dir")

		dir, err := cache.Dir()

		s.NoError(err)
		s.Equal("dir", dir)
	})
	s.Run("WithDirSingle", func() {
		cache := New("dir").
			WithDir("foo")

		dir, err := cache.Dir()

		s.NoError(err)
		s.Equal(filepath.Join("dir", "foo"), dir)
	})
	s.Run("WithDirMultiple", func() {
		cache := New("dir").
			WithDir("foo").
			WithDir("bar")

		dir, err := cache.Dir()

		s.NoError(err)
		s.Equal(filepath.Join("dir", "foo", "bar"), dir)
	})
	s.Run("WithUserDirIgnored", func() {
		cache := New("dir").
			WithUserDir("foo")

		dir, err := cache.Dir()

		s.NoError(err)
		s.Equal("dir", dir)
	})
	s.Run("WithUserDir", func() {
		userDir, _ := os.UserCacheDir()

		cache := New("").
			WithUserDir("foo")

		dir, err := cache.Dir()

		s.NoError(err)
		s.Equal(filepath.Join(userDir, "foo"), dir)
	})
	s.Run("WithUserDirAndDir", func() {
		userDir, _ := os.UserCacheDir()

		cache := New("").
			WithUserDir("foo").
			WithDir("bar")

		dir, err := cache.Dir()

		s.NoError(err)
		s.Equal(filepath.Join(userDir, "foo", "bar"), dir)
	})
	s.Run("WithHashDir", func() {
		cache := New("dir").
			WithHashDir("foo")

		dir, err := cache.Dir()

		s.NoError(err)
		s.Equal(filepath.Join("dir", "0808f64e60d58979fcb676c96ec938270dea42445aeefcd3a4e6f8db"), dir)
	})
	s.Run("WithUserDirAndHashDir", func() {
		userDir, _ := os.UserCacheDir()

		cache := New("").
			WithUserDir("foo").
			WithHashDir("bar")

		dir, err := cache.Dir()

		s.NoError(err)
		s.Equal(filepath.Join(userDir, "foo", "07daf010de7f7f0d8d76a76eb8d1eb40182c8d1e7a3877a6686c9bf0"), dir)
	})
}
