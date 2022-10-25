package template

import (
	"errors"
	internalReport "manala/internal/report"
	"regexp"
	"strconv"
	textTemplate "text/template"
)

func NewParsingError(err error) *ParsingError {
	return &ParsingError{
		error: err,
	}
}

type ParsingError struct {
	error
}

func (err *ParsingError) Unwrap() error {
	return err.error
}

func (err *ParsingError) Error() string {
	if matches := parsingErrorRegex.FindStringSubmatch(err.error.Error()); matches != nil {
		return matches[3]
	}

	return err.error.Error()
}

func (err *ParsingError) Report(report *internalReport.Report) {
	if matches := parsingErrorRegex.FindStringSubmatch(err.error.Error()); matches != nil {
		// Line
		if line, _err := strconv.Atoi(matches[2]); _err == nil {
			report.Compose(
				internalReport.WithField("line", line),
			)
		}
	}
}

// 1 : template
// 2 : line
// 3 : message
var parsingErrorRegex = regexp.MustCompile(`template: (.*):(\d+): (.*)`)

func NewExecutionError(err error) *ExecutionError {
	return &ExecutionError{
		error: err,
	}
}

type ExecutionError struct {
	error
}

func (err *ExecutionError) Unwrap() error {
	return err.error
}

func (err *ExecutionError) Error() string {
	if matches := executionErrorRegex.FindStringSubmatch(err.error.Error()); matches != nil {
		return matches[6]
	}

	return err.error.Error()
}

func (err *ExecutionError) Report(report *internalReport.Report) {
	var _execError textTemplate.ExecError
	if errors.As(err.error, &_execError) {
		if matches := executionErrorRegex.FindStringSubmatch(err.error.Error()); matches != nil {
			// Line
			if line, _err := strconv.Atoi(matches[2]); _err == nil {
				report.Compose(
					internalReport.WithField("line", line),
				)
			}
			// Column
			if column, _err := strconv.Atoi(matches[3]); _err == nil {
				report.Compose(
					internalReport.WithField("column", column),
				)
			}
			// Context
			report.Compose(
				internalReport.WithField("context", matches[5]),
			)
		}
	}
}

// 1 : template
// 2 : line
// 3 : column
// 4 : name
// 5 : context
// 6 : message
var executionErrorRegex = regexp.MustCompile(`template: (.*):(\d+):(\d+): executing "(.*)" at <(.*)>: (.*)`)
