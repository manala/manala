package manifest

import (
	"slices"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/recipe"
	"github.com/manala/manala/internal/schema"
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/yaml/parser"
	"github.com/manala/manala/internal/yaml/validator"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
)

const filename = ".manala.yaml"

type Manifest struct {
	Description string
	Icon        string
	Template    string
	Partials    []string
	Sync        recipe.Sync
	// Decoded separately.
	vars map[string]any
	// Inferred from vars.
	schema  schema.Schema
	options []app.RecipeOption
}

func New() *Manifest {
	return &Manifest{
		Sync:   recipe.Sync{},
		vars:   map[string]any{},
		schema: schema.Schema{},
	}
}

func (m *Manifest) Vars() map[string]any {
	return m.vars
}

func (m *Manifest) Schema() schema.Schema {
	return m.schema
}

func (m *Manifest) Options() []app.RecipeOption {
	return m.options
}

func (m *Manifest) Unmarshal(content []byte) error {
	// Parse content to node
	node, err := parser.Parse(content)
	if err != nil {
		return err
	}

	// Partition manala & vars
	i := slices.IndexFunc(node.Values, func(node *ast.MappingValueNode) bool {
		return node.Key.String() == "manala"
	})
	if i == -1 {
		return parser.ErrorFrom(
			serrors.New("missing manala property"),
		)
	}

	manalaNode := node.Values[i].Value
	node.Values = slices.Concat(node.Values[:i], node.Values[i+1:])

	// Decode manala
	if err = yaml.NodeToValue(manalaNode, m,
		yaml.Validator(manifestValidator{}),
		yaml.DisallowUnknownField(),
	); err != nil {
		return parser.ErrorFrom(err)
	}

	// Decode vars
	if err = yaml.NodeToValue(node, &m.vars); err != nil {
		return parser.ErrorFrom(err)
	}

	// Infer schema & options
	inf := Inferrer{
		Schema:  &m.schema,
		Options: &m.options,
	}
	if err = inf.Infer(node); err != nil {
		return err
	}

	return nil
}

type manifestValidator struct{}

func (v manifestValidator) Struct(s any) error {
	m, ok := s.(Manifest)
	if !ok {
		return nil
	}

	var errs validator.FieldErrors

	// Description (required, max=256)
	if m.Description == "" {
		errs = append(errs, validator.NewFieldError("Description", "missing manala description property"))
	} else if len(m.Description) > 256 {
		errs = append(errs, validator.NewFieldError("Description", "too long manala description field (max=256)"))
	}

	// Icon (optional, max=100)
	if len(m.Icon) > 100 {
		errs = append(errs, validator.NewFieldError("Icon", "too long manala icon field (max=100)"))
	}

	// Template (optional, max=100)
	if len(m.Template) > 100 {
		errs = append(errs, validator.NewFieldError("Template", "too long manala template field (max=100)"))
	}

	// Partials (optional, max=100 per entry)
	for _, partial := range m.Partials {
		if partial == "" {
			errs = append(errs, validator.NewFieldError("Partials", "empty partials entry"))
			break
		}
		if len(partial) > 100 {
			errs = append(errs, validator.NewFieldError("Partials", "too long partials entry (max=100)"))
			break
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}
