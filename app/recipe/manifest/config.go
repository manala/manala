package manifest

import (
	"github.com/manala/manala/app/recipe"
	"github.com/manala/manala/internal/yaml/validator"
)

type config struct {
	Description string
	Icon        string
	Template    string
	Partials    []string
	Sync        recipe.Sync
}

type configValidator struct{}

func (v configValidator) Struct(s any) error {
	cfg, ok := s.(config)
	if !ok {
		return nil
	}

	var errs validator.FieldErrors

	// Description (required, max=256)
	if cfg.Description == "" {
		errs = append(errs, validator.NewFieldError("Description", "missing manala description property"))
	} else if len(cfg.Description) > 256 {
		errs = append(errs, validator.NewFieldError("Description", "too long manala description field (max=256)"))
	}

	// Icon (optional, max=100)
	if len(cfg.Icon) > 100 {
		errs = append(errs, validator.NewFieldError("Icon", "too long manala icon field (max=100)"))
	}

	// Template (optional, max=100)
	if len(cfg.Template) > 100 {
		errs = append(errs, validator.NewFieldError("Template", "too long manala template field (max=100)"))
	}

	// Partials (optional, max=100 per entry)
	for _, partial := range cfg.Partials {
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
