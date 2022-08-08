package config

import (
	"github.com/caarlos0/log"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ConfigSuite struct{ suite.Suite }

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}

func (s *ConfigSuite) Test() {
	config := New()

	s.Run("Fields", func() {
		config.Set("key", "value")
		s.Equal(log.Fields{"key": "value"}, config.Fields())
	})

	s.Run("Set", func() {
		config.Set("set", "set")
		s.Equal("set", config.Get("set"))
	})

	s.Run("Get", func() {
		config.Set("get", "get")
		s.Equal("get", config.Get("get"))
	})

	s.Run("Get String", func() {
		config.Set("string", "string")
		s.Equal("string", config.GetString("string"))
	})

	s.Run("Get Bool", func() {
		config.Set("bool", true)
		s.Equal(true, config.GetBool("bool"))
	})

	s.Run("Bind Persistent Flags", func() {
		flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
		flags.String("flag", "flag", "test")
		_ = config.BindPFlags(flags)
		s.Equal("flag", config.GetString("flag"))
	})
}
