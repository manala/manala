package errors

import (
	"errors"
	"fmt"

	"github.com/goccy/go-yaml"
)

// From try to convert a github.com/goccy/go-yaml error into an Error, extracting token from the error itself.
func From(err error) error {
	// Type error
	if err, ok := errors.AsType[*yaml.TypeError](err); ok {
		return New(
			fmt.Errorf("field must be a %s", err.DstType),
			err.GetToken(),
		)
	}

	// Syntax error
	if err, ok := errors.AsType[*yaml.SyntaxError](err); ok {
		message := err.GetMessage()

		// Replace confusing tab message
		if message == "found character '\t' that cannot start any token" {
			message = "found a tab character where an indentation space is expected "
		}

		return New(
			errors.New(message),
			err.GetToken(),
		)
	}

	if err, ok := errors.AsType[yaml.Error](err); ok {
		return New(
			errors.New(err.GetMessage()),
			err.GetToken(),
		)
	}

	return err
}
