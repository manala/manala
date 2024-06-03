package serrors

import (
	"errors"
	htmlTemplate "html/template"
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
var textExecErrorRegex = regexp.MustCompile(`template: (.*):(\d+):(\d+): executing "(.*)" at <(.*)>: (.*)`)

// 1: template
// 2: line
// 3: message
var textParsingErrorRegex = regexp.MustCompile(`template: (.*):(\d+): (.*)`)

// 1: template
// 2: line
// 3: column
// 4: message
var htmlLineColumnErrorRegex = regexp.MustCompile(`html/template:(.*):(\d+):(\d+): (.*)`)

// 1: template
// 2: line
// 3: message
var htmlLineErrorRegex = regexp.MustCompile(`html/template:(.*):(\d+): (.*)`)

// 2: template (optional)
// 3: message
var htmlErrorRegex = regexp.MustCompile(`html/template:((.*):)? (.*)`)

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
	// Html error
	case errors.As(err, &_htmlError):
		if matches := htmlLineColumnErrorRegex.FindStringSubmatch(message); matches != nil {
			message = matches[4]
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
		} else if matches := htmlLineErrorRegex.FindStringSubmatch(message); matches != nil {
			message = matches[3]
			// Template
			if matches[1] != "" {
				arguments = append(arguments, "template", matches[1])
			}
			// Line
			if line, _err := strconv.Atoi(matches[2]); _err == nil {
				arguments = append(arguments, "line", line)
			}
		} else if matches := htmlErrorRegex.FindStringSubmatch(message); matches != nil {
			message = matches[3]
			// Template
			if matches[2] != "" {
				arguments = append(arguments, "template", matches[2])
			}
		}
	default:
		// Text parsing error
		if matches := textParsingErrorRegex.FindStringSubmatch(message); matches != nil {
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

	return New(message).
		WithArguments(arguments...)
}
