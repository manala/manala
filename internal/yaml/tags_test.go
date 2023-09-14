package yaml

func (s *Suite) TestParseCommentTags() {
	tags := &Tags{}
	ParseCommentTags(`
	  # @foo bar
	  # @bar baz
      # qux
	`, tags)
	s.Equal(
		&Tags{
			{Name: "foo", Value: "bar"},
			{Name: "bar", Value: "baz\nqux"},
		},
		tags,
	)
}

func (s *Suite) TestTagsFilter() {
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
