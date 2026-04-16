package parsing

import (
	"testing"

	"github.com/manala/manala/internal/testing/expect"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ErrorExpectation struct {
	Line   int
	Column int
	Err    expect.ErrorExpectation
}

func (a ErrorExpectation) Expect(t *testing.T, err error) {
	t.Helper()

	require.IsType(t, (*Error)(nil), err)
	e := err.(*Error)

	assert.Equal(t, a.Line, e.Line, "Line not equal")
	assert.Equal(t, a.Column, e.Column, "Column not equal")

	if a.Err != nil {
		a.Err.Expect(t, e.Err)
	}
}

type FlattenErrorExpectation ErrorExpectation

func (a FlattenErrorExpectation) Expect(t *testing.T, err error) {
	t.Helper()

	require.IsType(t, (*Error)(nil), err)
	e := err.(*Error)

	ErrorExpectation(a).Expect(t, e.Flatten())
}
