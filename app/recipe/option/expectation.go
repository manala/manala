package option

import (
	"testing"

	"github.com/manala/manala/app"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Expectation struct {
	Type  any
	Label string
	Name  string
	// String
	MaxLength int
	// Enum
	Values []any
}

func (a Expectation) Expect(t *testing.T, opt app.RecipeOption) {
	t.Helper()

	require.IsType(t, a.Type, opt)

	assert.Equal(t, a.Label, opt.Label(), "label not equal")
	assert.Equal(t, a.Name, opt.Name(), "name not equal")

	// String
	if opt, ok := opt.(*String); ok {
		assert.Equal(t, a.MaxLength, opt.MaxLength(), "maxLength not equal")
	}

	// Enum
	if opt, ok := opt.(*Enum); ok {
		assert.Equal(t, a.Values, opt.Values(), "values not equals")
	}
}

func ExpectOption(t *testing.T, expectation Expectation, opt app.RecipeOption) {
	t.Helper()
	expectation.Expect(t, opt)
}

type Expectations []Expectation

func (a Expectations) Expect(t *testing.T, opts []app.RecipeOption) {
	t.Helper()

	require.Len(t, opts, len(a), "options count not equal")

	for i, expectation := range a {
		expectation.Expect(t, opts[i])
	}
}

func ExpectOptions(t *testing.T, expectations Expectations, opts []app.RecipeOption) {
	t.Helper()
	expectations.Expect(t, opts)
}
