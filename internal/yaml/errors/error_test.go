package errors_test

import (
	"errors"
	"testing"

	"github.com/manala/manala/internal/testing/expectation"
	yamlerrors "github.com/manala/manala/internal/yaml/errors"
	yamlerrorstest "github.com/manala/manala/internal/yaml/errors/errorstest"

	"github.com/goccy/go-yaml/token"
	"github.com/stretchr/testify/suite"
)

type ErrorsSuite struct{ suite.Suite }

func TestErrorsSuite(t *testing.T) {
	suite.Run(t, new(ErrorsSuite))
}

func (s *ErrorsSuite) Test() {
	err := yamlerrors.New(
		errors.New("error"),
		&token.Token{
			Position: &token.Position{Line: 2, Column: 3},
		},
	)

	expectation.ExpectError(s.T(), yamlerrorstest.Expectation{
		Position: [2]int{2, 3},
		Err:      expectation.ErrorMessage("error"),
	}, err)
}
