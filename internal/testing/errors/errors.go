package errors

import "testing"

type Assertion interface {
	Assert(t *testing.T, err error)
}

func Equal(t *testing.T, assertion Assertion, err error) {
	t.Helper()
	assertion.Assert(t, err)
}
