package config

import (
	"github.com/caarlos0/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"strings"
)

const (
	defaultDebug      = false
	defaultRepository = "https://github.com/manala/manala-recipes.git"
	defaultCacheDir   = ""
)

const (
	debugKey      = "debug"
	repositoryKey = "repository"
	cacheDirKey   = "cache-dir"
)

func New() *Config {
	v := viper.New()

	v.AutomaticEnv()
	v.SetEnvPrefix("MANALA")
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	v.SetDefault(debugKey, defaultDebug)
	v.SetDefault(repositoryKey, defaultRepository)
	v.SetDefault(cacheDirKey, defaultCacheDir)

	return &Config{
		viper: v,
	}
}

type Config struct {
	viper *viper.Viper
}

// Fields implements caarlos0 log Fielder
func (conf *Config) Fields() log.Fields {
	return conf.viper.AllSettings()
}

func (conf *Config) Debug() bool {
	return conf.viper.GetBool(debugKey)
}

func (conf *Config) BindDebugFlag(flag *pflag.Flag) {
	_ = conf.viper.BindPFlag(debugKey, flag)
}

func (conf *Config) Repository() string {
	return conf.viper.GetString(repositoryKey)
}

func (conf *Config) CacheDir() string {
	return conf.viper.GetString(cacheDirKey)
}

func (conf *Config) BindCacheDirFlag(flag *pflag.Flag) {
	_ = conf.viper.BindPFlag(cacheDirKey, flag)
}
