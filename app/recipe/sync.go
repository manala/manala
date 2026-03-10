package recipe

import (
	"strings"

	"github.com/manala/manala/internal/sync"
)

type Sync []sync.UnitInterface

func (sync *Sync) UnmarshalYAML(unmarshal func(any) error) error {
	var values []string
	if err := unmarshal(&values); err != nil {
		return err
	}

	for _, value := range values {
		source, destination := value, value

		// Separate source / destination
		splits := strings.Split(source, " ")
		if len(splits) > 1 {
			source = splits[0]
			destination = splits[1]
		}

		*sync = append(*sync, NewSyncUnit(source, destination))
	}

	return nil
}

type SyncUnit struct {
	source      string
	destination string
}

func NewSyncUnit(source string, destination string) *SyncUnit {
	return &SyncUnit{
		source:      source,
		destination: destination,
	}
}

func (unit *SyncUnit) Source() string {
	return unit.source
}

func (unit *SyncUnit) Destination() string {
	return unit.destination
}
