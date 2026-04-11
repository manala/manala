package parsing

import (
	"errors"

	"github.com/manala/manala/internal/serrors"
)

type Error struct {
	Err    error
	Line   int
	Column int
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func (e *Error) Unwrap() error {
	return e.Err
}

// Flatten walks the chain of nested parsing.Error, accumulating
// line/column positions, and returns a single flat Error with the
// resolved position and the root cause.
func (e *Error) Flatten() *Error {
	next, ok := errors.AsType[*Error](e.Unwrap())
	if !ok {
		return e
	}

	flat := next.Flatten()

	line := e.Line
	if flat.Line > 0 {
		line = max(1, line) + flat.Line - 1
	}

	column := e.Column
	if flat.Column > 0 {
		column = max(1, column) + flat.Column - 1
	}

	return &Error{
		Err:    flat.Err,
		Line:   line,
		Column: column,
	}
}

func ErrorTo(serr serrors.Error, err *Error, options Options) serrors.Error {
	err = err.Flatten()

	if err.Line == 0 && err.Column == 0 {
		return serr.WithErrors(err)
	}

	serr = serr.WithArguments("line", err.Line, "column", err.Column)

	if options.Src != "" {
		serr = serr.WithDetailsFunc((&Dumper{
			Err:   err,
			Src:   options.Src,
			Lexer: options.Src,
		}).Dump)
	} else {
		serr = serr.WithErrors(err)
	}

	return serr
}

type Options struct {
	Src   string
	Lexer string
}
