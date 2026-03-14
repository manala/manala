package cmd

type CancelError struct{}

func (err *CancelError) Error() string { return "operation cancelled" }

type TerminalNotFoundError struct{}

func (err *TerminalNotFoundError) Error() string { return "terminal not found" }
