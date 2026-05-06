package manifest

import (
	"errors"
	"fmt"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/recipe/option"
	jsondecoder "github.com/manala/manala/internal/json/decoder"
	yamlannotation "github.com/manala/manala/internal/yaml/annotation"
	yamlerrors "github.com/manala/manala/internal/yaml/errors"

	"dario.cat/mergo"
	"github.com/goccy/go-yaml/ast"
)

type Inferrer struct {
	Schema  *map[string]any
	Options *[]app.RecipeOption
}

func (i *Inferrer) Infer(node ast.MapNode) error {
	sch, err := i.infer(node)
	if err != nil {
		return err
	}

	*i.Schema = sch

	return nil
}

func (i *Inferrer) infer(node ast.MapNode) (map[string]any, error) {
	// Init schema
	sch := map[string]any{
		"type":                 "object",
		"properties":           map[string]any{},
		"additionalProperties": false,
	}

	iter := node.MapRange()
	for iter.Next() {
		node := iter.KeyValue()

		// Start with empty property
		property := map[string]any{}

		// Infer type
		switch node := node.Value.(type) {
		case *ast.StringNode:
			property["type"] = "string"
		case *ast.IntegerNode:
			property["type"] = "integer"
		case *ast.FloatNode:
			property["type"] = "number"
		case *ast.BoolNode:
			property["type"] = "boolean"
		case ast.ArrayNode:
			property["type"] = "array"
		case ast.MapNode:
			var err error
			property, err = i.infer(node)
			if err != nil {
				return nil, err
			}
		case *ast.NullNode:
			// No type
		default:
			return nil, yamlerrors.New(
				errors.New("unable to infer schema value type"),
				node.GetToken(),
			)
		}

		// Parse comment annotations
		if comment := node.GetComment(); comment != nil {
			annotations, err := yamlannotation.Parse(comment.String())
			if err != nil {
				return nil, yamlerrors.New(err, comment.GetToken())
			}

			// Schema
			var propertySch map[string]any
			if err := annotations.JSONVar(&propertySch, "schema"); err != nil {
				return nil, yamlerrors.New(err, comment.GetToken())
			}

			if err := mergo.Merge(&property, propertySch, mergo.WithOverride); err != nil {
				return nil, yamlerrors.New(err, comment.GetToken())
			}

			// Enum makes type redundant
			if _, ok := property["enum"]; ok {
				delete(property, "type")
			}

			// Option
			if err := annotations.Func("option", func(a *yamlannotation.Annotation) error {
				// Stencil preserves source positions for accurate error reporting
				value := a.Value.Stencil()

				// Decode option type discriminator
				var disc struct {
					Type string `json:"type"`
				}
				if err := jsondecoder.Decode([]byte(value), &disc); err != nil {
					return err
				}

				// Resolve option type
				optionType := disc.Type
				if optionType == "" {
					if _, ok := property["enum"]; ok {
						optionType = option.ENUM
					} else if property["type"] == "string" {
						optionType = option.STRING
					} else {
						return yamlannotation.NewError(
							errors.New("unable to auto detect option type"),
							a.Value.Start(),
						)
					}
				}

				var opt app.RecipeOption
				switch optionType {
				case option.STRING:
					o, err := option.NewString(property, node.GetPath())
					if err != nil {
						return yamlannotation.NewError(err, a.Value.Start())
					}
					if err := o.UnmarshalJSON([]byte(value)); err != nil {
						return yamlannotation.NewError(err, a.Value.Start())
					}
					opt = o
				case option.ENUM:
					o, err := option.NewEnum(property, node.GetPath())
					if err != nil {
						return yamlannotation.NewError(err, a.Value.Start())
					}
					if err := o.UnmarshalJSON([]byte(value)); err != nil {
						return yamlannotation.NewError(err, a.Value.Start())
					}
					opt = o
				default:
					return yamlannotation.NewError(
						fmt.Errorf("unexpected \"%s\" option type", disc.Type),
						a.Value.Start(),
					)
				}

				*i.Options = append(*i.Options, opt)

				return nil
			}); err != nil {
				return nil, yamlerrors.New(err, comment.GetToken())
			}
		}

		sch["properties"].(map[string]any)[node.Key.String()] = property
	}

	return sch, nil
}
