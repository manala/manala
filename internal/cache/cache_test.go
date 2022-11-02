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
	s.Run("Dir Single", func() {
		cache := New("dir")

		dir, err := cache.Dir("foo")

		s.NoError(err)
		s.Equal(filepath.Join("dir", "foo"), dir)
	})
	s.Run("Dir Multiple", func() {
		cache := New("dir")

		dir, err := cache.Dir("foo", "bar")

		s.NoError(err)
		s.Equal(filepath.Join("dir", "foo", "bar"), dir)
	})
	s.Run("User Dir", func() {
		userDir, _ := os.UserCacheDir()

		cache := New("")

		dir, err := cache.Dir()

		s.NoError(err)
		s.Equal(userDir, dir)
	})
	s.Run("User Dir Single", func() {
		userDir, _ := os.UserCacheDir()

		cache := New("")

		dir, err := cache.Dir("foo")

		s.NoError(err)
		s.Equal(filepath.Join(userDir, "foo"), dir)
	})
	s.Run("User Dir Multiple", func() {
		userDir, _ := os.UserCacheDir()

		cache := New("")

		dir, err := cache.Dir("foo", "bar")

		s.NoError(err)
		s.Equal(filepath.Join(userDir, "foo", "bar"), dir)
	})
	s.Run("User Dir Suffix", func() {
		userDir, _ := os.UserCacheDir()

		cache := New(
			"",
			WithUserDir("foo"),
		)

		dir, err := cache.Dir()

		s.NoError(err)
		s.Equal(filepath.Join(userDir, "foo"), dir)
	})
	s.Run("User Dir Suffix Single", func() {
		userDir, _ := os.UserCacheDir()

		cache := New(
			"",
			WithUserDir("foo"),
		)

		dir, err := cache.Dir("bar")

		s.NoError(err)
		s.Equal(filepath.Join(userDir, "foo", "bar"), dir)
	})
	s.Run("User Dir Suffix Multiple", func() {
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
