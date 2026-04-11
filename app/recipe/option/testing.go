package option

import (
	"testing"

	"github.com/manala/manala/app"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Assertion struct {
	Type  any
	Label string
	Name  string
	// String
	MaxLength int
	// Enum
	Values []any
}

func (a *Assertion) Assert(t *testing.T, opt app.RecipeOption) {
	t.Helper()

	require.IsType(t, a.Type, opt)
	assert.Equal(t, a.Label, opt.Label(), "Label not equal")
	assert.Equal(t, a.Name, opt.Name(), "Name not equal")

	// String
	if opt, ok := opt.(*String); ok {
		assert.Equal(t, a.MaxLength, opt.MaxLength(), "MaxLength not equal")
	}

	// Enum
	if opt, ok := opt.(*Enum); ok {
		assert.Equal(t, a.Values, opt.Values(), "Values not equals")
	}
}

func Equal(t *testing.T, assertion Assertion, opt app.RecipeOption) {
	t.Helper()
	assertion.Assert(t, opt)
}

type Assertions []Assertion

func (a *Assertions) Assert(t *testing.T, opts []app.RecipeOption) {
	t.Helper()

	require.Len(t, opts, len(*a), "Incorrect options length")

	for i, assertion := range *a {
		assertion.Assert(t, opts[i])
	}
}

func Equals(t *testing.T, assertions Assertions, opts []app.RecipeOption) {
	t.Helper()
	assertions.Assert(t, opts)
}
