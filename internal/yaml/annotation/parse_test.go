package annotation_test

import (
	"slices"
	"testing"

	"github.com/manala/manala/internal/testing/expectation"
	"github.com/manala/manala/internal/testing/heredoc"
	yamlannotation "github.com/manala/manala/internal/yaml/annotation"

	"github.com/stretchr/testify/suite"
)

type ParseSuite struct{ suite.Suite }

func TestParseSuite(t *testing.T) {
	suite.Run(t, new(ParseSuite))
}

func (s *ParseSuite) Test() {
	src := heredoc.Doc(`
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
	`)
	annotations, err := yamlannotation.Parse(src)

	s.Require().NotNil(annotations)
	s.Require().NoError(err)

	tests := []struct {
		test string
		name string
		body string
	}{
		{
			test: "Line",
			name: "line",
			body: "foo",
		},
		{
			test: "Multiline",
			name: "multiline",
			body: "foo\nbar",
		},
		{
			test: "Continuation",
			name: "continuation",
			body: "foo\nbar",
		},
		{
			test: "MultipleHashes",
			name: "multiple_hashes",
			body: "foo",
		},
		{
			test: "MixedHashes",
			name: "mixed_hashes",
			body: "foo",
		},
		{
			test: "DashedUnderScored",
			name: "dashed-under_scored",
			body: "foo\n@123 invalid name",
		},
		{
			test: "Indented",
			name: "indented",
			body: "foo",
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			i := slices.IndexFunc(annotations, func(a *yamlannotation.Annotation) bool {
				return a.Name.String() == test.name
			})
			s.Require().NotEqual(-1, i)

			annot := annotations[i]

			s.Require().NotNil(annot.Body)
			s.Equal(test.body, annot.Body.String())
		})
	}

	s.Len(annotations, len(tests))
}

func (s *ParseSuite) TestNoBody() {
	src := heredoc.Doc(`
		# @foo
	`)
	annotations, err := yamlannotation.Parse(src)

	s.Require().NotNil(annotations)
	s.Require().NoError(err)

	i := slices.IndexFunc(annotations, func(a *yamlannotation.Annotation) bool {
		return a.Name.String() == "foo"
	})
	s.Require().NotEqual(-1, i)

	annot := annotations[i]

	s.Require().NotNil(annot)
	s.Nil(annot.Body)
}

func (s *ParseSuite) TestErrors() {
	tests := []struct {
		test     string
		src      string
		expected expectation.ErrorExpectation
	}{
		{
			test: "Duplicate",
			src: heredoc.Doc(`
				# @foo bar
				# @foo baz
			`),
			expected: yamlannotation.ErrorExpectation{
				Position: [2]int{2, 3},
				Err:      expectation.ErrorMessage("duplicate @foo annotation"),
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			set, err := yamlannotation.Parse(test.src)

			s.Nil(set)

			expectation.ExpectError(s.T(), test.expected, err)
		})
	}
}
