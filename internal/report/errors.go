package report

func NewError(err error) *Error {
	return &Error{
		error:  err,
		errors: []error{},
	}
}

type Error struct {
	error
	errors  []error
	options []Option
}

func (err *Error) Unwrap() error {
	return err.error
}

func (err *Error) Report(report *Report) {
	report.Compose(err.options...)

	// Sub errors
	if len(err.errors) > 0 {
		for _, _err := range err.errors {
			report.Add(NewErrorReport(_err))
		}
	}
}

func (err *Error) Add(_err error) *Error {
	err.errors = append(err.errors, _err)

	return err
}

func (err *Error) WithMessage(message string) *Error {
	err.options = append(err.options, WithMessage(message))

	return err
}

func (err *Error) WithField(key string, value interface{}) *Error {
	err.options = append(err.options, WithField(key, value))

	return err
}
