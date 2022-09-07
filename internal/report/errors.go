package report

func NewError(err error) *Error {
	return &Error{
		error: err,
	}
}

type Error struct {
	error
	options []Option
}

func (err *Error) Unwrap() error {
	return err.error
}

func (err *Error) Report(report *Report) {
	report.Compose(err.options...)
}

func (err *Error) WithMessage(message string) *Error {
	err.options = append(err.options, WithMessage(message))

	return err
}

func (err *Error) WithField(key string, value interface{}) *Error {
	err.options = append(err.options, WithField(key, value))

	return err
}
