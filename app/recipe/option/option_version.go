package option

import (
	"github.com/Masterminds/semver/v3"
	"manala/internal/path"
	"manala/internal/schema"
	"manala/internal/serrors"
	"manala/internal/ui/components"
	"manala/internal/validator"
)

func NewVersionOption(option *option, fields map[string]any) (*VersionOption, error) {
	// Option
	versionOption := &VersionOption{
		option: option,
	}

	// Constraint
	if constraint, ok := fields["version_constraint"].(string); ok {
		var err error
		versionOption.Constraint, err = semver.NewConstraint(constraint)
		if err != nil {
			return nil, serrors.New("invalid version constraint").
				WithArguments("constraint", constraint)
		}
	}

	return versionOption, nil
}

func NewVersionOptionUiFormField(option *VersionOption, vars *map[string]any) (components.FormField, error) {
	// Field
	field, err := components.NewFormFieldText(
		option.Name(),
		option.Label(),
		option.Help(),
		path.NewAccessor(
			option.Path(),
			vars,
		),
		validator.New(
			validator.WithValidators(
				schema.NewValidator(option.Schema()),
				option,
			),
		),
	)
	if err != nil {
		return nil, err
	}

	return field, nil
}

type VersionOption struct {
	*option
	Constraint *semver.Constraints
}

func (option *VersionOption) Validate(value any) (validator.Violations, error) {
	versionString, ok := value.(string)
	if !ok {
		return validator.Violations{{
			Message: "version must be a string",
		}}, nil
	}
	version, err := semver.StrictNewVersion(versionString)
	if err != nil {
		return validator.Violations{{
			Message: err.Error(),
		}}, nil
	}
	if option.Constraint != nil && !option.Constraint.Check(version) {
		return validator.Violations{{
			Message: "version does not meet the constraint " + option.Constraint.String(),
		}}, nil
	}
	return nil, nil
}
