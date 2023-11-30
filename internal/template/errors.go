package template

import (
	"errors"
	"manala/internal/serrors"
	"regexp"
	"strconv"
	"text/template"
)

// 1: template
// 2: line
// 3: column
// 4: name
// 5: context
// 6: message
var execErrorRegex = regexp.MustCompile(`template: (.*):(\d+):(\d+): executing "(.*)" at <(.*)>: (.*)`)

// 1: template
// 2: line
// 3: message
var parsingErrorRegex = regexp.MustCompile(`template: (.*):(\d+): (.*)`)

func NewError(err error) serrors.Error {
	message := err.Error()
	var arguments []any

	// Exec error
	var _execError template.ExecError
	if errors.As(err, &_execError) {
		if matches := execErrorRegex.FindStringSubmatch(message); matches != nil {
			message = matches[6]
			arguments = append(arguments, "context", matches[5])
			// Template
			if matches[1] != "" {
				arguments = append(arguments, "template", matches[1])
			}
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
			// Template
			if matches[1] != "" {
				arguments = append(arguments, "template", matches[1])
			}
			// Line
			if line, _err := strconv.Atoi(matches[2]); _err == nil {
				arguments = append(arguments, "line", line)
			}
		}
	}

	return serrors.New(message).
		WithArguments(arguments...)
}
