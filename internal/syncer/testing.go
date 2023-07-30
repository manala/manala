package syncer

import (
	"github.com/stretchr/testify/assert"
)

type UnitsAssert []struct {
	Source      string
	Destination string
}

func EqualUnits(s *assert.Assertions, assert *UnitsAssert, units []UnitInterface) {
	s.Len(units, len(*assert))
	for i, a := range *assert {
		s.Equal(a.Source, units[i].Source())
		s.Equal(a.Destination, units[i].Destination())
	}
}
