package validationtest

import (
	"testing"

	"github.com/manala/manala/internal/testing/expectation"
	"github.com/manala/manala/internal/validation"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ViolationExpectation struct {
	Location string
	Position [2]int
	Err      expectation.ErrorExpectation
}

func (a ViolationExpectation) Expect(t *testing.T, err error) {
	t.Helper()

	require.IsType(t, &validation.Violation{}, err)
	e := err.(*validation.Violation)

	assert.Equal(t, a.Location, e.Location(), "location not equal")

	line, column := e.Position()
	assert.Equal(t, a.Position[0], line, "line not equal")
	assert.Equal(t, a.Position[1], column, "column not equal")

	if a.Err != nil {
		a.Err.Expect(t, e.Unwrap())
	}
}
