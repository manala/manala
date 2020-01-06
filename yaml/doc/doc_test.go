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
		[]Tag{
			{Name: "foo", Value: "bar"},
			{Name: "bar", Value: "baz\nqux"},
		},
		tags,
	)
}
