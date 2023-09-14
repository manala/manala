package config

import (
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (mock *Mock) Debug() bool {
	args := mock.Called()
	return args.Bool(0)
}

func (mock *Mock) BindDebugFlag(flag *pflag.Flag) {
	mock.Called(flag)
}

func (mock *Mock) Repository() string {
	args := mock.Called()
	return args.String(0)
}

func (mock *Mock) CacheDir() string {
	args := mock.Called()
	return args.String(0)
}

func (mock *Mock) BindCacheDirFlag(flag *pflag.Flag) {
	mock.Called(flag)
}

func (mock *Mock) Args() []any {
	return []any{}
}
