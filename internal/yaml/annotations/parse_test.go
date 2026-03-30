package annotations_test

import (
	"testing"

	"github.com/manala/manala/internal/yaml/annotations"

	"github.com/stretchr/testify/suite"
)

type ParseSuite struct{ suite.Suite }

func TestParseSuite(t *testing.T) {
	suite.Run(t, new(ParseSuite))
}

func (s *ParseSuite) Test() {
	src := `
		# text before annotations
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
	annots, err := annotations.Parse(src)
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

	for i, test := range tests {
		s.Run(test.test, func() {
			s.Equal(test.name, annots[i].Name())
			s.Equal(test.value, annots[i].Value())
		})
	}

	s.Require().Len(annots, len(tests))
}
