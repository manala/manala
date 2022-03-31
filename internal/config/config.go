package config

import (
	"github.com/apex/log"
	"github.com/spf13/viper"
)

func New() *Config {
	return &Config{Viper: viper.New()}
}

type Config struct {
	*viper.Viper
}

// Fields implements apex log Fielder
func (config *Config) Fields() log.Fields {
	return config.AllSettings()
}
