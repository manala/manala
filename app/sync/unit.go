package sync

import (
	"strings"

	"github.com/manala/manala/internal/validation"
	yamlerrors "github.com/manala/manala/internal/yaml/errors"
	yamlvalidation "github.com/manala/manala/internal/yaml/validation"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
)

var unitValidator = validation.MustNewValidator(map[string]any{
	"type":      "string",
	"minLength": 1,
	"maxLength": 256,
})

type Unit struct {
	Source      string
	Destination string
}

func (u *Unit) UnmarshalYAML(node ast.Node) error {
	// Decode to string
	var data string
	if err := yaml.NodeToValue(node, &data); err != nil {
		return yamlerrors.From(err)
	}

	// Validate
	if violations, err := unitValidator.Validate(data, yamlvalidation.WithLocator(node)); err != nil {
		return err
	} else if violations != nil {
		return violations
	}

	// Separate source / destination
	u.Source, u.Destination = data, data
	splits := strings.Split(u.Source, " ")
	if len(splits) > 1 {
		u.Source = splits[0]
		u.Destination = splits[1]
	}

	return nil
}
