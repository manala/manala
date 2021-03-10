package logger

import (
	apex "github.com/apex/log"
	apexHandlersCli "github.com/apex/log/handlers/cli"
	"io"
	"manala/config"
)

type Fields map[string]interface{}

type Logger struct {
	apex *apex.Logger
	conf *config.Config
}

func New(conf *config.Config) *Logger {
	// Apex
	apex := &apex.Logger{
		Handler: apexHandlersCli.Default,
		Level:   apex.DebugLevel,
	}

	return &Logger{
		apex: apex,
		conf: conf,
	}
}

func (log *Logger) SetOut(out io.Writer) {
	log.apex.Handler = apexHandlersCli.New(out)
}

func (log *Logger) Debug(msg string) {
	if log.conf.Debug() {
		log.apex.Debug(msg)
	}
}

func (log *Logger) DebugWithField(msg string, key string, value interface{}) {
	if log.conf.Debug() {
		log.apex.WithField(key, value).Debug(msg)
	}
}

func (log *Logger) DebugWithFields(msg string, fields Fields) {
	if log.conf.Debug() {
		f := apex.Fields{}
		for k, v := range fields {
			f[k] = v
		}
		log.apex.WithFields(f).Debug(msg)
	}
}

func (log *Logger) Info(msg string) {
	log.apex.Info(msg)
}

func (log *Logger) InfoWithField(msg string, key string, value interface{}) {
	log.apex.WithField(key, value).Info(msg)
}

func (log *Logger) InfoWithFields(msg string, fields Fields) {
	f := apex.Fields{}
	for k, v := range fields {
		f[k] = v
	}
	log.apex.WithFields(f).Info(msg)
}

func (log *Logger) Error(msg string) {
	log.apex.Error(msg)
}

func (log *Logger) ErrorWithError(msg string, err error) {
	log.apex.WithError(err).Error(msg)
}
