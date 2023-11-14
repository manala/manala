package template

import (
	"errors"
	"manala/internal/serrors"
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

func NewError(err error) serrors.Error {
	message := err.Error()
	arguments := []any{}

	// Execution error
	var _execError textTemplate.ExecError
	if errors.As(err, &_execError) {
		if matches := executionErrorRegex.FindStringSubmatch(message); matches != nil {
			message = matches[6]
			arguments = append(arguments, "context", matches[5])
			// Line
			if line, _err := strconv.Atoi(matches[2]); _err == nil {
				arguments = append(arguments, "line", line)
			}
			// Column
			if column, _err := strconv.Atoi(matches[3]); _err == nil {
				arguments = append(arguments, "column", column)
			}
		}
	} else {
		// Parsing error
		if matches := parsingErrorRegex.FindStringSubmatch(message); matches != nil {
			message = matches[3]
			// Line
			if line, _err := strconv.Atoi(matches[2]); _err == nil {
				arguments = append(arguments, "line", line)
			}
		}
	}

	return serrors.New(message).
		WithArguments(arguments...)
}
