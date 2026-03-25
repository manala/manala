package serrors

import (
	"errors"
	htmlTemplate "html/template"
	"regexp"
	"strconv"
	textTemplate "text/template"
)

var (
	textExecErrorRegex       = regexp.MustCompile(`template: (?P<template>.*):(?P<line>\d+):(?P<column>\d+): executing "(?P<name>.*)" at <(?P<context>.*)>: (?P<message>.*)`)
	textParsingErrorRegex    = regexp.MustCompile(`template: (?P<template>.*):(?P<line>\d+): (?P<message>.*)`)
	htmlLineColumnErrorRegex = regexp.MustCompile(`html/template:(?P<template>.*):(?P<line>\d+):(?P<column>\d+): (?P<message>.*)`)
	htmlLineErrorRegex       = regexp.MustCompile(`html/template:(?P<template>.*):(?P<line>\d+): (?P<message>.*)`)
	htmlErrorRegex           = regexp.MustCompile(`html/template:(?:(?P<template>.*):)? (?P<message>.*)`)
)

func NewTemplate(err error) Error {
	var (
		arguments      []any
		_textExecError textTemplate.ExecError
		_htmlError     = &htmlTemplate.Error{}
	)

	message := err.Error()

	switch {
	// Text exec error
	case errors.As(err, &_textExecError):
		if matches := textExecErrorRegex.FindStringSubmatch(message); matches != nil {
			message = matches[textExecErrorRegex.SubexpIndex("message")]
			arguments = append(arguments, "context", matches[textExecErrorRegex.SubexpIndex("context")])
			// Template
			if template := matches[textExecErrorRegex.SubexpIndex("template")]; template != "" {
				arguments = append(arguments, "template", template)
			}
			// Line
			line, _ := strconv.Atoi(matches[textExecErrorRegex.SubexpIndex("line")])
			arguments = append(arguments, "line", line)
			// Column
			column, _ := strconv.Atoi(matches[textExecErrorRegex.SubexpIndex("column")])
			arguments = append(arguments, "column", column)
		}
	// Html error
	case errors.As(err, &_htmlError):
		if matches := htmlLineColumnErrorRegex.FindStringSubmatch(message); matches != nil {
			message = matches[htmlLineColumnErrorRegex.SubexpIndex("message")]
			// Template
			if template := matches[htmlLineColumnErrorRegex.SubexpIndex("template")]; template != "" {
				arguments = append(arguments, "template", template)
			}
			// Line
			line, _ := strconv.Atoi(matches[htmlLineColumnErrorRegex.SubexpIndex("line")])
			arguments = append(arguments, "line", line)
			// Column
			column, _ := strconv.Atoi(matches[htmlLineColumnErrorRegex.SubexpIndex("column")])
			arguments = append(arguments, "column", column)
		} else if matches := htmlLineErrorRegex.FindStringSubmatch(message); matches != nil {
			message = matches[htmlLineErrorRegex.SubexpIndex("message")]
			// Template
			if template := matches[htmlLineErrorRegex.SubexpIndex("template")]; template != "" {
				arguments = append(arguments, "template", template)
			}
			// Line
			line, _ := strconv.Atoi(matches[htmlLineErrorRegex.SubexpIndex("line")])
			arguments = append(arguments, "line", line)
		} else if matches := htmlErrorRegex.FindStringSubmatch(message); matches != nil {
			message = matches[htmlErrorRegex.SubexpIndex("message")]
			// Template
			if template := matches[htmlErrorRegex.SubexpIndex("template")]; template != "" {
				arguments = append(arguments, "template", template)
			}
		}
	default:
		// Text parsing error
		if matches := textParsingErrorRegex.FindStringSubmatch(message); matches != nil {
			message = matches[textParsingErrorRegex.SubexpIndex("message")]
			// Template
			if template := matches[textParsingErrorRegex.SubexpIndex("template")]; template != "" {
				arguments = append(arguments, "template", template)
			}
			// Line
			line, _ := strconv.Atoi(matches[textParsingErrorRegex.SubexpIndex("line")])
			arguments = append(arguments, "line", line)
		}
	}

	return New(message).
		WithArguments(arguments...)
}
