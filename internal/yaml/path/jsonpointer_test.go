package path_test

import (
	"testing"

	yamlpath "github.com/manala/manala/internal/yaml/path"

	"github.com/stretchr/testify/suite"
)

type JSONPointerSuite struct{ suite.Suite }

func TestJSONPointerSuite(t *testing.T) {
	suite.Run(t, new(JSONPointerSuite))
}

func (s *JSONPointerSuite) TestFromJSONPointer() {
	tests := []struct {
		test     string
		pointer  string
		expected string
	}{
		{test: "Empty", pointer: "", expected: "$"},
		{test: "Root", pointer: "/", expected: "$"},
		{test: "Object", pointer: "/foo", expected: "$.foo"},
		{test: "ObjectNested", pointer: "/foo/bar", expected: "$.foo.bar"},
		{test: "Array", pointer: "/foo/0", expected: "$.foo[0]"},
		{test: "ArrayFirst", pointer: "/0", expected: "$[0]"},
		{test: "ArrayNested", pointer: "/foo/0/bar", expected: "$.foo[0].bar"},
		{test: "Deep", pointer: "/a/b/c", expected: "$.a.b.c"},
		{test: "PointerEscapeTilde", pointer: "/a~0b", expected: "$.a~b"},
		{test: "PointerEscapeSlash", pointer: "/a~1b", expected: "$.a/b"},
	}
	for _, test := range tests {
		s.Run(test.test, func() {
			s.Equal(test.expected, yamlpath.FromJSONPointer(test.pointer))
		})
	}
}

func (s *JSONPointerSuite) TestToJSONPointer() {
	tests := []struct {
		test     string
		path     string
		expected string
	}{
		{test: "Root", path: "$", expected: ""},
		{test: "Object", path: "$.foo", expected: "/foo"},
		{test: "ObjectNested", path: "$.foo.bar", expected: "/foo/bar"},
		{test: "Array", path: "$.foo[0]", expected: "/foo/0"},
		{test: "ArrayFirst", path: "$[0]", expected: "/0"},
		{test: "ArrayNested", path: "$.foo[0].bar", expected: "/foo/0/bar"},
		{test: "Deep", path: "$.a.b.c", expected: "/a/b/c"},
		{test: "PointerEscapeTilde", path: "$.a~b", expected: "/a~0b"},
	}
	for _, test := range tests {
		s.Run(test.test, func() {
			s.Equal(test.expected, yamlpath.ToJSONPointer(test.path))
		})
	}
}
