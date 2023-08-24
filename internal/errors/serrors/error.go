package serrors

func New(message string) *Error {
	return &Error{
		message:   message,
		Arguments: NewArguments(),
		Details:   NewDetails(),
	}
}

type Error struct {
	message string
	*Arguments
	*Details
}

func (err *Error) Error() string {
	return err.message
}

func (err *Error) WithArguments(arguments ...any) *Error {
	err.AppendArguments(arguments...)
	return err
}

func (err *Error) WithDetails(details string) *Error {
	err.SetDetails(details)
	return err
}
