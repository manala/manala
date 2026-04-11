package filter_test

import (
	"testing"

	"github.com/manala/manala/internal/filepath/filter"

	"github.com/stretchr/testify/suite"
)

type FilterSuite struct{ suite.Suite }

func TestFilterSuite(t *testing.T) {
	suite.Run(t, new(FilterSuite))
}

func (s *FilterSuite) Test() {
	s.Run("Default", func() {
		filter := filter.New()

		s.False(filter.Excluded("foo"))
		s.False(filter.Excluded(".bar"))
	})
	s.Run("Without", func() {
		filter := filter.New(
			filter.Without(
				"foo",
				"baz",
			),
		)

		s.True(filter.Excluded("foo"))
		s.False(filter.Excluded(".bar"))
		s.True(filter.Excluded("baz"))
	})
	s.Run("Dotfiles", func() {
		filter := filter.New(
			filter.WithDotfiles(false),
		)

		s.False(filter.Excluded("foo"))
		s.True(filter.Excluded(".baz"))
	})
}
