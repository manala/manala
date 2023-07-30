package interfaces

import (
	"github.com/spf13/pflag"
)

type Config interface {
	Debug() bool
	BindDebugFlag(flag *pflag.Flag)
	Repository() string
	CacheDir() string
	BindCacheDirFlag(flag *pflag.Flag)
	Args() []any
}
