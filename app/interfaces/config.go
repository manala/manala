package interfaces

import (
	"github.com/caarlos0/log"
	"github.com/spf13/pflag"
)

type Config interface {
	log.Fielder
	Debug() bool
	BindDebugFlag(flag *pflag.Flag)
	Repository() string
	CacheDir() string
	BindCacheDirFlag(flag *pflag.Flag)
}
