package expectation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type ErrorExpectation interface {
	Expect(t *testing.T, err error)
}

func ExpectError(t *testing.T, expectation ErrorExpectation, err error) {
	t.Helper()

	if expectation == nil {
		require.NoError(t, err)
		return
	}

	expectation.Expect(t, err)
}

func ErrorMessage(msg string) ErrorMessageExpectation { return ErrorMessageExpectation(msg) }

type ErrorMessageExpectation string

func (a ErrorMessageExpectation) Expect(t *testing.T, err error) {
	t.Helper()

	require.EqualError(t, err, string(a))
}

func ErrorEqual(err error) ErrorEqualExpectation { return ErrorEqualExpectation{err} }

type ErrorEqualExpectation struct{ error } //nolint:errname

func (a ErrorEqualExpectation) Expect(t *testing.T, err error) {
	t.Helper()

	require.Equal(t, a.error, err)
}

func ErrorType(err error) ErrorTypeExpectation { return ErrorTypeExpectation{err} }

type ErrorTypeExpectation struct{ error } //nolint:errname

func (a ErrorTypeExpectation) Expect(t *testing.T, err error) {
	t.Helper()

	require.IsType(t, a.error, err)
}

func Errors(expectations ...ErrorExpectation) ErrorsExpectation { return expectations }

type ErrorsExpectation []ErrorExpectation

func (a ErrorsExpectation) Expect(t *testing.T, err error) {
	t.Helper()

	errs := a.flatten(err)
	require.Len(t, errs, len(a), "errors count not equal")
	for i, exp := range a {
		if exp != nil {
			exp.Expect(t, errs[i])
		}
	}
}

func (a ErrorsExpectation) flatten(err error) []error {
	if err == nil {
		return nil
	}
	if multi, ok := err.(interface{ Unwrap() []error }); ok {
		if children := multi.Unwrap(); len(children) > 0 {
			var result []error
			for _, e := range children {
				result = append(result, a.flatten(e)...)
			}
			return result
		}
	}
	return []error{err}
}
