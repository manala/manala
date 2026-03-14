package option

import (
	"github.com/manala/manala/internal/json"
	"github.com/manala/manala/internal/serrors"
)

type TextOption struct {
	*option

	MaxLength int
}

func NewTextOption(option *option, _ map[string]any) (*TextOption, error) {
	// Schema type *MUST* be string
	if t, ok := option.schema["type"]; !ok || t != "string" {
		return nil, serrors.New("invalid recipe option string type").
			WithArguments("label", option.label)
	}

	// Option
	textOption := &TextOption{
		option: option,
	}

	// Max length
	if maxLength, ok := json.NumberType(option.schema["maxLength"]); ok {
		textOption.MaxLength = maxLength.Int()
	}

	return textOption, nil
}
