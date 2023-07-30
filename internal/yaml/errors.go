package yaml

import (
	goYaml "github.com/goccy/go-yaml"
	goYamlAst "github.com/goccy/go-yaml/ast"
	goYamlPrinter "github.com/goccy/go-yaml/printer"
	"manala/internal/errors/serrors"
	"regexp"
	"strconv"
	"strings"
)

// 3: line (mutually optional with column)
// 4: column (mutually optional with line)
// 5: message
// 8: details (optional)
var errorRegex = regexp.MustCompile(`(?s)^(\x1b\[91m)?(\[(\d+):(\d+)] )?([^\n]*)(\x1b\[0m)?(\n(.*))?$`)

func NewError(err error) *Error {
	_err := &Error{
		message:   err.Error(),
		err:       err,
		Arguments: serrors.NewArguments(),
	}

	str := goYaml.FormatError(err, false, false)
	if matches := errorRegex.FindStringSubmatch(str); matches != nil {
		// Message
		_err.message = matches[5]
		// Line
		if line, __err := strconv.Atoi(matches[3]); __err == nil {
			_err.AppendArguments("line", line)
		}
		// Column
		if column, __err := strconv.Atoi(matches[4]); __err == nil {
			_err.AppendArguments("column", column)
		}
	}

	return _err
}

type Error struct {
	message string
	err     error
	*serrors.Arguments
}

func (err *Error) Error() string {
	return err.message
}

func (err *Error) ErrorDetails(ansi bool) string {
	str := goYaml.FormatError(err.err, ansi, true)
	if matches := errorRegex.FindStringSubmatch(str); matches != nil {
		if matches[8] != "" {
			return matches[8]
		}
	}

	return ""
}

func NewNodeError(message string, node goYamlAst.Node) *NodeError {
	_err := &NodeError{
		message:   message,
		node:      node,
		Arguments: serrors.NewArguments(),
	}

	if node != nil {
		token := node.GetToken()
		_err.AppendArguments(
			"line", token.Position.Line,
			"column", token.Position.Column,
		)
	}

	return _err
}

type NodeError struct {
	message string
	node    goYamlAst.Node
	*serrors.Arguments
}

func (err *NodeError) Error() string {
	return err.message
}

func (err *NodeError) ErrorDetails(ansi bool) string {
	if err.node == nil {
		return ""
	}

	var pp goYamlPrinter.Printer

	lines := strings.Split(
		strings.TrimRight(
			pp.PrintErrorToken(err.node.GetToken(), ansi),
			"\n",
		),
		"\n",
	)

	// Remove trailing spaces
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " ")
	}

	return strings.Join(lines, "\n") + "\n"
}

func (err *NodeError) WithArguments(arguments ...any) *NodeError {
	err.AppendArguments(arguments...)
	return err
}
