package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Expectation struct {
	Type  any
	Attrs [][2]any
}

func (a Expectation) Expect(t *testing.T, err error) {
	t.Helper()

	require.IsType(t, a.Type, err)

	// Attrs
	if _err, ok := err.(interface{ Attrs() [][2]any }); ok {
		assert.Equal(t, a.Attrs, _err.Attrs())
	} else if a.Attrs != nil {
		assert.Fail(t, "Error does not contain attrs")
	}
}
