package source

import (
	"testing"

	"github.com/manala/manala/internal/output"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Expectation string

func (a Expectation) Expect(t *testing.T, err error) {
	t.Helper()

	require.IsType(t, Error{}, err)
	e := err.(Error)

	assert.Equal(t, string(a), e.Render(output.Plain))
}
