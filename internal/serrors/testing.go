package serrors

import (
	"testing"

	"github.com/manala/manala/internal/testing/errors"
	"github.com/manala/manala/internal/testing/heredoc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Assertion struct {
	Type      any
	Message   string
	Arguments []any
	Details   string
	Errors    []errors.Assertion
}

func (a *Assertion) Assert(t *testing.T, err error) {
	t.Helper()

	if a.Type != nil {
		require.IsType(t, a.Type, err)
	} else {
		require.IsType(t, Error{}, err)
	}

	require.EqualError(t, err, a.Message)

	// Arguments
	if _err, ok := err.(ErrorArguments); ok {
		assert.Equal(t, a.Arguments, _err.ErrorArguments())
	} else if a.Arguments != nil {
		assert.Fail(t, "Error does not contains arguments")
	}

	// Details
	if _err, ok := err.(ErrorDetails); ok {
		heredoc.Equal(t, a.Details, _err.ErrorDetails(false))
	} else if a.Details != "" {
		assert.Fail(t, "Error does not contains details")
	}

	// Errors
	if _err, ok := err.(interface {
		Unwrap() []error
	}); ok {
		_errs := _err.Unwrap()
		if _errs == nil {
			if a.Errors != nil {
				assert.Fail(t, "Error contains nil errors")
			}
		} else {
			if a.Errors == nil {
				assert.Fail(t, "Error contains errors")
			} else {
				require.Len(t, _errs, len(a.Errors), "Incorrect error's errors length")

				for i, _assert := range a.Errors {
					_assert.Assert(t, _errs[i])
				}
			}
		}
	} else if a.Errors != nil {
		assert.Fail(t, "Error does not contains errors")
	}
}
