package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"strings"
)

func NewViperConfig() *ViperConfig {
	config := &ViperConfig{
		viper: viper.New(),
	}

	config.viper.AutomaticEnv()
	config.viper.SetEnvPrefix("MANALA")
	config.viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	config.viper.SetDefault(debugKey, defaultDebug)
	config.viper.SetDefault(repositoryKey, defaultRepository)
	config.viper.SetDefault(cacheDirKey, defaultCacheDir)

	return config
}

type ViperConfig struct {
	viper *viper.Viper
}

func (config *ViperConfig) Debug() bool {
	return config.viper.GetBool(debugKey)
}

func (config *ViperConfig) BindDebugFlag(flag *pflag.Flag) {
	_ = config.viper.BindPFlag(debugKey, flag)
}

func (config *ViperConfig) Repository() string {
	return config.viper.GetString(repositoryKey)
}

func (config *ViperConfig) CacheDir() string {
	return config.viper.GetString(cacheDirKey)
}

func (config *ViperConfig) BindCacheDirFlag(flag *pflag.Flag) {
	_ = config.viper.BindPFlag(cacheDirKey, flag)
}

func (config *ViperConfig) Args() []any {
	return []any{
		"repository", config.Repository(),
		"cache-dir", config.CacheDir(),
		"debug", config.Debug(),
	}
}
