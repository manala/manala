package validator

import (
	"strings"
)

// FieldError implements goccy/go-yaml FieldError interface.
type FieldError struct {
	field   string
	message string
}

func NewFieldError(field, message string) FieldError {
	return FieldError{field: field, message: message}
}

func (e FieldError) StructField() string { return e.field }
func (e FieldError) Error() string       { return e.message }

// FieldErrors is a slice of FieldError.
// goccy/go-yaml checks via reflect that the error returned by StructValidator.Struct()
// is a slice to extract individual FieldError values.
type FieldErrors []FieldError

func (errs FieldErrors) Error() string {
	messages := make([]string, len(errs))
	for i, err := range errs {
		messages[i] = err.Error()
	}
	return strings.Join(messages, "; ")
}
