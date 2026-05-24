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
		value    [2]int
		property [2]int
	}{
		{
			test:     "Empty",
			bytes:    `{"foo": "bar"}`,
			location: "",
			value:    [2]int{0, 0},
			property: [2]int{0, 0},
		},
		{
			test:     "NotFound",
			bytes:    `{"foo": "bar"}`,
			location: "/baz",
			value:    [2]int{0, 0},
			property: [2]int{0, 0},
		},
		{
			test:     "Object",
			bytes:    `{"foo": "bar"}`,
			location: "/foo",
			value:    [2]int{1, 9},
			property: [2]int{1, 2},
		},
		{
			test:     "ObjectSecond",
			bytes:    `{"foo": "bar", "baz": "qux"}`,
			location: "/baz",
			value:    [2]int{1, 23},
			property: [2]int{1, 16},
		},
		{
			test:     "ObjectNested",
			bytes:    `{"foo": {"bar": "baz"}}`,
			location: "/foo/bar",
			value:    [2]int{1, 17},
			property: [2]int{1, 10},
		},
		{
			test:     "ArrayFirst",
			bytes:    `{"foo": [1, 2, 3]}`,
			location: "/foo/0",
			value:    [2]int{1, 10},
			property: [2]int{1, 10},
		},
		{
			test:     "ArraySecond",
			bytes:    `{"foo": [1, 2, 3]}`,
			location: "/foo/1",
			value:    [2]int{1, 13},
			property: [2]int{1, 13},
		},
		{
			test:     "ArrayThird",
			bytes:    `{"foo": [1, 2, 3]}`,
			location: "/foo/2",
			value:    [2]int{1, 16},
			property: [2]int{1, 16},
		},
		{
			test: "Multiline",
			bytes: heredoc.Doc(`
				{
				  "foo": "bar"
				}
			`),
			location: "/foo",
			value:    [2]int{2, 10},
			property: [2]int{2, 3},
		},
		{
			test:     "PointerEscapeTilde",
			bytes:    `{"a~b": 1}`,
			location: "/a~0b",
			value:    [2]int{1, 9},
			property: [2]int{1, 2},
		},
		{
			test:     "PointerEscapeSlash",
			bytes:    `{"a~b": 1, "a/b": 2}`,
			location: "/a~1b",
			value:    [2]int{1, 19},
			property: [2]int{1, 12},
		},
		{
			test:     "RootArray",
			bytes:    `[1, 2, 3]`,
			location: "/0",
			value:    [2]int{1, 2},
			property: [2]int{1, 2},
		},
		{
			test:     "RootArraySecond",
			bytes:    `[1, 2, 3]`,
			location: "/1",
			value:    [2]int{1, 5},
			property: [2]int{1, 5},
		},
		{
			test:     "ArrayOfObjects",
			bytes:    `[{"foo": 1}, {"foo": 2}]`,
			location: "/1/foo",
			value:    [2]int{1, 22},
			property: [2]int{1, 15},
		},
		{
			test:     "ObjectDeep",
			bytes:    `{"a": {"b": {"c": 1}}}`,
			location: "/a/b/c",
			value:    [2]int{1, 19},
			property: [2]int{1, 14},
		},
		{
			test:     "ArrayOutOfBounds",
			bytes:    `{"foo": [1, 2, 3]}`,
			location: "/foo/5",
			value:    [2]int{0, 0},
			property: [2]int{0, 0},
		},
		{
			test:     "ValueBoolean",
			bytes:    `{"foo": true}`,
			location: "/foo",
			value:    [2]int{1, 9},
			property: [2]int{1, 2},
		},
		{
			test:     "ValueNull",
			bytes:    `{"foo": null}`,
			location: "/foo",
			value:    [2]int{1, 9},
			property: [2]int{1, 2},
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
			value:    [2]int{5, 5},
			property: [2]int{5, 5},
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
			value:    [2]int{3, 12},
			property: [2]int{3, 5},
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
			value:    [2]int{3, 10},
			property: [2]int{3, 3},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			l := validation.Locator{Bytes: []byte(test.bytes)}

			// Value
			line, column := l.ValueAt(test.location)
			s.Equal(test.value[0], line, "value line not equal")
			s.Equal(test.value[1], column, "value column not equal")

			// Property
			line, column = l.PropertyAt(test.location)
			s.Equal(test.property[0], line, "property line not equal")
			s.Equal(test.property[1], column, "property column not equal")
		})
	}
}
