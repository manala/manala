package syncer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type UnitsAssertion []struct {
	Source      string
	Destination string
}

func EqualUnits(t *testing.T, assertion *UnitsAssertion, units []UnitInterface) {
	assert.Len(t, units, len(*assertion))
	for i, a := range *assertion {
		assert.Equal(t, a.Source, units[i].Source())
		assert.Equal(t, a.Destination, units[i].Destination())
	}
}
