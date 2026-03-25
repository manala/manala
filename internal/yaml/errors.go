package yaml

import (
	"regexp"
	"strconv"

	"github.com/manala/manala/internal/serrors"

	goYaml "github.com/goccy/go-yaml"
	goYamlAst "github.com/goccy/go-yaml/ast"
)

var errorRegex = regexp.MustCompile(`(?s)^(?:\x1b\[91m)?(?:\[(?P<line>\d+):(?P<column>\d+)] )?(?P<message>[^\n]*)(?:\x1b\[0m)?(?:\n(?P<details>.*))?$`)

func NewError(err error) serrors.Error {
	var arguments []any
	message := err.Error()

	str := goYaml.FormatError(err, false, false)
	if matches := errorRegex.FindStringSubmatch(str); matches != nil {
		// Message
		message = matches[errorRegex.SubexpIndex("message")]
		// Line
		if line, _ := strconv.Atoi(matches[errorRegex.SubexpIndex("line")]); line != 0 {
			arguments = append(arguments, "line", line)
		}
		// Column
		if column, _ := strconv.Atoi(matches[errorRegex.SubexpIndex("column")]); column != 0 {
			arguments = append(arguments, "column", column)
		}
	}

	// Details
	detailsFunc := func(ansi bool) string {
		str := goYaml.FormatError(err, ansi, true)
		if matches := errorRegex.FindStringSubmatch(str); matches != nil {
			if details := matches[errorRegex.SubexpIndex("details")]; details != "" {
				return details
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
