package mocks

import (
	"github.com/caarlos0/log"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/mock"
)

func MockConfig() *ConfigMock {
	return &ConfigMock{}
}

type ConfigMock struct {
	mock.Mock
}

func (conf *ConfigMock) Fields() log.Fields {
	args := conf.Called()
	return args.Get(0).(log.Fields)
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

func (conf *ConfigMock) WebPort() int {
	args := conf.Called()
	return args.Int(0)
}

func (conf *ConfigMock) BindWebPortFlag(flag *pflag.Flag) {
	conf.Called(flag)
}
