package errorstest

import (
	"testing"

	"github.com/manala/manala/internal/template/errors"
	"github.com/manala/manala/internal/testing/expectation"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Expectation struct {
	Position [2]int
	Err      expectation.ErrorExpectation
}

func (a Expectation) Expect(t *testing.T, err error) {
	t.Helper()

	require.IsType(t, errors.Error{}, err)
	e := err.(errors.Error)

	line, column := e.Position()
	assert.Equal(t, a.Position[0], line, "line not equal")
	assert.Equal(t, a.Position[1], column, "column not equal")

	if a.Err != nil {
		a.Err.Expect(t, e.Unwrap())
	}
}
