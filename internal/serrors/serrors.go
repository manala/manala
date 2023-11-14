package serrors

func New(message string) Error {
	return Error{
		message: message,
	}
}

type Error struct {
	message     string
	arguments   []any
	detailsFunc func(ansi bool) string
	errors      []error
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

func (err Error) ErrorDetails(ansi bool) string {
	if err.detailsFunc == nil {
		return ""
	}
	return err.detailsFunc(ansi)
}

func (err Error) WithDetails(details string) Error {
	err.detailsFunc = func(ansi bool) string {
		return details
	}
	return err
}

func (err Error) WithDetailsFunc(detailsFunc func(ansi bool) string) Error {
	err.detailsFunc = detailsFunc
	return err
}

func (err Error) Unwrap() []error {
	return err.errors
}

func (err Error) WithErrors(errors ...error) Error {
	err.errors = append(err.errors, errors...)
	return err
}
