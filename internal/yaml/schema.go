package yaml

import (
	"encoding/json"
	yamlAst "github.com/goccy/go-yaml/ast"
	internalReport "manala/internal/report"
)

type schemaInferrerInterface interface {
	Infer(node yamlAst.Node, schema map[string]interface{}) error
}

func NewSchemaChainInferrer(inferrers ...schemaInferrerInterface) *SchemaChainInferrer {
	return &SchemaChainInferrer{
		inferrers: inferrers,
	}
}

type SchemaChainInferrer struct {
	inferrers []schemaInferrerInterface
}

func (inferrer *SchemaChainInferrer) Infer(node yamlAst.Node, schema map[string]interface{}) error {
	// Range over inferrers
	for _, inferrer := range inferrer.inferrers {
		// Inferrer schema
		if err := inferrer.Infer(node, schema); err != nil {
			return err
		}
	}

	return nil
}

func NewSchemaTypeInferrer() *SchemaTypeInferrer {
	return &SchemaTypeInferrer{}
}

type SchemaTypeInferrer struct{}

func (inferrer *SchemaTypeInferrer) Infer(node yamlAst.Node, schema map[string]interface{}) error {
	if n, ok := node.(*yamlAst.MappingValueNode); ok {
		// Infer schema type based on node value type
		switch v := n.Value.(type) {
		case *yamlAst.StringNode:
			schema["type"] = "string"
		case *yamlAst.IntegerNode:
			schema["type"] = "integer"
		case *yamlAst.FloatNode:
			schema["type"] = "number"
		case *yamlAst.BoolNode:
			schema["type"] = "boolean"
		case yamlAst.ArrayNode:
			schema["type"] = "array"
		case yamlAst.MapNode:
			schema["type"] = "object"
		case *yamlAst.NullNode:
			// No type
		default:
			return NewNodeError("unable to infer schema value type", v)
		}
	} else {
		return NewNodeError("unable to infer schema type", node)
	}

	return nil
}

func NewSchemaTagsInferrer(tags *Tags) *SchemaTagsInferrer {
	return &SchemaTagsInferrer{
		tags: tags,
	}
}

type SchemaTagsInferrer struct {
	tags *Tags
}

func (inferrer *SchemaTagsInferrer) Infer(node yamlAst.Node, schema map[string]interface{}) error {
	for _, tag := range *inferrer.tags {
		if err := json.Unmarshal([]byte(tag.Value), &schema); err != nil {
			if node != nil {
				return internalReport.NewError(NewNodeError(err.Error(), node)).
					WithMessage("unable to unmarshal json")
			}
			return internalReport.NewError(err).
				WithMessage("unable to unmarshal json")
		}
	}

	return nil
}

func NewSchemaCallbackInferrer(callback func(node yamlAst.Node, schema map[string]interface{}) error) *SchemaCallbackInferrer {
	return &SchemaCallbackInferrer{
		callback: callback,
	}
}

type SchemaCallbackInferrer struct {
	callback func(node yamlAst.Node, schema map[string]interface{}) error
}

func (inferrer *SchemaCallbackInferrer) Infer(node yamlAst.Node, schema map[string]interface{}) error {
	return inferrer.callback(node, schema)
}

func NewSchemaInferrer() *SchemaInferrer {
	return &SchemaInferrer{}
}

type SchemaInferrer struct {
	schema map[string]interface{}
	err    error
}

func (inferrer *SchemaInferrer) Infer(node yamlAst.Node, schema map[string]interface{}) error {
	if _, ok := interface{}(node).(yamlAst.MapNode); !ok {
		return NewNodeError("unable to infer schema type", node)
	}

	inferrer.schema = schema

	yamlAst.Walk(inferrer, node)

	return inferrer.err
}

func (inferrer *SchemaInferrer) Visit(node yamlAst.Node) yamlAst.Visitor {
	schemaTags := &Tags{}

	// Get schema comment tags
	comment := node.GetComment()
	if comment != nil {
		var tags Tags
		ParseCommentTags(comment.String(), &tags)
		schemaTags = tags.Filter("schema")
	}

	if n, ok := node.(*yamlAst.MappingValueNode); ok {
		// Get property key
		propertyKey := n.Key.GetToken().Value

		// Infer property schema
		propertySchema := map[string]interface{}{}
		if err := NewSchemaChainInferrer(
			NewSchemaTypeInferrer(),
			NewSchemaCallbackInferrer(func(node yamlAst.Node, schema map[string]interface{}) error {
				// Only mapping value
				if n, ok := node.(*yamlAst.MappingValueNode); ok {
					if _, ok := n.Value.(yamlAst.MapNode); ok {
						return NewSchemaInferrer().Infer(n.Value, schema)
					}

					return nil
				}

				return NewNodeError("unable to infer schema type", node)
			}),
			NewSchemaTagsInferrer(schemaTags),
		).Infer(n, propertySchema); err != nil {
			inferrer.err = err

			return nil
		}

		// Ensure schema is set
		inferrer.schema["type"] = "object"
		inferrer.schema["additionalProperties"] = false
		if _, ok := inferrer.schema["properties"]; !ok {
			inferrer.schema["properties"] = map[string]interface{}{}
		}

		// Set schema property
		inferrer.schema["properties"].(map[string]interface{})[propertyKey] = propertySchema

		// Stop visiting when map nodes
		if _, ok := n.Value.(yamlAst.MapNode); ok {
			return nil
		}
	} else {
		// Misplaced tag
		if len(*schemaTags) > 0 {
			inferrer.err = NewNodeError("misplaced schema tag", node)
			return nil
		}
	}

	return inferrer
}
