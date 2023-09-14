package heredoc

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type Suite struct{ suite.Suite }

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestDoc() {
	tests := []struct {
		test     string
		doc      string
		args     []any
		expected string
	}{
		{
			test:     "Empty",
			doc:      ``,
			expected: "",
		},
		{
			test: "Eol",
			doc: `
			`,
			expected: "\n",
		},
		{
			test:     "String",
			doc:      `Foo Bar`,
			expected: "Foo Bar",
		},
		{
			test: "Simple",
			doc: `
				Foo
				Bar
			`,
			expected: "Foo\nBar\n",
		},
		{
			test: "WithoutTrailingLines",
			doc: `Foo
				Bar`,
			expected: "Foo\nBar",
		},
		{
			test: "SpaceIndentation",
			doc: `
				Foo
				 Bar
				  Baz
			`,
			expected: "Foo\n Bar\n  Baz\n",
		},
		{
			test: "MultipleIndentation",
			doc: `
				Foo
					Bar
						Hoge
							`,
			expected: "Foo\n\tBar\n\t\tHoge\n",
		},
		{
			test:     "Args",
			doc:      `%[1]s %[2]d`,
			args:     []any{"foo", 123},
			expected: "foo 123",
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			s.Equal(test.expected, Doc(test.doc, test.args...))
		})
	}
}
