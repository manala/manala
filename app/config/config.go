package config

type Config struct {
	Debug      bool   `mapstructure:"debug"`
	Repository string `mapstructure:"repository"`
	CacheDir   string `mapstructure:"cache-dir"`
}
