package parser_test

import (
	"reflect"
	"testing"

	"github.com/manala/manala/internal/parsing"
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/testing/errors"
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
		serrors.New("error"),
		&token.Token{
			Position: &token.Position{Line: 2, Column: 3},
		},
	)

	errors.Equal(s.T(), &parsing.ErrorAssertion{
		Line:   2,
		Column: 3,
		Err: &serrors.Assertion{
			Message: "error",
		},
	}, err)
}

func (s *ErrorsSuite) TestErrorFrom() {
	tests := []struct {
		test     string
		err      error
		expected errors.Assertion
	}{
		{
			test: "Unknown",
			err:  serrors.New("unknown"),
			expected: &parsing.ErrorAssertion{
				Line:   0,
				Column: 0,
				Err: &serrors.Assertion{
					Message: "unknown",
				},
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
			expected: &parsing.ErrorAssertion{
				Line:   3,
				Column: 5,
				Err: &serrors.Assertion{
					Message: "field must be a string",
				},
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
			expected: &parsing.ErrorAssertion{
				Line:   2,
				Column: 4,
				Err: &serrors.Assertion{
					Message: "syntax error",
				},
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
			expected: &parsing.ErrorAssertion{
				Line:   4,
				Column: 1,
				Err: &serrors.Assertion{
					Message: "duplicate key",
				},
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
			expected: &parsing.ErrorAssertion{
				Line:   1,
				Column: 3,
				Err: &serrors.Assertion{
					Message: "cannot unmarshal 999 into Go value of type int8 ( overflow )",
				},
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
			expected: &parsing.ErrorAssertion{
				Line:   5,
				Column: 2,
				Err: &serrors.Assertion{
					Message: "unknown field",
				},
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
			expected: &parsing.ErrorAssertion{
				Line:   2,
				Column: 1,
				Err: &serrors.Assertion{
					Message: "string was used where mapping is expected",
				},
			},
		},
		{
			test: "ExceededMaxDepth",
			err:  yaml.ErrExceededMaxDepth,
			expected: &parsing.ErrorAssertion{
				Line:   0,
				Column: 0,
				Err: &serrors.Assertion{
					Message: "yaml exceeded max depth",
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			err := parser.ErrorFrom(test.err)

			errors.Equal(s.T(), test.expected, err)
		})
	}
}
