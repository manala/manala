package yaml_test

import (
	"testing"

	"github.com/manala/manala/internal/yaml"

	"github.com/stretchr/testify/suite"
)

type TagsSuite struct{ suite.Suite }

func TestTagsSuite(t *testing.T) {
	suite.Run(t, new(TagsSuite))
}

func (s *TagsSuite) TestParseComment() {
	tags := &yaml.Tags{}
	yaml.ParseCommentTags(`
	  # @foo bar
	  # @bar baz
      # qux
	`, tags)
	s.Equal(
		&yaml.Tags{
			{Name: "foo", Value: "bar"},
			{Name: "bar", Value: "baz\nqux"},
		},
		tags,
	)
}

func (s *TagsSuite) TestFilter() {
	tags := &yaml.Tags{
		{Name: "foo", Value: "bar"},
		{Name: "bar", Value: "baz"},
		{Name: "foo", Value: "baz"},
	}

	s.Equal(
		&yaml.Tags{
			{Name: "foo", Value: "bar"},
			{Name: "foo", Value: "baz"},
		},
		tags.Filter("foo"),
	)
	s.Equal(
		&yaml.Tags{
			{Name: "bar", Value: "baz"},
		},
		tags.Filter("bar"),
	)
	s.Equal(
		&yaml.Tags{},
		tags.Filter("baz"),
	)
}
