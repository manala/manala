package config

import (
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/suite"
	"testing"
)

/*********/
/* Suite */
/*********/

type ConfigTestSuite struct {
	suite.Suite
}

func TestConfigTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(ConfigTestSuite))
}

/*********/
/* Tests */
/*********/

func (s *ConfigTestSuite) Test() {
	conf := New()

	s.Equal("", conf.Version())
	cacheDir, _ := conf.CacheDir()
	s.NotEmpty(cacheDir)
	s.Equal("", conf.MainRepository())
	s.Equal("", conf.Repository())
	s.False(conf.Debug())
}

func (s *ConfigTestSuite) TestWithVersion() {
	s.Run("Dir", func() {
		version := "version"

		conf := New(WithVersion(version))

		s.Equal(version, conf.Version())
	})
}

func (s *ConfigTestSuite) TestWithMainRepository() {
	s.Run("Dir", func() {
		repository := "repository"

		conf := New(WithMainRepository(repository))

		s.Equal(repository, conf.MainRepository())
		s.Equal(repository, conf.Repository())
	})
}

func (s *ConfigTestSuite) TestWithDebug() {
	s.Run("True", func() {
		debug := true

		conf := New(WithDebug(debug))

		s.Equal(debug, conf.Debug())
	})

	s.Run("False", func() {
		debug := false

		conf := New(WithDebug(debug))

		s.Equal(debug, conf.Debug())
	})
}

func (s *ConfigTestSuite) TestWithCacheDir() {
	s.Run("Dir", func() {
		dir := "dir"

		conf := New(WithCacheDir(dir))

		cacheDir, _ := conf.CacheDir()
		s.Equal(dir, cacheDir)
	})
}

func (s *ConfigTestSuite) TestBind() {
	conf := New()
	f := new(pflag.FlagSet)

	s.Run("Cache dir", func() {
		dir := "dir"
		f.String("cache-dir", "", "")
		conf.BindCacheDirFlag(f.Lookup("cache-dir"))

		_ = f.Set("cache-dir", dir)

		cacheDir, _ := conf.CacheDir()
		s.Equal(dir, cacheDir)
	})

	s.Run("Repository", func() {
		repository := "repository"
		f.String("repository", "", "")
		conf.BindRepositoryFlag(f.Lookup("repository"))

		_ = f.Set("repository", repository)

		s.Equal(repository, conf.Repository())
	})

	s.Run("Debug", func() {
		f.Bool("debug", false, "")
		conf.BindDebugFlag(f.Lookup("debug"))

		_ = f.Set("debug", "1")

		s.Equal(true, conf.Debug())
	})
}
