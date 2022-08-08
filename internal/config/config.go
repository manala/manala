package config

import (
	"github.com/caarlos0/log"
	"github.com/spf13/viper"
)

func New() *Config {
	return &Config{Viper: viper.New()}
}

type Config struct {
	*viper.Viper
}

// Fields implements caarlos0 log Fielder
func (config *Config) Fields() log.Fields {
	return config.AllSettings()
}
