package logger

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/discard"
	"io"
	"manala/config"
)

// Create a logger
func New(opts ...func(logger *logger)) Logger {
	logger := &logger{
		log: &log.Logger{
			Handler: cli.Default,
			Level:   log.DebugLevel,
		},
	}

	// Debug false by default
	WithDebug(false)(logger)

	// Options
	for _, opt := range opts {
		opt(logger)
	}

	return logger
}

func WithConfig(conf config.Config) func(logger *logger) {
	return func(logger *logger) {
		logger.debug = func() bool {
			return conf.Debug()
		}
	}
}

func WithDebug(debug bool) func(logger *logger) {
	return func(logger *logger) {
		logger.debug = func() bool {
			return debug
		}
	}
}

func WithWriter(writer io.Writer) func(logger *logger) {
	return func(logger *logger) {
		logger.log.Handler = cli.New(writer)
	}
}

func WithDiscardment() func(logger *logger) {
	return func(logger *logger) {
		logger.log.Handler = discard.New()
	}
}

type Logger interface {
	Debug(msg string, fields ...*field)
	Info(msg string, fields ...*field)
	WithField(key string, value interface{}) *field
	Error(msg string, errs ...error)
}

type field struct {
	key   string
	value interface{}
}

type logger struct {
	log   *log.Logger
	debug func() bool
}

func (logger *logger) Debug(msg string, fields ...*field) {
	if logger.debug() {
		entry := log.NewEntry(logger.log)

		// Fields
		for _, field := range fields {
			entry = entry.WithField(field.key, field.value)
		}

		entry.Debug(msg)
	}
}

func (logger *logger) Info(msg string, fields ...*field) {
	entry := log.NewEntry(logger.log)

	// Fields
	for _, field := range fields {
		entry = entry.WithField(field.key, field.value)
	}

	entry.Info(msg)
}

func (logger *logger) WithField(key string, value interface{}) *field {
	return &field{
		key:   key,
		value: value,
	}
}

func (logger *logger) Error(msg string, errs ...error) {
	entry := log.NewEntry(logger.log)

	// Errors
	for _, err := range errs {
		entry = entry.WithError(err)
	}

	entry.Error(msg)
}
