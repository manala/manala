package validation

import (
	"cmp"
	"errors"
	"fmt"
	"slices"

	"github.com/go-openapi/jsonpointer"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/santhosh-tekuri/jsonschema/v6/kind"
)

type Validator struct {
	schema *jsonschema.Schema
}

func NewValidator(schema map[string]any) (*Validator, error) {
	compiler := jsonschema.NewCompiler()

	if err := compiler.AddResource("urn:manala", schema); err != nil {
		return nil, err
	}

	// Formats
	compiler.RegisterFormat(GitRepoFormat)
	compiler.RegisterFormat(FilePathFormat)
	compiler.RegisterFormat(DomainFormat)
	compiler.AssertFormat()

	compiled, err := compiler.Compile("urn:manala")
	if err != nil {
		return nil, err
	}

	return &Validator{
		schema: compiled,
	}, nil
}

func MustNewValidator(schema map[string]any) *Validator {
	validator, err := NewValidator(schema)
	if err != nil {
		panic(err)
	}
	return validator
}

func (v *Validator) Validate(value any, opts ...ValidateOption) error {
	// Config
	cfg := &validateConfig{
		locator: zeroLocator{},
	}

	// Options
	for _, opt := range opts {
		opt(cfg)
	}

	err := v.schema.Validate(value)
	if err == nil {
		return nil
	}

	verr, ok := errors.AsType[*jsonschema.ValidationError](err)
	if !ok {
		return err
	}

	errs := verr.BasicOutput().Errors
	if len(errs) == 0 {
		return nil
	}

	var violations Violations
	for _, unit := range errs {
		location := unit.InstanceLocation

		switch k := unit.Error.Kind.(type) {
		case *kind.AdditionalProperties:
			slices.Sort(k.Properties)
			for _, property := range k.Properties {
				propertyLocation := location + "/" + jsonpointer.Escape(property)
				line, column := cfg.locator.PropertyAt(propertyLocation)
				violations = append(violations, &Violation{
					error:    fmt.Errorf("additional property '%s' not allowed", property),
					location: propertyLocation,
					line:     line,
					column:   column,
				})
			}
		default:
			line, column := cfg.locator.ValueAt(location)
			violations = append(violations, &Violation{
				error:    errors.New(unit.Error.String()),
				location: location,
				line:     line,
				column:   column,
			})
		}
	}

	// Keep deterministic violations order
	slices.SortFunc(violations, func(a, b *Violation) int {
		// By location
		if c := cmp.Compare(a.location, b.location); c != 0 {
			return c
		}
		// By error message
		return cmp.Compare(a.error.Error(), b.error.Error())
	})

	return violations
}

type validateConfig struct {
	locator Locator
}
type ValidateOption func(cfg *validateConfig)

func WithLocator(locator Locator) ValidateOption {
	return func(cfg *validateConfig) { cfg.locator = locator }
}
