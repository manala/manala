package option

import (
	"github.com/manala/manala/internal/json/number"
	"github.com/manala/manala/internal/json/unmarshaler"
	"github.com/manala/manala/internal/path"
	"github.com/manala/manala/internal/schema"
	"github.com/manala/manala/internal/serrors"

	"github.com/gosimple/slug"
)

type SelectOption struct {
	name   string
	label  string
	help   string
	schema schema.Schema
	path   path.Path

	Values []any
}

func NewSelectOption(sch schema.Schema, p path.Path) (*SelectOption, error) {
	// Enum
	enum, ok := sch["enum"].([]any)
	if !ok {
		return nil, serrors.New("invalid recipe option enum")
	}

	if len(enum) == 0 {
		return nil, serrors.New("empty recipe option enum")
	}

	// Values
	values := make([]any, len(enum))
	for i := range enum {
		if value, ok := number.NumberType(enum[i]); ok {
			values[i] = value.Normalize()
		} else {
			values[i] = enum[i]
		}
	}

	return &SelectOption{
		schema: sch,
		path:   p,
		Values: values,
	}, nil
}

func (o *SelectOption) Name() string          { return o.name }
func (o *SelectOption) Label() string         { return o.label }
func (o *SelectOption) Help() string          { return o.help }
func (o *SelectOption) Path() path.Path       { return o.path }
func (o *SelectOption) Schema() schema.Schema { return o.schema }

func (o *SelectOption) UnmarshalJSON(data []byte) error {
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
		return serrors.New("missing option label property")
	} else if len(env.Label) > 100 {
		return serrors.New("too long option label field (max=100)")
	}

	// Help (optional, max=100)
	if len(env.Help) > 100 {
		return serrors.New("too long option help field (max=100)")
	}

	// Name (optional, max=100)
	if len(env.Name) > 100 {
		return serrors.New("too long option name field (max=100)")
	}

	o.label = env.Label
	o.help = env.Help
	o.name = env.Name

	if o.name == "" {
		o.name = slug.Make(o.label)
	}

	return nil
}
