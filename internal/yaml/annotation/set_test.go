package annotation_test

import (
	"testing"

	jsonerrors "github.com/manala/manala/internal/json/errors"
	"github.com/manala/manala/internal/testing/expectation"
	yamlannotation "github.com/manala/manala/internal/yaml/annotation"

	"github.com/stretchr/testify/suite"
)

type SetSuite struct{ suite.Suite }

func TestSetSuite(t *testing.T) {
	suite.Run(t, new(SetSuite))
}

func (s *SetSuite) TestLookup() {
	src := `
		# @foo bar
		# @bar baz
	`
	set, err := yamlannotation.Parse(src)
	s.Require().NoError(err)
	s.Equal(2, set.Len())

	var annot *yamlannotation.Annotation
	var ok bool

	annot, ok = set.Lookup("foo")
	s.True(ok)
	s.Equal("foo", annot.Name.String())
	s.Equal("bar", annot.Value.String())

	annot, ok = set.Lookup("bar")
	s.True(ok)
	s.Equal("bar", annot.Name.String())
	s.Equal("baz", annot.Value.String())

	annot, ok = set.Lookup("baz")
	s.False(ok)
	s.Nil(annot)
}

func (s *SetSuite) TestJSONVar() {
	s.Run("Found", func() {
		src := `
# @foo {"bar": "baz"}
`
		set, err := yamlannotation.Parse(src)
		s.Require().NoError(err)

		var foo map[string]any
		s.Require().NoError(set.JSONVar(&foo, "foo"))

		s.Equal("baz", foo["bar"])
	})

	s.Run("NotFound", func() {
		src := `
# @foo {"bar": "baz"}
`
		set, err := yamlannotation.Parse(src)
		s.Require().NoError(err)

		var bar map[string]any
		s.Require().NoError(set.JSONVar(&bar, "bar"))

		s.Nil(bar)
	})

	s.Run("InvalidJSON", func() {
		src := `
# @foo bar
`
		set, err := yamlannotation.Parse(src)
		s.Require().NoError(err)

		var foo map[string]any
		err = set.JSONVar(&foo, "foo")

		expectation.ExpectError(s.T(), jsonerrors.Expectation{
			Position: [2]int{2, 8},
			Err:      expectation.ErrorMessage("invalid character 'b' looking for beginning of value"),
		}, err)
	})
}
