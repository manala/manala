package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

type Config struct {
	version string
	viper   *viper.Viper
}

func New(version string, repository string) *Config {
	// Viper
	v := viper.New()

	v.SetEnvPrefix("manala")
	v.AutomaticEnv()

	v.SetDefault("repository", repository)
	v.SetDefault("debug", false)

	return &Config{
		version: version,
		viper:   v,
	}
}

func (conf *Config) Version() string {
	return conf.version
}

func (conf *Config) CacheDir() (string, error) {
	if !conf.viper.IsSet("cache_dir") {
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			return "", err
		}
		conf.viper.Set("cache_dir", filepath.Join(cacheDir, "manala"))
	}

	return conf.viper.GetString("cache_dir"), nil
}

func (conf *Config) SetCacheDir(dir string) {
	conf.viper.Set("cache_dir", dir)
}

func (conf *Config) BindCacheDirFlag(flag *pflag.Flag) {
	_ = conf.viper.BindPFlag("cache_dir", flag)
}

func (conf *Config) Debug() bool {
	return conf.viper.GetBool("debug")
}

func (conf *Config) SetDebug(debug bool) {
	conf.viper.Set("debug", debug)
}

func (conf *Config) BindDebugFlag(flag *pflag.Flag) {
	_ = conf.viper.BindPFlag("debug", flag)
}

func (conf *Config) Repository() string {
	return conf.viper.GetString("repository")
}
