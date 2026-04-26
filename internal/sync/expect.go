package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type UnitsExpectation []struct {
	Source      string
	Destination string
}

func (a UnitsExpectation) Expect(t *testing.T, units []UnitInterface) {
	t.Helper()

	require.Len(t, units, len(a))

	for i, u := range a {
		assert.Equal(t, u.Source, units[i].Source())
		assert.Equal(t, u.Destination, units[i].Destination())
	}
}

func ExpectUnits(t *testing.T, expectation UnitsExpectation, units []UnitInterface) {
	t.Helper()

	expectation.Expect(t, units)
}
