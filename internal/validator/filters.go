package validator

import "regexp"

type Filter struct {
	Path              string
	PathRegex         *regexp.Regexp
	Type              ViolationType
	Property          string
	Message           string
	StructuredMessage string
}

func (filter Filter) Format(violation *Violation) {
	// Path
	path := violation.Path.String()

	// Try to match on path
	if filter.Path != "" && filter.Path != path {
		return
	}
	// Try to match on path regex
	if filter.PathRegex != nil && !filter.PathRegex.MatchString(path) {
		return
	}
	// Try to match on type
	if filter.Type != 0 && filter.Type != violation.Type {
		return
	}
	// Try to match on property
	if filter.Property != "" && filter.Property != violation.Property {
		return
	}

	if filter.Message != "" {
		violation.Message = filter.Message
	}
	if filter.StructuredMessage != "" {
		violation.StructuredMessage = filter.StructuredMessage
	}
}

func WithFilters(filters Filters) Option {
	return func(validator *validator) {
		validator.formatters = append(validator.formatters, filters)
	}
}

type Filters []Filter

func (filters Filters) Format(violation *Violation) {
	// Apply filters
	for _, filter := range filters {
		filter.Format(violation)
	}
}
