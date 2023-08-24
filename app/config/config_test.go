package config

import (
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ConfigSuite struct{ suite.Suite }

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}

func (s *ConfigSuite) Test() {
	conf := New()

	s.Run("Args", func() {
		s.Equal([]any{
			"repository", "https://github.com/manala/manala-recipes.git",
			"cache-dir", "",
			"debug", false,
		}, conf.Args())
	})

	s.Run("BindFlags", func() {
		flags := pflag.NewFlagSet("test", pflag.ContinueOnError)

		flags.Bool("foo", conf.Debug(), "test")
		conf.BindDebugFlag(flags.Lookup("foo"))

		flags.String("bar", conf.CacheDir(), "test")
		conf.BindCacheDirFlag(flags.Lookup("bar"))

		_ = flags.Set("foo", "1")
		_ = flags.Set("bar", "dir")

		s.Equal(true, conf.Debug())
		s.Equal("dir", conf.CacheDir())
	})
}
