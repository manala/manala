package template

import (
	"errors"
	"manala/internal/errors/serrors"
	"regexp"
	"strconv"
	textTemplate "text/template"
)

// 1: template
// 2: line
// 3: column
// 4: name
// 5: context
// 6: message
var executionErrorRegex = regexp.MustCompile(`template: (.*):(\d+):(\d+): executing "(.*)" at <(.*)>: (.*)`)

// 1: template
// 2: line
// 3: message
var parsingErrorRegex = regexp.MustCompile(`template: (.*):(\d+): (.*)`)

func NewError(err error) *Error {
	_err := &Error{
		message:   err.Error(),
		Arguments: serrors.NewArguments(),
	}

	// Execution error
	var _execError textTemplate.ExecError
	if errors.As(err, &_execError) {
		if matches := executionErrorRegex.FindStringSubmatch(_err.message); matches != nil {
			// Message
			_err.message = matches[6]
			// Context
			_err.AppendArguments("context", matches[5])
			// Line
			if line, __err := strconv.Atoi(matches[2]); __err == nil {
				_err.AppendArguments("line", line)
			}
			// Column
			if column, __err := strconv.Atoi(matches[3]); __err == nil {
				_err.AppendArguments("column", column)
			}
		}
	} else {
		// Parsing error
		if matches := parsingErrorRegex.FindStringSubmatch(_err.message); matches != nil {
			_err.message = matches[3]
			// Line
			if line, __err := strconv.Atoi(matches[2]); __err == nil {
				_err.AppendArguments("line", line)
			}
		}
	}

	return _err
}

type Error struct {
	message string
	*serrors.Arguments
}

func (err *Error) Error() string {
	return err.message
}

func (err *Error) WithFile(file string) *Error {
	err.PrependArguments("file", file)
	return err
}
