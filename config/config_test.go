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
	conf := New("foo", "bar")
	s.Equal("foo", conf.Version())
	s.Equal("bar", conf.MainRepository())
	s.Equal("bar", conf.Repository())
	cacheDir, _ := conf.CacheDir()
	s.NotEmpty(cacheDir)
	s.False(conf.Debug())
}

func (s *ConfigTestSuite) TestSet() {
	conf := New("foo", "bar")

	conf.SetCacheDir("baz")
	cacheDir, _ := conf.CacheDir()
	s.Equal("baz", cacheDir)

	conf.SetDebug(true)
	s.Equal(true, conf.Debug())
}

func (s *ConfigTestSuite) TestBind() {
	conf := New("foo", "bar")

	f := new(pflag.FlagSet)
	f.String("cache-dir", "", "")
	f.String("repository", "", "")
	f.Bool("debug", false, "")

	conf.BindCacheDirFlag(f.Lookup("cache-dir"))
	conf.BindRepositoryFlag(f.Lookup("repository"))
	conf.BindDebugFlag(f.Lookup("debug"))

	f.Set("cache-dir", "baz")
	cacheDir, _ := conf.CacheDir()
	s.Equal("baz", cacheDir)

	f.Set("repository", "qux")
	s.Equal("qux", conf.Repository())

	f.Set("debug", "1")
	s.True(conf.Debug())
}
