package parsing

import (
	"testing"

	"github.com/manala/manala/internal/testing/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Assertion struct {
	Type   any
	Line   int
	Column int
	Err    errors.Assertion
}

func (a *Assertion) AssertError(t *testing.T, err error) {
	t.Helper()

	if a.Type != nil {
		require.IsType(t, a.Type, err)
	} else {
		require.IsType(t, (*Error)(nil), err)
	}

	e := err.(*Error)

	assert.Equal(t, a.Line, e.Line, "Line not equal")
	assert.Equal(t, a.Column, e.Column, "Column not equal")

	if a.Err != nil {
		a.Err.AssertError(t, e.Err)
	}
}

type FlattenAssertion Assertion

func (a *FlattenAssertion) AssertError(t *testing.T, err error) {
	t.Helper()

	require.IsType(t, (*Error)(nil), err)

	(*Assertion)(a).AssertError(t, err.(*Error).Flatten())
}
