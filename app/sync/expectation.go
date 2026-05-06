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

func (a UnitsExpectation) Expect(t *testing.T, units []Unit) {
	t.Helper()

	require.Len(t, units, len(a))

	for i, u := range a {
		assert.Equal(t, u.Source, units[i].Source, "source not equal")
		assert.Equal(t, u.Destination, units[i].Destination, "destination not equal")
	}
}

func ExpectUnits(t *testing.T, expectation UnitsExpectation, units []Unit) {
	t.Helper()

	expectation.Expect(t, units)
}
