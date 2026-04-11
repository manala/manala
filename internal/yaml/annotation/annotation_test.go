package annotation_test

import (
	"testing"

	"github.com/manala/manala/internal/yaml/annotation"

	"github.com/stretchr/testify/suite"
)

type AnnotationSuite struct{ suite.Suite }

func TestAnnotationSuite(t *testing.T) {
	suite.Run(t, new(AnnotationSuite))
}

func (s *AnnotationSuite) TestName() {
	name := annotation.Name{
		Token: annotation.Token{Value: "foo"},
	}

	s.Equal("foo", name.String())
}

func (s *AnnotationSuite) TestValueString() {
	s.Run("Single", func() {
		value := annotation.Value{
			Tokens: []annotation.Token{
				{Value: "bar", Line: 1, Column: 7},
			},
		}
		s.Equal("bar", value.String())
	})

	s.Run("Multiline", func() {
		value := annotation.Value{
			Tokens: []annotation.Token{
				{Value: "bar", Line: 1, Column: 7},
				{Value: "baz", Line: 2, Column: 3},
			},
		}
		s.Equal("bar\nbaz", value.String())
	})

	s.Run("Empty", func() {
		value := annotation.Value{}
		s.Empty(value.String())
	})
}

func (s *AnnotationSuite) TestValueStencil() {
	s.Run("Single", func() {
		value := annotation.Value{
			Tokens: []annotation.Token{
				{Value: "bar", Line: 1, Column: 8},
			},
		}
		s.Equal("       bar", value.Stencil())
	})

	s.Run("Multiline", func() {
		value := annotation.Value{
			Tokens: []annotation.Token{
				{Value: "{", Line: 2, Column: 8},
				{Value: "\"bar\":", Line: 3, Column: 5},
				{Value: "123", Line: 3, Column: 12},
				{Value: "}", Line: 4, Column: 3},
			},
		}
		s.Equal("\n       {\n    \"bar\":           123\n  }", value.Stencil())
	})
}
