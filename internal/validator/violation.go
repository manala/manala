package validator

import (
	"errors"
	"manala/internal/path"
	"manala/internal/serrors"
)

type ViolationType int

const (
	REQUIRED ViolationType = iota + 1
	INVALID_TYPE
	ADDITIONAL_PROPERTY_NOT_ALLOWED
	STRING_GTE
	STRING_LTE
)

func NewViolation(message string) Violation {
	return Violation{
		Message: message,
	}
}

type Violation struct {
	Type              ViolationType
	Message           string
	StructuredMessage string
	Arguments         []any
	Path              path.Path
	Property          string
	Line              int
	Column            int
	DetailsFunc       func(ansi bool) string
}

func (violation Violation) Error() error {
	return errors.New(violation.Message)
}

func (violation Violation) StructuredError() serrors.Error {
	// Message
	message := violation.Message
	if violation.StructuredMessage != "" {
		message = violation.StructuredMessage
	}
	err := serrors.New(message)

	// Arguments
	err = err.WithArguments(violation.Arguments...)

	// Path
	if violation.Path != "" {
		err = err.WithArguments("path", violation.Path.String())
	}

	// Property
	if violation.Property != "" {
		err = err.WithArguments("property", violation.Property)
	}

	// Line
	if violation.Line != 0 {
		err = err.WithArguments("line", violation.Line)
	}

	// Column
	if violation.Column != 0 {
		err = err.WithArguments("column", violation.Column)
	}

	// Details
	if violation.DetailsFunc != nil {
		err = err.WithDetailsFunc(violation.DetailsFunc)
	}

	return err
}

type Violations []Violation

func (violations Violations) Errors() []error {
	errs := make([]error, len(violations))
	for i := range violations {
		errs[i] = violations[i].Error()
	}
	return errs
}

func (violations Violations) StructuredErrors() []error {
	errs := make([]error, len(violations))
	for i := range violations {
		errs[i] = violations[i].StructuredError()
	}
	return errs
}

func CompareViolations(a Violation, b Violation) int {
	if a.Column == b.Column && a.Line == b.Line {
		return 0
	}
	if a.Column >= b.Column {
		if a.Line < b.Line {
			return -1
		}
		return 1
	} else {
		if a.Line > b.Line {
			return 1
		}
		return -1
	}
}
