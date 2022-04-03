package yaml

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

/******************/
/* DocTag - Suite */
/******************/

type DocTagTestSuite struct{ suite.Suite }

func TestDocTagTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(DocTagTestSuite))
}

/******************/
/* DocTag - Tests */
/******************/

func (s *DocTagTestSuite) TestFilter() {
	tags := &DocTags{
		&DocTag{Name: "foo", Value: "bar"},
		&DocTag{Name: "bar", Value: "baz"},
		&DocTag{Name: "foo", Value: "baz"},
	}

	s.Equal(
		[]*DocTag{
			{Name: "foo", Value: "bar"},
			{Name: "foo", Value: "baz"},
		},
		tags.Filter("foo"),
	)
	s.Equal(
		[]*DocTag{
			{Name: "bar", Value: "baz"},
		},
		tags.Filter("bar"),
	)
	s.Equal(
		[]*DocTag{},
		tags.Filter("baz"),
	)
}
