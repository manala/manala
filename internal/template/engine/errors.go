package engine

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	textTemplate "text/template"

	"github.com/manala/manala/internal/parsing"
	"github.com/manala/manala/internal/serrors"
)

// ErrorAt creates a parsing.Error positioned at the given line and 0-based byte column.
// The byte column is converted to a 1-based rune column by iterating over runes (mirroring
// json/unmarshaler ErrorAt). A zero byte column means no column info — Column is left at 0.
func ErrorAt(err error, src string, line, byteColumn int) *parsing.Error {
	e := &parsing.Error{
		Err:  err,
		Line: line,
	}

	if src == "" || line <= 0 {
		return e
	}

	// Extract target line, then count runes up to byteColumn (0-based byte offset)
	lines := strings.SplitN(src, "\n", line+1)
	if line > len(lines) {
		return e
	}

	lineContent := lines[line-1]
	e.Column = 1
	for range lineContent[:min(byteColumn, len(lineContent))] {
		e.Column++
	}

	return e
}

var (
	textExecErrorRegex    = regexp.MustCompile(`template: (?P<template>.*):(?P<line>\d+):(?P<column>\d+): executing "(?P<name>.*)" at <(?P<context>.*)>: (?P<message>.*)`)
	textParsingErrorRegex = regexp.MustCompile(`template: (?P<template>.*):(?P<line>\d+): (?P<message>.*)`)
)

// ErrorFrom converts a text/template error into a parsing.Error, extracting position from the error itself.
func ErrorFrom(err error, src string) *parsing.Error {
	message := err.Error()

	// Exec error
	if _, ok := errors.AsType[textTemplate.ExecError](err); ok {
		if matches := textExecErrorRegex.FindStringSubmatch(message); matches != nil {
			e := serrors.New(matches[textExecErrorRegex.SubexpIndex("message")])

			// Context
			if context := matches[textExecErrorRegex.SubexpIndex("context")]; context != "" {
				e = e.WithArguments("context", context)
			}

			// Template
			if template := matches[textExecErrorRegex.SubexpIndex("template")]; template != "" {
				e = e.WithArguments("template", template)
			}

			// Line
			line, _ := strconv.Atoi(matches[textExecErrorRegex.SubexpIndex("line")])

			// Column (0-based byte offset)
			column, _ := strconv.Atoi(matches[textExecErrorRegex.SubexpIndex("column")])

			return ErrorAt(e, src, line, column)
		}
	}

	// Parsing error
	if matches := textParsingErrorRegex.FindStringSubmatch(message); matches != nil {
		e := serrors.New(matches[textParsingErrorRegex.SubexpIndex("message")])

		// Template
		if template := matches[textParsingErrorRegex.SubexpIndex("template")]; template != "" {
			e = e.WithArguments("template", template)
		}

		// Line
		line, _ := strconv.Atoi(matches[textParsingErrorRegex.SubexpIndex("line")])

		return &parsing.Error{
			Err:  e,
			Line: line,
		}
	}

	return &parsing.Error{Err: err}
}
