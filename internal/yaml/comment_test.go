package yaml

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type CommentSuite struct{ suite.Suite }

func TestCommentSuite(t *testing.T) {
	suite.Run(t, new(CommentSuite))
}

func (s *CommentSuite) TestParseCommentTags() {
	docTags := &Tags{}
	ParseCommentTags(`
	  # @foo bar
	  # @bar baz
      # qux
	`, docTags)
	s.Equal(
		&Tags{
			{Name: "foo", Value: "bar"},
			{Name: "bar", Value: "baz\nqux"},
		},
		docTags,
	)
}

func (s *CommentSuite) TestTagsFilter() {
	tags := &Tags{
		&Tag{Name: "foo", Value: "bar"},
		&Tag{Name: "bar", Value: "baz"},
		&Tag{Name: "foo", Value: "baz"},
	}

	s.Equal(
		&Tags{
			{Name: "foo", Value: "bar"},
			{Name: "foo", Value: "baz"},
		},
		tags.Filter("foo"),
	)
	s.Equal(
		&Tags{
			{Name: "bar", Value: "baz"},
		},
		tags.Filter("bar"),
	)
	s.Equal(
		&Tags{},
		tags.Filter("baz"),
	)
}
