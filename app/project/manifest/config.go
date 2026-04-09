package manifest

import "github.com/manala/manala/internal/yaml/validator"

type config struct {
	Recipe     string
	Repository string
}

type configValidator struct{}

func (v configValidator) Struct(s any) error {
	cfg, ok := s.(config)
	if !ok {
		return nil
	}

	var errs validator.FieldErrors

	// Recipe (required, max=100)
	if cfg.Recipe == "" {
		errs = append(errs, validator.NewFieldError("Recipe", "missing manala recipe property"))
	} else if len(cfg.Recipe) > 100 {
		errs = append(errs, validator.NewFieldError("Recipe", "too long manala recipe field (max=100)"))
	}

	// Repository (optional, max=256)
	if len(cfg.Repository) > 256 {
		errs = append(errs, validator.NewFieldError("Repository", "too long manala repository field (max=256)"))
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}
