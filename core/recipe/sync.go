package recipe

import (
	internalSyncer "manala/internal/syncer"
	"strings"
)

type sync []internalSyncer.UnitInterface

func (sync *sync) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var values []string
	if err := unmarshal(&values); err != nil {
		return err
	}

	for _, value := range values {
		unit := &syncUnit{}

		unit.source = value
		unit.destination = unit.source

		// Separate source / destination
		splits := strings.Split(unit.source, " ")
		if len(splits) > 1 {
			unit.source = splits[0]
			unit.destination = splits[1]
		}

		*sync = append(*sync, unit)
	}

	return nil
}

type syncUnit struct {
	source      string
	destination string
}

func (unit *syncUnit) Source() string {
	return unit.source
}

func (unit *syncUnit) Destination() string {
	return unit.destination
}
