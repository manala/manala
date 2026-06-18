package serrortest

import (
	"testing"

	"github.com/manala/manala/internal/errors/serror"
	"github.com/manala/manala/internal/testing/expectation"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Expectation struct {
	Msg   string
	Attrs [][2]any
	Dump  string
	Err   expectation.ErrorExpectation
}

func (a Expectation) Expect(t *testing.T, err error) {
	t.Helper()

	require.IsType(t, serror.Error{}, err)
	e := err.(serror.Error)

	require.EqualError(t, e, a.Msg, "msg not equal")

	// Attrs
	assert.Equal(t, a.Attrs, e.Attrs(), "attrs not equal")

	// Dump
	assert.Equal(t, a.Dump, e.Dump(), "dump not equal")

	// Err
	if a.Err != nil {
		a.Err.Expect(t, e.Err())
	} else {
		assert.NoError(t, e.Err())
	}
}
