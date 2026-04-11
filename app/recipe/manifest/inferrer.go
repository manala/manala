package manifest

import (
	"strings"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/recipe/option"
	"github.com/manala/manala/internal/schema"
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/yaml"
	"github.com/manala/manala/internal/yaml/annotation"
	"github.com/manala/manala/internal/yaml/parser"

	"dario.cat/mergo"
	"github.com/goccy/go-yaml/ast"
)

type Inferrer struct {
	Schema  *schema.Schema
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
			return nil, parser.ErrorAt(
				serrors.New("unable to infer schema value type"),
				node.GetToken(),
			)
		}

		// Parse comment annotations
		if comment := node.GetComment(); comment != nil {
			annotations, err := annotation.Parse(comment.String())
			if err != nil {
				return nil, parser.ErrorAt(err, comment.GetToken())
			}

			// Schema
			var propertySch map[string]any
			if err := annotations.JSONVar(&propertySch, "schema"); err != nil {
				return nil, parser.ErrorAt(err, comment.GetToken())
			}

			if err := mergo.Merge(&property, propertySch, mergo.WithOverride); err != nil {
				return nil, parser.ErrorAt(err, comment.GetToken())
			}

			// Enum makes type redundant
			if _, ok := property["enum"]; ok {
				delete(property, "type")
			}

			// Option
			if err := annotations.Func("option", func(a *annotation.Annotation) error {
				opt, err := option.New(strings.NewReader(a.Value.String()), property, yaml.NewNodePath(node))
				if err != nil {
					return annotation.ErrorAt(err, a.Value.Tokens[0])
				}
				*i.Options = append(*i.Options, opt)
				return nil
			}); err != nil {
				return nil, parser.ErrorAt(err, comment.GetToken())
			}
		}

		sch["properties"].(map[string]any)[node.Key.String()] = property
	}

	return sch, nil
}
