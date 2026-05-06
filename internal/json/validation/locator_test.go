package validation_test

import (
	"testing"

	"github.com/manala/manala/internal/json/validation"
	"github.com/manala/manala/internal/testing/heredoc"

	"github.com/stretchr/testify/suite"
)

type LocatorSuite struct{ suite.Suite }

func TestLocatorSuite(t *testing.T) {
	suite.Run(t, new(LocatorSuite))
}

func (s *LocatorSuite) TestAt() {
	tests := []struct {
		test     string
		bytes    string
		location string
		line     int
		column   int
	}{
		{
			test:     "Empty",
			bytes:    `{"foo": "bar"}`,
			location: "",
			line:     0,
			column:   0,
		},
		{
			test:     "NotFound",
			bytes:    `{"foo": "bar"}`,
			location: "/baz",
			line:     0,
			column:   0,
		},
		{
			test:     "Object",
			bytes:    `{"foo": "bar"}`,
			location: "/foo",
			line:     1,
			column:   9,
		},
		{
			test:     "ObjectSecond",
			bytes:    `{"foo": "bar", "baz": "qux"}`,
			location: "/baz",
			line:     1,
			column:   23,
		},
		{
			test:     "ObjectNested",
			bytes:    `{"foo": {"bar": "baz"}}`,
			location: "/foo/bar",
			line:     1,
			column:   17,
		},
		{
			test:     "ArrayFirst",
			bytes:    `{"foo": [1, 2, 3]}`,
			location: "/foo/0",
			line:     1,
			column:   10,
		},
		{
			test:     "ArraySecond",
			bytes:    `{"foo": [1, 2, 3]}`,
			location: "/foo/1",
			line:     1,
			column:   13,
		},
		{
			test:     "ArrayThird",
			bytes:    `{"foo": [1, 2, 3]}`,
			location: "/foo/2",
			line:     1,
			column:   16,
		},
		{
			test: "Multiline",
			bytes: heredoc.Doc(`
				{
				  "foo": "bar"
				}
			`),
			location: "/foo",
			line:     2,
			column:   10,
		},
		{
			test:     "PointerEscapeTilde",
			bytes:    `{"a~b": 1}`,
			location: "/a~0b",
			line:     1,
			column:   9,
		},
		{
			test:     "PointerEscapeSlash",
			bytes:    `{"a~b": 1, "a/b": 2}`,
			location: "/a~1b",
			line:     1,
			column:   19,
		},
		{
			test:     "RootArray",
			bytes:    `[1, 2, 3]`,
			location: "/0",
			line:     1,
			column:   2,
		},
		{
			test:     "RootArraySecond",
			bytes:    `[1, 2, 3]`,
			location: "/1",
			line:     1,
			column:   5,
		},
		{
			test:     "ArrayOfObjects",
			bytes:    `[{"foo": 1}, {"foo": 2}]`,
			location: "/1/foo",
			line:     1,
			column:   22,
		},
		{
			test:     "ObjectDeep",
			bytes:    `{"a": {"b": {"c": 1}}}`,
			location: "/a/b/c",
			line:     1,
			column:   19,
		},
		{
			test:     "ArrayOutOfBounds",
			bytes:    `{"foo": [1, 2, 3]}`,
			location: "/foo/5",
			line:     0,
			column:   0,
		},
		{
			test:     "ValueBoolean",
			bytes:    `{"foo": true}`,
			location: "/foo",
			line:     1,
			column:   9,
		},
		{
			test:     "ValueNull",
			bytes:    `{"foo": null}`,
			location: "/foo",
			line:     1,
			column:   9,
		},
		{
			test: "MultilineArray",
			bytes: heredoc.Doc(`
				{
				  "foo": [
				    1,
				    2,
				    3
				  ]
				}
			`),
			location: "/foo/2",
			line:     5,
			column:   5,
		},
		{
			test: "MultilineNested",
			bytes: heredoc.Doc(`
				{
				  "foo": {
				    "bar": "baz"
				  }
				}
			`),
			location: "/foo/bar",
			line:     3,
			column:   12,
		},
		{
			test: "MultilineSecondKey",
			bytes: heredoc.Doc(`
				{
				  "foo": "bar",
				  "baz": "qux"
				}
			`),
			location: "/baz",
			line:     3,
			column:   10,
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			l := validation.Locator{Bytes: []byte(test.bytes)}
			line, column := l.At(test.location)
			s.Equal(test.line, line)
			s.Equal(test.column, column)
		})
	}
}
