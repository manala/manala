package errors_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/manala/manala/internal/testing/expectation"
	yamlerrors "github.com/manala/manala/internal/yaml/errors"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/token"
	"github.com/stretchr/testify/suite"
)

type FromSuite struct{ suite.Suite }

func TestFromSuite(t *testing.T) {
	suite.Run(t, new(FromSuite))
}

func (s *FromSuite) Test() {
	tests := []struct {
		test     string
		err      error
		expected expectation.ErrorExpectation
	}{
		{
			test:     "Unknown",
			err:      errors.New("unknown"),
			expected: expectation.ErrorEqual(errors.New("unknown")),
		},
		{
			test: "TypeError",
			err: &yaml.TypeError{
				DstType: reflect.TypeFor[string](),
				Token: &token.Token{
					Position: &token.Position{Line: 3, Column: 5},
				},
			},
			expected: yamlerrors.Expectation{
				Position: [2]int{3, 5},
				Err:      expectation.ErrorMessage("field must be a string"),
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
			expected: yamlerrors.Expectation{
				Position: [2]int{2, 4},
				Err:      expectation.ErrorMessage("syntax error"),
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
			expected: yamlerrors.Expectation{
				Position: [2]int{4, 1},
				Err:      expectation.ErrorMessage("duplicate key"),
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
			expected: yamlerrors.Expectation{
				Position: [2]int{1, 3},
				Err:      expectation.ErrorMessage("cannot unmarshal 999 into Go value of type int8 ( overflow )"),
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
			expected: yamlerrors.Expectation{
				Position: [2]int{5, 2},
				Err:      expectation.ErrorMessage("unknown field"),
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
			expected: yamlerrors.Expectation{
				Position: [2]int{2, 1},
				Err:      expectation.ErrorMessage("string was used where mapping is expected"),
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			err := yamlerrors.From(test.err)

			expectation.ExpectError(s.T(), test.expected, err)
		})
	}
}
