package template

import (
	"errors"
	internalErrors "manala/internal/errors"
	"regexp"
	"strconv"
	textTemplate "text/template"
)

var errorParsingRegex = regexp.MustCompile(`template: (?P<template>.*):(?P<line>\d+): (?P<message>.*)`)

func ParsingError(err error) *internalErrors.InternalError {
	_err := internalErrors.New("template error")

	if match := errorParsingRegex.FindStringSubmatch(err.Error()); match != nil {
		if line, err := strconv.Atoi(match[2]); err == nil {
			_ = _err.WithField("line", line)
		}
		_ = _err.WithField("message", match[3])
	} else {
		_ = _err.WithError(err)
	}

	return _err
}

func ParsingPathError(path string, err error) *internalErrors.InternalError {
	return ParsingError(err).
		WithField("path", path)
}

var errorExecutionRegex = regexp.MustCompile(`template: (?P<template>.*):(?P<line>\d+):(?P<column>\d+): executing "(?P<name>.*)" at <(?P<context>.*)>: (?P<message>.*)`)

func ExecutionError(err error) *internalErrors.InternalError {
	_err := internalErrors.New("template error")

	var execError textTemplate.ExecError
	if errors.As(err, &execError) {
		if match := errorExecutionRegex.FindStringSubmatch(execError.Error()); match != nil {
			if line, err := strconv.Atoi(match[2]); err == nil {
				_ = _err.WithField("line", line)
			}
			if column, err := strconv.Atoi(match[3]); err == nil {
				_ = _err.WithField("column", column)
			}
			_ = _err.
				WithField("context", match[5]).
				WithField("message", match[6])
		} else {
			_ = _err.WithError(execError.Err)
		}
	} else {
		_ = _err.WithError(err)
	}

	return _err
}

func ExecutionPathError(path string, err error) *internalErrors.InternalError {
	return ExecutionError(err).
		WithField("path", path)
}
