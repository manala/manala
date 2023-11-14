package ui

type CancelError struct{}

func (err *CancelError) Error() string { return "cancel execution" }
