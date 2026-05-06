package log

import (
	"strings"

	"github.com/manala/manala/internal/output"
)

func (l *Log) Error(err error) {
	for _, e := range l.flattenError(err) {
		l.out.Print(l.error(e, 0))
	}
}

func (l *Log) error(err error, depth int) string {
	var b strings.Builder

	// Self-rendering bypasses the default pipeline (Attrs, Dump, Children).
	if r, ok := err.(Render); ok {
		return l.indent(r.Render(l.out.Profile), 3+depth*2) + "\n"
	}

	// Attrs
	var attrs [][2]any
	if a, ok := err.(Attrs); ok {
		attrs = a.Attrs()
	}

	b.WriteString(l.indent(
		l.log(Error, err.Error(), attrs), depth*2,
	) + "\n")

	// Dump
	if e, ok := err.(Dump); ok {
		if dump := e.Dump(); dump != "" {
			b.WriteString(l.indent(l.block(dump), 3+depth*2) + "\n")
		}
	}

	// Children
	if e, ok := err.(Err); ok {
		for _, child := range l.flattenError(e.Err()) {
			b.WriteString(l.error(child, depth+1))
		}
	}

	return b.String()
}

// flattenError recursively unwraps multi-errors (Unwrap() []error) into a flat slice of errors.
// Single-unwrap errors (Unwrap() error) are treated as leaves — not traversed — to preserve their message.
// An error whose Unwrap() returns empty is also treated as a leaf — it has children semantically but none in practice.
func (l *Log) flattenError(err error) []error {
	if err == nil {
		return nil
	}
	if multi, ok := err.(interface{ Unwrap() []error }); ok {
		if children := multi.Unwrap(); len(children) > 0 {
			var result []error
			for _, e := range children {
				result = append(result, l.flattenError(e)...)
			}
			return result
		}
	}
	return []error{err}
}

// Err is implemented by errors that expose a child error for structured display.
// The child may itself be a multi-error — it will be recursively flattened before rendering.
type Err interface {
	Err() error
}

// Dump is implemented by errors that carry raw content to display below their message.
// The content is rendered as a bordered block, indented under the error line.
type Dump interface {
	Dump() string
}

// Render is implemented by errors that produce their own full output.
// It bypasses the default pipeline entirely — no Attrs, no Dump, no Children.
type Render interface {
	Render(profile output.Profile) string
}
