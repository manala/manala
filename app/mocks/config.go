package mocks

import (
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/mock"
)

func MockConfig() *ConfigMock {
	return &ConfigMock{}
}

type ConfigMock struct {
	mock.Mock
}

func (conf *ConfigMock) Fields() map[string]interface{} {
	args := conf.Called()
	return args.Get(0).(map[string]interface{})
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
