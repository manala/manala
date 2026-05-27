package annotation_test

import (
	"testing"

	yamlannotation "github.com/manala/manala/internal/yaml/annotation"

	"github.com/stretchr/testify/suite"
)

type AnnotationSuite struct{ suite.Suite }

func TestAnnotationSuite(t *testing.T) {
	suite.Run(t, new(AnnotationSuite))
}

func (s *AnnotationSuite) TestName() {
	name := yamlannotation.Name{
		Token: yamlannotation.Token{Value: "foo"},
	}

	s.Equal("foo", name.String())
}

func (s *AnnotationSuite) TestBodyString() {
	s.Run("Single", func() {
		body := yamlannotation.Body{
			Tokens: []yamlannotation.Token{
				{Value: "bar", Line: 1, Column: 7},
			},
		}
		s.Equal("bar", body.String())
	})

	s.Run("Multiline", func() {
		body := yamlannotation.Body{
			Tokens: []yamlannotation.Token{
				{Value: "bar", Line: 1, Column: 7},
				{Value: "baz", Line: 2, Column: 3},
			},
		}
		s.Equal("bar\nbaz", body.String())
	})

	s.Run("Empty", func() {
		body := yamlannotation.Body{}
		s.Empty(body.String())
	})
}

func (s *AnnotationSuite) TestBodyStencil() {
	s.Run("Single", func() {
		body := yamlannotation.Body{
			Tokens: []yamlannotation.Token{
				{Value: "bar", Line: 1, Column: 8},
			},
		}
		s.Equal("       bar", body.Stencil())
	})

	s.Run("Multiline", func() {
		body := yamlannotation.Body{
			Tokens: []yamlannotation.Token{
				{Value: "{", Line: 2, Column: 8},
				{Value: "\"bar\":", Line: 3, Column: 5},
				{Value: "123", Line: 3, Column: 12},
				{Value: "}", Line: 4, Column: 3},
			},
		}
		s.Equal("\n       {\n    \"bar\":           123\n  }", body.Stencil())
	})
}

func (s *AnnotationSuite) TestBodyStart() {
	token1 := yamlannotation.Token{Value: "foo", Line: 2, Column: 5}
	token2 := yamlannotation.Token{Value: "bar", Line: 3, Column: 3}

	body := yamlannotation.Body{
		Tokens: []yamlannotation.Token{token1, token2},
	}

	s.Equal(token1, body.Start())
}
