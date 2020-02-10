package doc

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

/*****************/
/* Parse - Suite */
/*****************/

type ParseTestSuite struct{ suite.Suite }

func TestParseTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(ParseTestSuite))
}

/*****************/
/* Parse - Tests */
/*****************/

func (s *ParseTestSuite) TestParseCommentTags() {
	tags := ParseCommentTags(`
	  # @foo bar
	  # @bar baz
      # qux
	`)
	s.Equal(
		TagList{
			tags: []*Tag{
				{Name: "foo", Value: "bar"},
				{Name: "bar", Value: "baz\nqux"},
			},
		},
		tags,
	)
}

/****************/
/* List - Suite */
/****************/

type ListTestSuite struct{ suite.Suite }

func TestListTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(ListTestSuite))
}

/****************/
/* List - Tests */
/****************/

func (s *ListTestSuite) TestListAll() {
	list := &TagList{}
	list.Add(&Tag{Name: "foo", Value: "bar"})
	list.Add(&Tag{Name: "bar", Value: "baz"})

	s.Equal(
		[]*Tag{
			{Name: "foo", Value: "bar"},
			{Name: "bar", Value: "baz"},
		},
		list.All(),
	)
}

func (s *ListTestSuite) TestListFilter() {
	list := &TagList{}
	list.Add(&Tag{Name: "foo", Value: "bar"})
	list.Add(&Tag{Name: "bar", Value: "baz"})
	list.Add(&Tag{Name: "foo", Value: "baz"})

	s.Equal(
		[]*Tag{
			{Name: "foo", Value: "bar"},
			{Name: "foo", Value: "baz"},
		},
		list.Filter("foo"),
	)
	s.Equal(
		[]*Tag{
			{Name: "bar", Value: "baz"},
		},
		list.Filter("bar"),
	)
	s.Equal(
		[]*Tag{},
		list.Filter("baz"),
	)
}
