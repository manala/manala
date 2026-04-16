package option

import (
	"errors"

	"github.com/manala/manala/internal/accessor"
	"github.com/manala/manala/internal/json/number"
	"github.com/manala/manala/internal/json/unmarshaler"
	"github.com/manala/manala/internal/path"
	"github.com/manala/manala/internal/schema"

	"github.com/gosimple/slug"
)

const ENUM = "enum"

type Enum struct {
	name   string
	label  string
	help   string
	schema schema.Schema
	path   path.Path
}

func NewEnum(sch schema.Schema, p path.Path) (*Enum, error) {
	// Enum
	enum, ok := sch["enum"].([]any)
	if !ok {
		return nil, errors.New("invalid recipe option enum")
	}

	if len(enum) == 0 {
		return nil, errors.New("empty recipe option enum")
	}

	return &Enum{
		schema: sch,
		path:   p,
	}, nil
}

func (o *Enum) Name() string  { return o.name }
func (o *Enum) Label() string { return o.label }
func (o *Enum) Help() string  { return o.help }

func (o *Enum) Values() []any {
	enum := o.schema["enum"].([]any)
	values := make([]any, len(enum))
	for i := range enum {
		if value, ok := number.NumberType(enum[i]); ok {
			values[i] = value.Normalize()
		} else {
			values[i] = enum[i]
		}
	}
	return values
}

func (o *Enum) Accessor(data any) accessor.Accessor {
	return path.NewAccessor(o.path, data)
}

func (o *Enum) UnmarshalJSON(data []byte) error {
	var env struct {
		Name  string `json:"name"`
		Label string `json:"label"`
		Help  string `json:"help"`
	}
	if err := unmarshaler.Unmarshal(data, &env); err != nil {
		return err
	}

	// Label (required, max=100)
	if env.Label == "" {
		return errors.New("missing option label property")
	} else if len(env.Label) > 100 {
		return errors.New("too long option label field (max=100)")
	}

	// Help (optional, max=100)
	if len(env.Help) > 100 {
		return errors.New("too long option help field (max=100)")
	}

	// Name (optional, max=100)
	if len(env.Name) > 100 {
		return errors.New("too long option name field (max=100)")
	}

	o.label = env.Label
	o.help = env.Help
	o.name = env.Name

	if o.name == "" {
		o.name = slug.Make(o.label)
	}

	return nil
}
