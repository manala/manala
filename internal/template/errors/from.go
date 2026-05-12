package errors

import (
	"errors"
	"regexp"
	"strconv"
	"text/template"
)

var (
	textExecErrorRegex    = regexp.MustCompile(`template: (?P<template>.*):(?P<line>\d+):(?P<column>\d+): executing "(?P<name>.*)" at <(?P<context>.*)>: (?P<message>.*)`)
	textParsingErrorRegex = regexp.MustCompile(`template: (?P<template>.*):(?P<line>\d+): (?P<message>.*)`)
)

// From try to convert a text/template error into an Error, extracting position from the error itself.
func From(err error, src string) error {
	message := err.Error()

	// Exec error
	if _, ok := errors.AsType[template.ExecError](err); ok {
		if matches := textExecErrorRegex.FindStringSubmatch(message); matches != nil {
			e := errors.New(matches[textExecErrorRegex.SubexpIndex("message")])

			// Line
			line, _ := strconv.Atoi(matches[textExecErrorRegex.SubexpIndex("line")])

			// Column (0-based byte offset)
			column, _ := strconv.Atoi(matches[textExecErrorRegex.SubexpIndex("column")])

			return At(e, src, line, column)
		}
	}

	// Parsing error
	if matches := textParsingErrorRegex.FindStringSubmatch(message); matches != nil {
		e := errors.New(matches[textParsingErrorRegex.SubexpIndex("message")])

		// Line
		line, _ := strconv.Atoi(matches[textParsingErrorRegex.SubexpIndex("line")])

		return Error{
			error:  e,
			line:   line,
			column: 0,
		}
	}

	return err
}
