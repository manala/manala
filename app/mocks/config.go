package mocks

import (
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/mock"
)

type ConfigMock struct {
	mock.Mock
}

func (conf *ConfigMock) Debug() bool {
	args := conf.Called()
	return args.Bool(0)
}

func (conf *ConfigMock) BindDebugFlag(flag *pflag.Flag) {
	conf.Called(flag)
}

func (conf *ConfigMock) Repository() string {
	args := conf.Called()
	return args.String(0)
}

func (conf *ConfigMock) CacheDir() string {
	args := conf.Called()
	return args.String(0)
}

func (conf *ConfigMock) BindCacheDirFlag(flag *pflag.Flag) {
	conf.Called(flag)
}

func (conf *ConfigMock) Args() []any {
	return []any{}
}
