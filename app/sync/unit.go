package sync

import (
	"strings"

	yamlerrors "github.com/manala/manala/internal/yaml/errors"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
)

type Unit struct {
	Source      string
	Destination string
}

func (u *Unit) UnmarshalYAML(node ast.Node) error {
	// Decode to string
	var value string
	if err := yaml.NodeToValue(node, &value); err != nil {
		return yamlerrors.From(err)
	}

	// Separate source / destination
	u.Source, u.Destination = value, value
	splits := strings.Split(u.Source, " ")
	if len(splits) > 1 {
		u.Source = splits[0]
		u.Destination = splits[1]
	}

	return nil
}
