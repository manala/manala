package serrors

func Wrap(message string, err error) *WrapError {
	return &WrapError{
		message:   message,
		Arguments: NewArguments(),
		Details:   NewDetails(),
		err:       err,
	}
}

type WrapError struct {
	message string
	*Arguments
	*Details
	err error
}

func (err *WrapError) Error() string {
	return err.message
}

func (err *WrapError) Unwrap() error {
	return err.err
}

func (err *WrapError) WithArguments(arguments ...any) *WrapError {
	err.AppendArguments(arguments...)
	return err
}

func (err *WrapError) WithDetails(details string) *WrapError {
	err.SetDetails(details)
	return err
}
