package yaml

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

/*******************/
/* Comment - Suite */
/*******************/

type CommentTestSuite struct{ suite.Suite }

func TestCommentTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(CommentTestSuite))
}

/*******************/
/* Comment - Tests */
/*******************/

func (s *CommentTestSuite) TestParse() {
	docTags := &DocTags{}
	ParseComment(`
	  # @foo bar
	  # @bar baz
      # qux
	`, docTags)
	s.Equal(
		&DocTags{
			{Name: "foo", Value: "bar"},
			{Name: "bar", Value: "baz\nqux"},
		},
		docTags,
	)
}
