package config

import (
	"github.com/spf13/pflag"
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

type Config interface {
	Debug() bool
	BindDebugFlag(flag *pflag.Flag)
	Repository() string
	CacheDir() string
	BindCacheDirFlag(flag *pflag.Flag)
	Args() []any
}
