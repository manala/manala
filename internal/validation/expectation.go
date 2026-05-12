package validation

import (
	"testing"

	"github.com/manala/manala/internal/testing/expectation"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ViolationExpectation struct {
	Position [2]int
	Err      expectation.ErrorExpectation
}

func (a ViolationExpectation) Expect(t *testing.T, err error) {
	t.Helper()

	require.IsType(t, &Violation{}, err)
	e := err.(*Violation)

	assert.Equal(t, a.Position[0], e.line, "line not equal")
	assert.Equal(t, a.Position[1], e.column, "column not equal")

	if a.Err != nil {
		a.Err.Expect(t, e.error)
	}
}
