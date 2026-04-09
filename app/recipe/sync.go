package recipe

import (
	"strings"

	"github.com/manala/manala/internal/sync"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
)

type Sync []sync.UnitInterface

func (s *Sync) UnmarshalYAML(node ast.Node) error {
	seq, ok := node.(*ast.SequenceNode)
	if !ok {
		return &yaml.SyntaxError{
			Message: "sync field must be a sequence",
			Token:   node.GetToken(),
		}
	}

	for _, entry := range seq.Values {
		var unit SyncUnit

		switch entry.Type() {
		case ast.StringType:
			value := entry.(*ast.StringNode).Value

			// Validate
			if value == "" {
				return &yaml.SyntaxError{
					Message: "empty sync entry",
					Token:   node.GetToken(),
				}
			}
			if len(value) > 256 {
				return &yaml.SyntaxError{
					Message: "too long sync entry (max=256)",
					Token:   node.GetToken(),
				}
			}

			// Separate source / destination
			source, destination := value, value
			splits := strings.Split(source, " ")
			if len(splits) > 1 {
				source = splits[0]
				destination = splits[1]
			}
			unit = SyncUnit{source: source, destination: destination}
		default:
			return &yaml.SyntaxError{
				Message: "sync entry must be a string",
				Token:   node.GetToken(),
			}
		}

		*s = append(*s, &unit)
	}

	return nil
}

type SyncUnit struct {
	source      string
	destination string
}

func (unit *SyncUnit) Source() string {
	return unit.source
}

func (unit *SyncUnit) Destination() string {
	return unit.destination
}
