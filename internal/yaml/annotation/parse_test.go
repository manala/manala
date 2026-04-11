package annotation_test

import (
	"testing"

	"github.com/manala/manala/internal/parsing"
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/testing/errors"
	"github.com/manala/manala/internal/yaml/annotation"

	"github.com/stretchr/testify/suite"
)

type ParseSuite struct{ suite.Suite }

func TestParseSuite(t *testing.T) {
	suite.Run(t, new(ParseSuite))
}

func (s *ParseSuite) TestErrors() {
	tests := []struct {
		test     string
		src      string
		expected errors.Assertion
	}{
		{
			test: "Duplicate",
			src: `
# @foo bar
# @foo baz
`,
			expected: &parsing.ErrorAssertion{
				Line:   3,
				Column: 3,
				Err: &serrors.Assertion{
					Message: "duplicate annotation @foo",
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			_, err := annotation.Parse(test.src)

			errors.Equal(s.T(), test.expected, err)
		})
	}
}

func (s *ParseSuite) Test() {
	src := `
		# text before annotation
		# @line foo
		# @multiline foo
		# bar
		# @continuation foo
		#
		# bar
		## @multiple_hashes foo
		# ## # @mixed_hashes foo
		# @dashed-under_scored foo
		# @123 invalid name
  		 # @indented foo
		# @empty
	`
	set, err := annotation.Parse(src)
	s.Require().NoError(err)

	tests := []struct {
		test  string
		name  string
		value string
	}{
		{
			test:  "Line",
			name:  "line",
			value: "foo",
		},
		{
			test:  "Multiline",
			name:  "multiline",
			value: "foo\nbar",
		},
		{
			test:  "Continuation",
			name:  "continuation",
			value: "foo\nbar",
		},
		{
			test:  "MultipleHashes",
			name:  "multiple_hashes",
			value: "foo",
		},
		{
			test:  "MixedHashes",
			name:  "mixed_hashes",
			value: "foo",
		},
		{
			test:  "DashedUnderScored",
			name:  "dashed-under_scored",
			value: "foo\n@123 invalid name",
		},
		{
			test:  "Indented",
			name:  "indented",
			value: "foo",
		},
		{
			test:  "Empty",
			name:  "empty",
			value: "",
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			annot, ok := set.Lookup(test.name)
			s.Require().True(ok)
			s.Equal(test.value, annot.Value.String())
		})
	}

	s.Equal(len(tests), set.Len())
}
