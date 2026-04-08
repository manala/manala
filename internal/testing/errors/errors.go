package errors

import "testing"

type Assertion interface {
	AssertError(t *testing.T, err error)
}

func Equal(t *testing.T, assertion Assertion, err error) {
	t.Helper()
	assertion.AssertError(t, err)
}
