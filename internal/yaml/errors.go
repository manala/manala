package yaml

import (
	goYaml "github.com/goccy/go-yaml"
	goYamlAst "github.com/goccy/go-yaml/ast"
	"manala/internal/serrors"
	"regexp"
	"strconv"
)

// 3: line (mutually optional with column)
// 4: column (mutually optional with line)
// 5: message
// 8: details (optional)
var errorRegex = regexp.MustCompile(`(?s)^(\x1b\[91m)?(\[(\d+):(\d+)] )?([^\n]*)(\x1b\[0m)?(\n(.*))?$`)

func NewError(err error) serrors.Error {
	message := err.Error()
	arguments := []any{}

	str := goYaml.FormatError(err, false, false)
	if matches := errorRegex.FindStringSubmatch(str); matches != nil {
		// Message
		message = matches[5]
		// Line
		if line, __err := strconv.Atoi(matches[3]); __err == nil {
			arguments = append(arguments, "line", line)
		}
		// Column
		if column, __err := strconv.Atoi(matches[4]); __err == nil {
			arguments = append(arguments, "column", column)
		}
	}

	// Details
	detailsFunc := func(ansi bool) string {
		str := goYaml.FormatError(err, ansi, true)
		if matches := errorRegex.FindStringSubmatch(str); matches != nil {
			if matches[8] != "" {
				return matches[8]
			}
		}

		return ""
	}

	return serrors.New(message).
		WithArguments(arguments...).
		WithDetailsFunc(detailsFunc)

}

func NewNodeError(message string, node goYamlAst.Node) serrors.Error {
	err := serrors.New(message)

	if node == nil {
		return err
	}

	// Trace
	trace := NewNodeTrace(node)

	return err.
		WithArguments(
			"line", trace.Line,
			"column", trace.Column,
		).
		WithDetailsFunc(trace.DetailsFunc)
}
