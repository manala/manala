package expect

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type ErrorExpectation interface {
	Expect(t *testing.T, err error)
}

func Error(t *testing.T, expectation ErrorExpectation, err error) {
	t.Helper()

	if expectation == nil {
		require.NoError(t, err)
		return
	}

	expectation.Expect(t, err)
}

type ErrorMessageExpectation string

func (a ErrorMessageExpectation) Expect(t *testing.T, err error) {
	t.Helper()

	require.EqualError(t, err, string(a))
}
