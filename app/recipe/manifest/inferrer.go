package manifest

import (
	"errors"
	"fmt"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/recipe/option"
	jsondecoder "github.com/manala/manala/internal/json/decoder"
	jsonvalidation "github.com/manala/manala/internal/json/validation"
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
			annotations := yamlannotation.NewSet()

			// Schema
			annotations.BodyFunc("schema", func(body *yamlannotation.Body) error {
				// Decode property schema
				var propertySch map[string]any
				if err := jsondecoder.Decode([]byte(body.Stencil()), &propertySch); err != nil {
					return err
				}

				// Merge in property
				if err := mergo.Merge(&property, propertySch, mergo.WithOverride); err != nil {
					return err
				}

				// Enum makes type redundant
				if _, ok := property["enum"]; ok {
					delete(property, "type")
				}

				return nil
			})

			// Option
			annotations.BodyFunc("option", func(body *yamlannotation.Body) error {
				value := []byte(body.Stencil())

				// Decode option value
				var optionValue map[string]any
				if err := jsondecoder.Decode(value, &optionValue); err != nil {
					return err
				}

				// Validate option value
				if err := option.Validator.Validate(optionValue, jsonvalidation.WithLocator(value)); err != nil {
					return err
				}

				// Auto-detect type if not provided
				if _, ok := optionValue["type"]; !ok {
					if _, ok := property["enum"]; ok {
						optionValue["type"] = option.ENUM
					} else if property["type"] == "string" {
						optionValue["type"] = option.STRING
					} else {
						return yamlannotation.NewError(
							errors.New("unable to auto-detect option type"),
							body.Start(),
						)
					}
				}

				var opt app.RecipeOption
				switch optionValue["type"] {
				case option.STRING:
					o, err := option.NewString(property, node.GetPath())
					if err != nil {
						return yamlannotation.NewError(err, body.Start())
					}
					if err := o.UnmarshalJSON(value); err != nil {
						return err
					}
					opt = o
				case option.ENUM:
					o, err := option.NewEnum(property, node.GetPath())
					if err != nil {
						return yamlannotation.NewError(err, body.Start())
					}
					if err := o.UnmarshalJSON(value); err != nil {
						return err
					}
					opt = o
				default:
					return yamlannotation.NewError(
						fmt.Errorf("unexpected \"%s\" option type", optionValue["type"]),
						body.Start(),
					)
				}

				*i.Options = append(*i.Options, opt)

				return nil
			})

			if err := annotations.Parse(comment.String()); err != nil {
				return nil, yamlerrors.New(err, comment.GetToken())
			}
		}

		sch["properties"].(map[string]any)[node.Key.String()] = property
	}

	return sch, nil
}
