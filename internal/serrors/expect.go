package serrors

import (
	"testing"

	"github.com/manala/manala/internal/output"
	"github.com/manala/manala/internal/testing/expect"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Expectation struct {
	Message string
	Attrs   [][2]any
	Dump    string
	Errors  []expect.ErrorExpectation
}

func (a Expectation) Expect(t *testing.T, err error) {
	t.Helper()

	require.IsType(t, Error{}, err)
	e := err.(Error)

	require.EqualError(t, e, a.Message)

	// Attrs
	assert.Equal(t, a.Attrs, e.Attrs())

	// Dumper
	var dump string
	if dumper := e.Dumper(); dumper != nil {
		dump = dumper.Dump(output.Plain)
	}
	assert.Equal(t, a.Dump, dump)

	// Errors
	if _err, ok := err.(interface{ Unwrap() []error }); ok {
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
					_assert.Expect(t, _errs[i])
				}
			}
		}
	} else if a.Errors != nil {
		assert.Fail(t, "Error does not contain errors")
	}
}
