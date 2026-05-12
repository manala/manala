package annotation

import (
	"testing"

	"github.com/manala/manala/internal/testing/expectation"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ErrorExpectation struct {
	Position [2]int
	Err      expectation.ErrorExpectation
}

func (a ErrorExpectation) Expect(t *testing.T, err error) {
	t.Helper()

	require.IsType(t, Error{}, err)
	e := err.(Error)

	line, column := e.Position()
	assert.Equal(t, a.Position[0], line, "line not equal")
	assert.Equal(t, a.Position[1], column, "column not equal")

	if a.Err != nil {
		a.Err.Expect(t, e.error)
	}
}
