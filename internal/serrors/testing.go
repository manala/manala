package serrors

import (
	"github.com/stretchr/testify/assert"
	"manala/internal/testing/heredoc"
	"testing"
)

type Assertion struct {
	Type      any
	Message   string
	Arguments []any
	Details   string
	Errors    []*Assertion
}

func Equal(t *testing.T, assertion *Assertion, err error) {
	if assertion.Type != nil {
		assert.IsType(t, assertion.Type, err)
	}

	assert.EqualError(t, err, assertion.Message)

	// Arguments
	if _err, ok := err.(ErrorArguments); ok {
		assert.Equal(t, assertion.Arguments, _err.ErrorArguments())
	} else if assertion.Arguments != nil {
		assert.Fail(t, "Error does not contains arguments")

	}

	// Details
	if _err, ok := err.(ErrorDetails); ok {
		heredoc.Equal(t, assertion.Details, _err.ErrorDetails(false))
	} else if assertion.Details != "" {
		assert.Fail(t, "Error does not contains details")
	}

	// Errors
	if _err, ok := err.(interface {
		Unwrap() []error
	}); ok {
		__errs := _err.Unwrap()
		if __errs == nil {
			if assertion.Errors != nil {
				assert.Fail(t, "Error contains nil errors")
			}
		} else {
			if assertion.Errors == nil {
				assert.Fail(t, "Error contains errors")
			} else {
				assert.Len(t, __errs, len(assertion.Errors), "Incorrect error's errors length")
				for i, _assert := range assertion.Errors {
					Equal(t, _assert, __errs[i])
				}
			}
		}
	} else if assertion.Errors != nil {
		assert.Fail(t, "Error does not contains errors")
	}
}
