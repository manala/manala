package validation

import (
	"errors"

	"github.com/santhosh-tekuri/jsonschema/v6"
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

func (v *Validator) Validate(value any, opts ...ValidateOption) (Violations, error) {
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
		return nil, nil
	}

	verr, ok := errors.AsType[*jsonschema.ValidationError](err)
	if !ok {
		return nil, err
	}

	errs := verr.BasicOutput().Errors
	if len(errs) == 0 {
		return nil, nil
	}

	var violations Violations
	for _, unit := range errs {
		location := unit.InstanceLocation
		line, column := cfg.locator.At(location)
		violations = append(violations, &Violation{
			error:  errors.New(unit.Error.String()),
			line:   line,
			column: column,
		})
	}

	return violations, nil
}

type validateConfig struct {
	locator Locator
}
type ValidateOption func(cfg *validateConfig)

func WithLocator(locator Locator) ValidateOption {
	return func(cfg *validateConfig) { cfg.locator = locator }
}
