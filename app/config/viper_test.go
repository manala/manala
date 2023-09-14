package config

import (
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ViperConfigSuite struct{ suite.Suite }

func TestViperConfigSuite(t *testing.T) {
	suite.Run(t, new(ViperConfigSuite))
}

func (s *ViperConfigSuite) Test() {
	config := NewViperConfig()

	s.Run("Args", func() {
		s.Equal([]any{
			"repository", "https://github.com/manala/manala-recipes.git",
			"cache-dir", "",
			"debug", false,
		}, config.Args())
	})

	s.Run("BindFlags", func() {
		flags := pflag.NewFlagSet("test", pflag.ContinueOnError)

		flags.Bool("foo", config.Debug(), "test")
		config.BindDebugFlag(flags.Lookup("foo"))

		flags.String("bar", config.CacheDir(), "test")
		config.BindCacheDirFlag(flags.Lookup("bar"))

		_ = flags.Set("foo", "1")
		_ = flags.Set("bar", "dir")

		s.Equal(true, config.Debug())
		s.Equal("dir", config.CacheDir())
	})
}
