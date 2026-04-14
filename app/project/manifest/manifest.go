package manifest

import (
	"slices"

	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/yaml/parser"
	"github.com/manala/manala/internal/yaml/validator"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
)

const filename = ".manala.yaml"

type Manifest struct {
	Recipe     string
	Repository string
	// Decoded separately.
	vars map[string]any
}

func New() *Manifest {
	return &Manifest{
		vars: map[string]any{},
	}
}

func (m *Manifest) Vars() map[string]any {
	return m.vars
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

	return nil
}

type manifestValidator struct{}

func (v manifestValidator) Struct(s any) error {
	m, ok := s.(Manifest)
	if !ok {
		return nil
	}

	var errs validator.FieldErrors

	// Recipe (required, max=100)
	if m.Recipe == "" {
		errs = append(errs, validator.NewFieldError("Recipe", "missing manala recipe property"))
	} else if len(m.Recipe) > 100 {
		errs = append(errs, validator.NewFieldError("Recipe", "too long manala recipe field (max=100)"))
	}

	// Repository (optional, max=256)
	if len(m.Repository) > 256 {
		errs = append(errs, validator.NewFieldError("Repository", "too long manala repository field (max=256)"))
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}
