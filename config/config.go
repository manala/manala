package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

// Create a config
func New(opts ...func(config *config)) Config {
	// Viper
	v := viper.New()

	v.SetEnvPrefix("manala")
	v.AutomaticEnv()

	config := &config{
		viper: v,
	}

	// Options
	for _, opt := range opts {
		opt(config)
	}

	return config
}

func WithVersion(version string) func(config *config) {
	return func(config *config) {
		config.viper.Set("version", version)
	}
}

func WithMainRepository(repository string) func(config *config) {
	return func(config *config) {
		config.viper.Set("main_repository", repository)
		config.viper.SetDefault("repository", repository)
	}
}

func WithDebug(debug bool) func(config *config) {
	return func(config *config) {
		config.viper.SetDefault("debug", debug)
	}
}

func WithCacheDir(dir string) func(config *config) {
	return func(config *config) {
		config.viper.Set("cache_dir", dir)
	}
}

type Config interface {
	Version() string
	CacheDir() (string, error)
	BindCacheDirFlag(flag *pflag.Flag)
	MainRepository() string
	Repository() string
	BindRepositoryFlag(flag *pflag.Flag)
	Debug() bool
	BindDebugFlag(flag *pflag.Flag)
}

type config struct {
	viper *viper.Viper
}

func (config *config) Version() string {
	return config.viper.GetString("version")
}

func (config *config) CacheDir() (string, error) {
	if !config.viper.IsSet("cache_dir") {
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			return "", err
		}
		config.viper.Set("cache_dir", filepath.Join(cacheDir, "manala"))
	}

	return config.viper.GetString("cache_dir"), nil
}

func (config *config) BindCacheDirFlag(flag *pflag.Flag) {
	_ = config.viper.BindPFlag("cache_dir", flag)
}

func (config *config) MainRepository() string {
	return config.viper.GetString("main_repository")
}

func (config *config) Repository() string {
	return config.viper.GetString("repository")
}

func (config *config) BindRepositoryFlag(flag *pflag.Flag) {
	_ = config.viper.BindPFlag("repository", flag)
}

func (config *config) Debug() bool {
	return config.viper.GetBool("debug")
}

func (config *config) BindDebugFlag(flag *pflag.Flag) {
	_ = config.viper.BindPFlag("debug", flag)
}
