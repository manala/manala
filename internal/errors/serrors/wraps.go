package serrors

func Wraps(message string, errs ...error) *WrapsError {
	return &WrapsError{
		message:   message,
		Arguments: NewArguments(),
		Details:   NewDetails(),
		errs:      errs,
	}
}

type WrapsError struct {
	message string
	*Arguments
	*Details
	errs []error
}

func (err *WrapsError) Error() string {
	return err.message
}

func (err *WrapsError) Unwrap() []error {
	return err.errs
}

func (err *WrapsError) WithArguments(arguments ...any) *WrapsError {
	err.AppendArguments(arguments...)
	return err
}

func (err *WrapsError) WithDetails(details string) *WrapsError {
	err.SetDetails(details)
	return err
}
