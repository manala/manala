package filter

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type Suite struct{ suite.Suite }

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) Test() {
	s.Run("Default", func() {
		filter := New()

		s.False(filter.Excluded("foo"))
		s.False(filter.Excluded(".bar"))
	})
	s.Run("Without", func() {
		filter := New(
			Without(
				"foo",
				"baz",
			),
		)

		s.True(filter.Excluded("foo"))
		s.False(filter.Excluded(".bar"))
		s.True(filter.Excluded("baz"))
	})
	s.Run("Dotfiles", func() {
		filter := New(
			WithDotfiles(false),
		)

		s.False(filter.Excluded("foo"))
		s.True(filter.Excluded(".baz"))
	})
}
