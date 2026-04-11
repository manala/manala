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

type Text struct {
	name   string
	label  string
	help   string
	schema schema.Schema
	path   path.Path
}

func NewText(sch schema.Schema, p path.Path) (*Text, error) {
	// Schema type *MUST* be string
	if t, ok := sch["type"]; !ok || t != "string" {
		return nil, serrors.New("invalid recipe option string type")
	}

	return &Text{
		schema: sch,
		path:   p,
	}, nil
}

func (o *Text) Name() string  { return o.name }
func (o *Text) Label() string { return o.label }
func (o *Text) Help() string  { return o.help }

func (o *Text) MaxLength() int {
	if maxLength, ok := number.NumberType(o.schema["maxLength"]); ok {
		return maxLength.Int()
	}
	return 0
}

func (o *Text) Accessor(data any) accessor.Accessor {
	return path.NewAccessor(o.path, data)
}

func (o *Text) Validator() *schema.Validator {
	return schema.NewValidator(o.schema)
}

func (o *Text) UnmarshalJSON(data []byte) error {
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
