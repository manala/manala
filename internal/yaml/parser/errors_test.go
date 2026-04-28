package parser_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/manala/manala/internal/parsing"
	"github.com/manala/manala/internal/testing/expect"
	"github.com/manala/manala/internal/yaml/parser"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/token"
	"github.com/stretchr/testify/suite"
)

type ErrorsSuite struct{ suite.Suite }

func TestErrorsSuite(t *testing.T) {
	suite.Run(t, new(ErrorsSuite))
}

func (s *ErrorsSuite) TestErrorAt() {
	err := parser.ErrorAt(
		errors.New("error"),
		&token.Token{
			Position: &token.Position{Line: 2, Column: 3},
		},
	)

	expect.Error(s.T(), parsing.ErrorExpectation{
		Line:   2,
		Column: 3,
		Err:    expect.ErrorMessageExpectation("error"),
	}, err)
}

func (s *ErrorsSuite) TestErrorFrom() {
	tests := []struct {
		test     string
		err      error
		expected expect.ErrorExpectation
	}{
		{
			test: "Unknown",
			err:  errors.New("unknown"),
			expected: parsing.ErrorExpectation{
				Line:   0,
				Column: 0,
				Err:    expect.ErrorMessageExpectation("unknown"),
			},
		},
		{
			test: "TypeError",
			err: &yaml.TypeError{
				DstType: reflect.TypeFor[string](),
				Token: &token.Token{
					Position: &token.Position{Line: 3, Column: 5},
				},
			},
			expected: parsing.ErrorExpectation{
				Line:   3,
				Column: 5,
				Err:    expect.ErrorMessageExpectation("field must be a string"),
			},
		},
		{
			test: "SyntaxError",
			err: &yaml.SyntaxError{
				Message: "syntax error",
				Token: &token.Token{
					Position: &token.Position{Line: 2, Column: 4},
				},
			},
			expected: parsing.ErrorExpectation{
				Line:   2,
				Column: 4,
				Err:    expect.ErrorMessageExpectation("syntax error"),
			},
		},
		{
			test: "DuplicateKeyError",
			err: &yaml.DuplicateKeyError{
				Message: "duplicate key",
				Token: &token.Token{
					Position: &token.Position{Line: 4, Column: 1},
				},
			},
			expected: parsing.ErrorExpectation{
				Line:   4,
				Column: 1,
				Err:    expect.ErrorMessageExpectation("duplicate key"),
			},
		},
		{
			test: "OverflowError",
			err: &yaml.OverflowError{
				DstType: reflect.TypeFor[int8](),
				SrcNum:  "999",
				Token: &token.Token{
					Position: &token.Position{Line: 1, Column: 3},
				},
			},
			expected: parsing.ErrorExpectation{
				Line:   1,
				Column: 3,
				Err:    expect.ErrorMessageExpectation("cannot unmarshal 999 into Go value of type int8 ( overflow )"),
			},
		},
		{
			test: "UnknownFieldError",
			err: &yaml.UnknownFieldError{
				Message: "unknown field",
				Token: &token.Token{
					Position: &token.Position{Line: 5, Column: 2},
				},
			},
			expected: parsing.ErrorExpectation{
				Line:   5,
				Column: 2,
				Err:    expect.ErrorMessageExpectation("unknown field"),
			},
		},
		{
			test: "UnexpectedNodeTypeError",
			err: &yaml.UnexpectedNodeTypeError{
				Actual:   ast.StringType,
				Expected: ast.MappingType,
				Token: &token.Token{
					Position: &token.Position{Line: 2, Column: 1},
				},
			},
			expected: parsing.ErrorExpectation{
				Line:   2,
				Column: 1,
				Err:    expect.ErrorMessageExpectation("string was used where mapping is expected"),
			},
		},
		{
			test: "ExceededMaxDepth",
			err:  yaml.ErrExceededMaxDepth,
			expected: parsing.ErrorExpectation{
				Line:   0,
				Column: 0,
				Err:    expect.ErrorMessageExpectation("yaml exceeded max depth"),
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			err := parser.ErrorFrom(test.err)

			expect.Error(s.T(), test.expected, err)
		})
	}
}
