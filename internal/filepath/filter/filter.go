package filter

import (
	"slices"
)

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

type Filter struct {
	exclusions []string
	dotfiles   bool
}

func (filter *Filter) Excluded(path string) bool {
	if !filter.dotfiles && path[0] == '.' {
		return true
	}

	return slices.Contains(filter.exclusions, path)
}

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
