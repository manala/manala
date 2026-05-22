package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type UnitExpectation struct {
	Source      string
	Destination string
}

func (a UnitExpectation) Expect(t *testing.T, unit Unit) {
	t.Helper()

	assert.Equal(t, a.Source, unit.Source, "source not equal")
	assert.Equal(t, a.Destination, unit.Destination, "destination not equal")
}

func ExpectUnit(t *testing.T, expectation UnitExpectation, unit Unit) {
	t.Helper()

	expectation.Expect(t, unit)
}

type UnitExpectations []UnitExpectation

func (a UnitExpectations) Expect(t *testing.T, units []Unit) {
	t.Helper()

	require.Len(t, units, len(a), "units count not equal")

	for i, expectation := range a {
		expectation.Expect(t, units[i])
	}
}

func ExpectUnits(t *testing.T, expectations UnitExpectations, units []Unit) {
	t.Helper()

	expectations.Expect(t, units)
}
