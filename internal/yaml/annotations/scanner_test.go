package annotations_test

import (
	"testing"

	"github.com/manala/manala/internal/yaml/annotations"

	"github.com/stretchr/testify/suite"
)

type ScannerSuite struct{ suite.Suite }

func TestScannerSuite(t *testing.T) {
	suite.Run(t, new(ScannerSuite))
}

func (s *ScannerSuite) Test() {
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
	scanner := annotations.NewScanner(src)

	tests := []struct {
		test     string
		expected annotations.Token
	}{
		{
			test:     "Text",
			expected: annotations.Token{Kind: annotations.TokenText, Value: "text before annotations", Line: 2, Column: 5},
		},
		{
			test:     "Line",
			expected: annotations.Token{Kind: annotations.TokenName, Value: "line", Line: 3, Column: 5},
		},
		{
			test:     "LineValue",
			expected: annotations.Token{Kind: annotations.TokenText, Value: "foo", Line: 3, Column: 11},
		},
		{
			test:     "Multiline",
			expected: annotations.Token{Kind: annotations.TokenName, Value: "multiline", Line: 4, Column: 5},
		},
		{
			test:     "MultilineValue1",
			expected: annotations.Token{Kind: annotations.TokenText, Value: "foo", Line: 4, Column: 16},
		},
		{
			test:     "MultilineValue2",
			expected: annotations.Token{Kind: annotations.TokenText, Value: "bar", Line: 5, Column: 5},
		},
		{
			test:     "Continuation",
			expected: annotations.Token{Kind: annotations.TokenName, Value: "continuation", Line: 6, Column: 5},
		},
		{
			test:     "ContinuationValue1",
			expected: annotations.Token{Kind: annotations.TokenText, Value: "foo", Line: 6, Column: 19},
		},
		{
			test:     "ContinuationValue2",
			expected: annotations.Token{Kind: annotations.TokenText, Value: "bar", Line: 8, Column: 5},
		},
		{
			test:     "MultipleHashes",
			expected: annotations.Token{Kind: annotations.TokenName, Value: "multiple_hashes", Line: 9, Column: 6},
		},
		{
			test:     "MultipleHashesValue",
			expected: annotations.Token{Kind: annotations.TokenText, Value: "foo", Line: 9, Column: 23},
		},
		{
			test:     "MixedHashes",
			expected: annotations.Token{Kind: annotations.TokenName, Value: "mixed_hashes", Line: 10, Column: 10},
		},
		{
			test:     "MixedHashesValue",
			expected: annotations.Token{Kind: annotations.TokenText, Value: "foo", Line: 10, Column: 24},
		},
		{
			test:     "DashedUnderscored",
			expected: annotations.Token{Kind: annotations.TokenName, Value: "dashed-under_scored", Line: 11, Column: 5},
		},
		{
			test:     "DashedUnderscoredValue",
			expected: annotations.Token{Kind: annotations.TokenText, Value: "foo", Line: 11, Column: 26},
		},
		{
			test:     "InvalidName",
			expected: annotations.Token{Kind: annotations.TokenText, Value: "@123 invalid name", Line: 12, Column: 5},
		},
		{
			test:     "Indented",
			expected: annotations.Token{Kind: annotations.TokenName, Value: "indented", Line: 13, Column: 8},
		},
		{
			test:     "IndentedValue",
			expected: annotations.Token{Kind: annotations.TokenText, Value: "foo", Line: 13, Column: 18},
		},
		{
			test:     "Empty",
			expected: annotations.Token{Kind: annotations.TokenName, Value: "empty", Line: 14, Column: 5},
		},
		{
			test:     "EOF",
			expected: annotations.Token{Kind: annotations.TokenEOF, Line: 15, Column: 2},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			actual := scanner.Scan()
			s.Equal(test.expected, actual)
		})
	}
}
