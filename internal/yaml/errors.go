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
	return &Error{
		error: err,
	}
}

type Error struct {
	error
}

func (err *Error) Unwrap() error {
	return err.error
}

func (err *Error) Error() string {
	color := true
	if termenv.EnvColorProfile() == termenv.Ascii {
		color = false
	}

	message := yaml.FormatError(err.error, color, true)
	if matches := errorRegex.FindStringSubmatch(message); matches != nil {
		return matches[5]
	}

	return err.error.Error()
}

func (err *Error) Report(report *internalReport.Report) {
	color := true
	if termenv.EnvColorProfile() == termenv.Ascii {
		color = false
	}

	message := yaml.FormatError(err.error, color, true)
	if matches := errorRegex.FindStringSubmatch(message); matches != nil {
		// Line
		if line, _err := strconv.Atoi(matches[3]); _err == nil {
			report.Compose(
				internalReport.WithField("line", line),
			)
		}
		// Column
		if column, _err := strconv.Atoi(matches[4]); _err == nil {
			report.Compose(
				internalReport.WithField("column", column),
			)
		}
		// Trace
		if matches[8] != "" {
			report.Compose(
				internalReport.WithTrace(matches[8]),
			)
		}
	}
}

// 3 : line (mutually optional with column)
// 4 : column (mutually optional with line)
// 5 : message
// 8 : trace (optional)
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
