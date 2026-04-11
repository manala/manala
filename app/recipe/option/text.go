package option

import (
	"github.com/manala/manala/internal/json/number"
	"github.com/manala/manala/internal/json/unmarshaler"
	"github.com/manala/manala/internal/path"
	"github.com/manala/manala/internal/schema"
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/validator"

	"github.com/gosimple/slug"
)

type TextOption struct {
	name   string
	label  string
	help   string
	schema schema.Schema
	path   path.Path

	MaxLength int
}

func NewTextOption(sch schema.Schema, p path.Path) (*TextOption, error) {
	// Schema type *MUST* be string
	if t, ok := sch["type"]; !ok || t != "string" {
		return nil, serrors.New("invalid recipe option string type")
	}

	// Max length
	maxLength := 0
	if length, ok := number.NumberType(sch["maxLength"]); ok {
		maxLength = length.Int()
	}

	return &TextOption{
		schema:    sch,
		path:      p,
		MaxLength: maxLength,
	}, nil
}

func (o *TextOption) Name() string          { return o.name }
func (o *TextOption) Label() string         { return o.label }
func (o *TextOption) Help() string          { return o.help }
func (o *TextOption) Path() path.Path       { return o.path }
func (o *TextOption) Schema() schema.Schema { return o.schema }

func (o *TextOption) Validate(_ any) (validator.Violations, error) {
	return nil, nil
}

func (o *TextOption) UnmarshalJSON(data []byte) error {
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
