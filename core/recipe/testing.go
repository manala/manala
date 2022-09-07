package recipe

import (
	"github.com/stretchr/testify/suite"
	internalSyncer "manala/internal/syncer"
)

type syncAssert []struct {
	Source      string
	Destination string
}

func (assert *syncAssert) Equal(s *suite.Suite, sync []internalSyncer.UnitInterface) {
	s.Len(sync, len(*assert))
	for i, a := range *assert {
		s.Equal(a.Source, sync[i].Source())
		s.Equal(a.Destination, sync[i].Destination())
	}
}
