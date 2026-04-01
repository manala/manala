package annotation_test

import (
	"testing"

	"github.com/manala/manala/internal/yaml/annotation"

	"github.com/stretchr/testify/suite"
)

type ScannerSuite struct{ suite.Suite }

func TestScannerSuite(t *testing.T) {
	suite.Run(t, new(ScannerSuite))
}

func (s *ScannerSuite) Test() {
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
	scanner := annotation.NewScanner(src)

	tests := []struct {
		test     string
		expected annotation.Token
	}{
		{
			test:     "Text",
			expected: annotation.Token{Kind: annotation.TokenText, Value: "text before annotation", Line: 2, Column: 5},
		},
		{
			test:     "Line",
			expected: annotation.Token{Kind: annotation.TokenName, Value: "line", Line: 3, Column: 5},
		},
		{
			test:     "LineValue",
			expected: annotation.Token{Kind: annotation.TokenText, Value: "foo", Line: 3, Column: 11},
		},
		{
			test:     "Multiline",
			expected: annotation.Token{Kind: annotation.TokenName, Value: "multiline", Line: 4, Column: 5},
		},
		{
			test:     "MultilineValue1",
			expected: annotation.Token{Kind: annotation.TokenText, Value: "foo", Line: 4, Column: 16},
		},
		{
			test:     "MultilineValue2",
			expected: annotation.Token{Kind: annotation.TokenText, Value: "bar", Line: 5, Column: 5},
		},
		{
			test:     "Continuation",
			expected: annotation.Token{Kind: annotation.TokenName, Value: "continuation", Line: 6, Column: 5},
		},
		{
			test:     "ContinuationValue1",
			expected: annotation.Token{Kind: annotation.TokenText, Value: "foo", Line: 6, Column: 19},
		},
		{
			test:     "ContinuationValue2",
			expected: annotation.Token{Kind: annotation.TokenText, Value: "bar", Line: 8, Column: 5},
		},
		{
			test:     "MultipleHashes",
			expected: annotation.Token{Kind: annotation.TokenName, Value: "multiple_hashes", Line: 9, Column: 6},
		},
		{
			test:     "MultipleHashesValue",
			expected: annotation.Token{Kind: annotation.TokenText, Value: "foo", Line: 9, Column: 23},
		},
		{
			test:     "MixedHashes",
			expected: annotation.Token{Kind: annotation.TokenName, Value: "mixed_hashes", Line: 10, Column: 10},
		},
		{
			test:     "MixedHashesValue",
			expected: annotation.Token{Kind: annotation.TokenText, Value: "foo", Line: 10, Column: 24},
		},
		{
			test:     "DashedUnderscored",
			expected: annotation.Token{Kind: annotation.TokenName, Value: "dashed-under_scored", Line: 11, Column: 5},
		},
		{
			test:     "DashedUnderscoredValue",
			expected: annotation.Token{Kind: annotation.TokenText, Value: "foo", Line: 11, Column: 26},
		},
		{
			test:     "InvalidName",
			expected: annotation.Token{Kind: annotation.TokenText, Value: "@123 invalid name", Line: 12, Column: 5},
		},
		{
			test:     "Indented",
			expected: annotation.Token{Kind: annotation.TokenName, Value: "indented", Line: 13, Column: 8},
		},
		{
			test:     "IndentedValue",
			expected: annotation.Token{Kind: annotation.TokenText, Value: "foo", Line: 13, Column: 18},
		},
		{
			test:     "Empty",
			expected: annotation.Token{Kind: annotation.TokenName, Value: "empty", Line: 14, Column: 5},
		},
		{
			test:     "EOF",
			expected: annotation.Token{Kind: annotation.TokenEOF, Line: 15, Column: 2},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			actual := scanner.Scan()
			s.Equal(test.expected, actual)
		})
	}
}
