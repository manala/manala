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
	s.Run("Dir", func() {
		cache := New("dir")

		dir, err := cache.Dir()

		s.NoError(err)
		s.Equal("dir", dir)
	})
	s.Run("DirSingle", func() {
		cache := New("dir")

		dir, err := cache.Dir("foo")

		s.NoError(err)
		s.Equal(filepath.Join("dir", "foo"), dir)
	})
	s.Run("DirMultiple", func() {
		cache := New("dir")

		dir, err := cache.Dir("foo", "bar")

		s.NoError(err)
		s.Equal(filepath.Join("dir", "foo", "bar"), dir)
	})
	s.Run("UserDir", func() {
		userDir, _ := os.UserCacheDir()

		cache := New("")

		dir, err := cache.Dir()

		s.NoError(err)
		s.Equal(userDir, dir)
	})
	s.Run("UserDirSingle", func() {
		userDir, _ := os.UserCacheDir()

		cache := New("")

		dir, err := cache.Dir("foo")

		s.NoError(err)
		s.Equal(filepath.Join(userDir, "foo"), dir)
	})
	s.Run("UserDirMultiple", func() {
		userDir, _ := os.UserCacheDir()

		cache := New("")

		dir, err := cache.Dir("foo", "bar")

		s.NoError(err)
		s.Equal(filepath.Join(userDir, "foo", "bar"), dir)
	})
	s.Run("UserDirSuffix", func() {
		userDir, _ := os.UserCacheDir()

		cache := New(
			"",
			WithUserDir("foo"),
		)

		dir, err := cache.Dir()

		s.NoError(err)
		s.Equal(filepath.Join(userDir, "foo"), dir)
	})
	s.Run("UserDirSuffixSingle", func() {
		userDir, _ := os.UserCacheDir()

		cache := New(
			"",
			WithUserDir("foo"),
		)

		dir, err := cache.Dir("bar")

		s.NoError(err)
		s.Equal(filepath.Join(userDir, "foo", "bar"), dir)
	})
	s.Run("UserDirSuffixMultiple", func() {
		userDir, _ := os.UserCacheDir()

		cache := New(
			"",
			WithUserDir("foo"),
		)

		dir, err := cache.Dir("bar", "baz")

		s.NoError(err)
		s.Equal(filepath.Join(userDir, "foo", "bar", "baz"), dir)
	})
}
