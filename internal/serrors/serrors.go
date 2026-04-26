package serrors

import (
	"io"
)

func New(msg string) Error {
	return Error{
		msg: msg,
	}
}

type Error struct {
	msg    string
	attrs  [][2]any
	dumper Dumper
	errs   []error
}

func (err Error) Error() string {
	return err.msg
}

func (err Error) Attrs() [][2]any {
	return err.attrs
}

func (err Error) With(args ...any) Error {
	a := args
	for len(a) >= 2 {
		if key, ok := a[0].(string); ok {
			err.attrs = append(err.attrs, [2]any{key, a[1]})
		}
		a = a[2:]
	}

	return err
}

func (err Error) Dumper() Dumper {
	return err.dumper
}

func (err Error) WithDumper(dumper Dumper) Error {
	err.dumper = dumper

	return err
}

func (err Error) WithDump(dump string) Error {
	err.dumper = StringDumper(dump)

	return err
}

func (err Error) Unwrap() []error {
	return err.errs
}

func (err Error) WithErrors(errs ...error) Error {
	for _, e := range errs {
		if e != nil {
			err.errs = append(err.errs, e)
		}
	}

	return err
}

type Attrs interface {
	Attrs() [][2]any
}

type Dumper = interface {
	Dump(w io.Writer)
}
