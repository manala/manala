package yaml

import (
	goYamlAst "github.com/goccy/go-yaml/ast"
	goYamlParser "github.com/goccy/go-yaml/parser"
	"manala/internal/path"
	"manala/internal/serrors"
)

func (s *Suite) TestNodePath() {
	node := goYamlAst.Null(nil)

	tests := []struct {
		test     string
		path     string
		expected string
	}{
		{
			test:     "Root",
			path:     "$",
			expected: "",
		},
		{
			test:     "FirstLevel",
			path:     "$.foo",
			expected: "foo",
		},
		{
			test:     "Levels",
			path:     "$.foo.bar",
			expected: "foo.bar",
		},
		{
			test:     "Index",
			path:     "$.foo[0].bar",
			expected: "foo[0].bar",
		},
		{
			test:     "IndexLast",
			path:     "$.foo[0]",
			expected: "foo[0]",
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			node.Path = test.path

			path := NewNodePath(node)

			s.Equal(test.expected, path.String())
		})
	}
}

func (s *Suite) TestNodePathAccessorGetErrors() {
	node, _ := goYamlParser.ParseBytes([]byte(`
sequence_empty: {}
sequence_single:
  first: foo
sequence_multiple:
  first: foo
  last: bar
`), 0)

	tests := []struct {
		test     string
		path     string
		expected *serrors.Assert
	}{
		{
			test: "Root",
			path: "baz",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "unable to access yaml path",
				Arguments: []any{
					"path", "baz",
				},
			},
		},
		{
			test: "Leaf",
			path: "bar.bar",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "unable to access yaml path",
				Arguments: []any{
					"path", "bar.bar",
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			accessor := NewNodePathAccessor(path.Path(test.path))

			_node, err := accessor.Get(node.Docs[0].Body)

			serrors.Equal(s.Assert(), test.expected, err)
			s.Nil(_node)
		})
	}
}

func (s *Suite) TestNodePathAccessorGet() {
	node, _ := goYamlParser.ParseBytes([]byte(`
sequence_empty: {}
sequence_single:
  first: foo
sequence_multiple:
  first: foo
  last: bar
`), 0)

	tests := []struct {
		test           string
		path           string
		expectedLine   int
		expectedColumn int
	}{
		{
			test:           "SequenceEmpty",
			path:           "sequence_empty",
			expectedLine:   2,
			expectedColumn: 17,
		},
		{
			test:           "SequenceSingle",
			path:           "sequence_single",
			expectedLine:   4,
			expectedColumn: 8,
		},
		{
			test:           "SequenceSingleFirst",
			path:           "sequence_single.first",
			expectedLine:   4,
			expectedColumn: 10,
		},
		{
			test:           "SequenceMultiple",
			path:           "sequence_multiple",
			expectedLine:   6,
			expectedColumn: 8,
		},
		{
			test:           "SequenceMultipleFirst",
			path:           "sequence_multiple.first",
			expectedLine:   6,
			expectedColumn: 10,
		},
		{
			test:           "SequenceMultipleLast",
			path:           "sequence_multiple.last",
			expectedLine:   7,
			expectedColumn: 9,
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			accessor := NewNodePathAccessor(path.Path(test.path))

			_node, err := accessor.Get(node.Docs[0].Body)

			s.NoError(err)
			s.NotNil(node)

			token := _node.GetToken()

			s.Equal(test.expectedLine, token.Position.Line)
			s.Equal(test.expectedColumn, token.Position.Column)
		})
	}
}
