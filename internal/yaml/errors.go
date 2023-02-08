package yaml

import (
	"github.com/goccy/go-yaml"
	yamlAst "github.com/goccy/go-yaml/ast"
	"github.com/muesli/termenv"
	internalReport "manala/internal/report"
	"regexp"
	"strconv"
)

func NewError(err error) *Error {
	newError := &Error{
		error: err,
	}

	color := !(termenv.EnvColorProfile() == termenv.Ascii)

	message := yaml.FormatError(err, color, true)
	if matches := errorRegex.FindStringSubmatch(message); matches != nil {
		// Message
		newError.message = matches[5]
		// Line
		if line, _err := strconv.Atoi(matches[3]); _err == nil {
			newError.line = line
		}
		// Column
		if column, _err := strconv.Atoi(matches[4]); _err == nil {
			newError.column = column
		}
		// Trace
		if matches[8] != "" {
			newError.trace = matches[8]
		}
	}

	return newError
}

type Error struct {
	error
	message string
	line    int
	column  int
	trace   string
}

func (err *Error) Unwrap() error {
	return err.error
}

func (err *Error) Error() string {
	if err.message != "" {
		return err.message
	}

	return err.error.Error()
}

func (err *Error) Report(report *internalReport.Report) {
	// Line
	if err.line != 0 {
		report.Compose(
			internalReport.WithField("line", err.line),
		)
	}
	// Column
	if err.column != 0 {
		report.Compose(
			internalReport.WithField("column", err.column),
		)
	}
	// Trace
	if err.trace != "" {
		report.Compose(
			internalReport.WithTrace(err.trace),
		)
	}
}

// 3: line (mutually optional with column)
// 4: column (mutually optional with line)
// 5: message
// 8: trace (optional)
var errorRegex = regexp.MustCompile(`(?s)^(\x1b\[91m)?(\[(\d+):(\d+)] )?([^\n]*)(\x1b\[0m)?(\n(.*))?$`)

func NewNodeError(message string, node yamlAst.Node) *NodeError {
	return &NodeError{
		message:  message,
		Reporter: NewReporter(node),
	}
}

type NodeError struct {
	message string
	*Reporter
}

func (err *NodeError) Error() string {
	return err.message
}
