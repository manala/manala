package validator

type Formatter interface {
	Format(violation *Violation)
}

func WithFormatters(formatters ...Formatter) Option {
	return func(validator *validator) {
		validator.formatters = append(validator.formatters, formatters...)
	}
}

func Formatters(fs ...Formatter) Formatter {
	return formatters(fs)
}

type formatters []Formatter

func (formatters formatters) Format(violation *Violation) {
	for _, formatter := range formatters {
		formatter.Format(violation)
	}
}
