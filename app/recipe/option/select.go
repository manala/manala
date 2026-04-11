package option

import (
	"github.com/manala/manala/internal/accessor"
	"github.com/manala/manala/internal/json/number"
	"github.com/manala/manala/internal/json/unmarshaler"
	"github.com/manala/manala/internal/path"
	"github.com/manala/manala/internal/schema"
	"github.com/manala/manala/internal/serrors"

	"github.com/gosimple/slug"
)

type Select struct {
	name   string
	label  string
	help   string
	schema schema.Schema
	path   path.Path
}

func NewSelect(sch schema.Schema, p path.Path) (*Select, error) {
	// Enum
	enum, ok := sch["enum"].([]any)
	if !ok {
		return nil, serrors.New("invalid recipe option enum")
	}

	if len(enum) == 0 {
		return nil, serrors.New("empty recipe option enum")
	}

	return &Select{
		schema: sch,
		path:   p,
	}, nil
}

func (o *Select) Name() string  { return o.name }
func (o *Select) Label() string { return o.label }
func (o *Select) Help() string  { return o.help }

func (o *Select) Values() []any {
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

func (o *Select) Accessor(data any) accessor.Accessor {
	return path.NewAccessor(o.path, data)
}

func (o *Select) UnmarshalJSON(data []byte) error {
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
