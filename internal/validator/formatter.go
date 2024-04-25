package validator

type Formatter interface {
	Format(violation *Violation)
}

func WithFormatters(formatters ...Formatter) Option {
	return func(validator *validator) {
		validator.formatters = append(validator.formatters, formatters...)
	}
}
