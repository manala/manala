package serrors

import (
	"testing"

	"github.com/manala/manala/internal/testing/heredoc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Assertion struct {
	Type      any
	Message   string
	Arguments []any
	Details   string
	Errors    []*Assertion
}

func Equal(t *testing.T, assertion *Assertion, err error) {
	t.Helper()

	if assertion.Type != nil {
		require.IsType(t, assertion.Type, err)
	}

	require.EqualError(t, err, assertion.Message)

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
		_errs := _err.Unwrap()
		if _errs == nil {
			if assertion.Errors != nil {
				assert.Fail(t, "Error contains nil errors")
			}
		} else {
			if assertion.Errors == nil {
				assert.Fail(t, "Error contains errors")
			} else {
				require.Len(t, _errs, len(assertion.Errors), "Incorrect error's errors length")

				for i, _assert := range assertion.Errors {
					Equal(t, _assert, _errs[i])
				}
			}
		}
	} else if assertion.Errors != nil {
		assert.Fail(t, "Error does not contains errors")
	}
}
