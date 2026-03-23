package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type UnitsAssertion []struct {
	Source      string
	Destination string
}

func EqualUnits(t *testing.T, assertion *UnitsAssertion, units []UnitInterface) {
	t.Helper()

	require.Len(t, units, len(*assertion))

	for i, a := range *assertion {
		assert.Equal(t, a.Source, units[i].Source())
		assert.Equal(t, a.Destination, units[i].Destination())
	}
}
