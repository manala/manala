package filter

import (
	"slices"
)

// New creates a new filter.
func New(opts ...Option) *Filter {
	filter := &Filter{
		dotfiles: true,
	}

	// Options
	for _, opt := range opts {
		opt(filter)
	}

	return filter
}

// Filter is a filter.
type Filter struct {
	exclusions []string
	dotfiles   bool
}

// Excluded returns true if the path is excluded.
func (filter *Filter) Excluded(path string) bool {
	if !filter.dotfiles && path[0] == '.' {
		return true
	}

	return slices.Contains(filter.exclusions, path)
}

// Option is an option.
type Option func(filter *Filter)

func Without(paths ...string) Option {
	return func(filter *Filter) {
		filter.exclusions = append(filter.exclusions, paths...)
	}
}

func WithDotfiles(dotfiles bool) Option {
	return func(filter *Filter) {
		filter.dotfiles = dotfiles
	}
}
