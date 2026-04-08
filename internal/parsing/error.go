package parsing

import "github.com/manala/manala/internal/serrors"

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

func ErrorTo(serr serrors.Error, err *Error, options Options) serrors.Error {
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
