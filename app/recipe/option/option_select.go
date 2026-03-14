package option

import (
	"github.com/manala/manala/internal/json"
	"github.com/manala/manala/internal/serrors"
)

type SelectOption struct {
	*option

	Values []any
}

func NewSelectOption(option *option, _ map[string]any) (*SelectOption, error) {
	// Option
	selectOption := &SelectOption{
		option: option,
	}

	// Enum
	enum, ok := option.schema["enum"].([]any)
	if !ok {
		return nil, serrors.New("invalid recipe option enum").
			WithArguments("label", option.label)
	}

	if len(enum) == 0 {
		return nil, serrors.New("empty recipe option enum").
			WithArguments("label", option.label)
	}

	// Values
	selectOption.Values = make([]any, len(enum))

	for i := range enum {
		if value, ok := json.NumberType(enum[i]); ok {
			selectOption.Values[i] = value.Normalize()
		} else {
			selectOption.Values[i] = enum[i]
		}
	}

	return selectOption, nil
}
