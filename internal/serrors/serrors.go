package serrors

import "slices"

func New(message string) Error {
	return Error{
		message: message,
	}
}

type Error struct {
	message   string
	arguments []any
	dumper    Dumper
	errors    []error
}

func (err Error) Error() string {
	return err.message
}

func (err Error) WithMessage(message string) Error {
	err.message = message

	return err
}

func (err Error) ErrorArguments() []any {
	return err.arguments
}

func (err Error) WithArguments(arguments ...any) Error {
	err.arguments = append(err.arguments, arguments...)

	return err
}

func (err Error) ErrorDump(ansi bool) string {
	if err.dumper == nil {
		return ""
	}

	return err.dumper.Dump(ansi)
}

func (err Error) WithDump(dump string) Error {
	err.dumper = StringDumper(dump)

	return err
}

func (err Error) WithDumper(dumper Dumper) Error {
	err.dumper = dumper

	return err
}

func (err Error) Unwrap() []error {
	return err.errors
}

func (err Error) WithErrors(errors ...error) Error {
	// Add only non nil errors
	err.errors = append(err.errors, slices.DeleteFunc(
		errors,
		func(err error) bool {
			return err == nil
		},
	)...)

	return err
}
