package recipe

import (
	"encoding/json"
	goYamlAst "github.com/goccy/go-yaml/ast"
	"manala/internal/yaml"
)

type schemaInferrerInterface interface {
	Infer(node goYamlAst.Node, schema map[string]interface{}) error
}

func NewSchemaChainInferrer(inferrers ...schemaInferrerInterface) *SchemaChainInferrer {
	return &SchemaChainInferrer{
		inferrers: inferrers,
	}
}

type SchemaChainInferrer struct {
	inferrers []schemaInferrerInterface
}

func (inferrer *SchemaChainInferrer) Infer(node goYamlAst.Node, schema map[string]interface{}) error {
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

func (inferrer *SchemaTypeInferrer) Infer(node goYamlAst.Node, schema map[string]interface{}) error {
	if n, ok := node.(*goYamlAst.MappingValueNode); ok {
		// Infer schema type based on node value type
		switch v := n.Value.(type) {
		case *goYamlAst.StringNode:
			schema["type"] = "string"
		case *goYamlAst.IntegerNode:
			schema["type"] = "integer"
		case *goYamlAst.FloatNode:
			schema["type"] = "number"
		case *goYamlAst.BoolNode:
			schema["type"] = "boolean"
		case goYamlAst.ArrayNode:
			schema["type"] = "array"
		case goYamlAst.MapNode:
			schema["type"] = "object"
		case *goYamlAst.NullNode:
			// No type
		default:
			return yaml.NewNodeError("unable to infer schema value type", v)
		}
	} else {
		return yaml.NewNodeError("unable to infer schema type", node)
	}

	return nil
}

func NewSchemaTagsInferrer(tags *yaml.Tags) *SchemaTagsInferrer {
	return &SchemaTagsInferrer{
		tags: tags,
	}
}

type SchemaTagsInferrer struct {
	tags *yaml.Tags
}

func (inferrer *SchemaTagsInferrer) Infer(node goYamlAst.Node, schema map[string]interface{}) error {
	for _, tag := range *inferrer.tags {
		if err := json.Unmarshal([]byte(tag.Value), &schema); err != nil {
			if node != nil {
				return yaml.NewNodeError(err.Error(), node.GetComment())
			}

			return err
		}
	}

	return nil
}

func NewSchemaCallbackInferrer(callback func(node goYamlAst.Node, schema map[string]interface{}) error) *SchemaCallbackInferrer {
	return &SchemaCallbackInferrer{
		callback: callback,
	}
}

type SchemaCallbackInferrer struct {
	callback func(node goYamlAst.Node, schema map[string]interface{}) error
}

func (inferrer *SchemaCallbackInferrer) Infer(node goYamlAst.Node, schema map[string]interface{}) error {
	return inferrer.callback(node, schema)
}

func NewSchemaInferrer() *SchemaInferrer {
	return &SchemaInferrer{}
}

type SchemaInferrer struct {
	schema map[string]interface{}
	err    error
}

func (inferrer *SchemaInferrer) Infer(node goYamlAst.Node, schema map[string]interface{}) error {
	if _, ok := interface{}(node).(goYamlAst.MapNode); !ok {
		return yaml.NewNodeError("unable to infer schema type", node)
	}

	inferrer.schema = schema

	goYamlAst.Walk(inferrer, node)

	return inferrer.err
}

func (inferrer *SchemaInferrer) Visit(node goYamlAst.Node) goYamlAst.Visitor {
	schemaTags := &yaml.Tags{}

	// Get schema comment tags
	comment := node.GetComment()
	if comment != nil {
		var tags yaml.Tags
		yaml.ParseCommentTags(comment.String(), &tags)
		schemaTags = tags.Filter("schema")
	}

	if n, ok := node.(*goYamlAst.MappingValueNode); ok {
		// Get property key
		propertyKey := n.Key.GetToken().Value

		// Infer property schema
		propertySchema := map[string]interface{}{}
		if err := NewSchemaChainInferrer(
			NewSchemaTypeInferrer(),
			NewSchemaCallbackInferrer(func(node goYamlAst.Node, schema map[string]interface{}) error {
				// Only mapping value
				if n, ok := node.(*goYamlAst.MappingValueNode); ok {
					if _, ok := n.Value.(goYamlAst.MapNode); ok {
						return NewSchemaInferrer().Infer(n.Value, schema)
					}

					return nil
				}

				return yaml.NewNodeError("unable to infer schema type", node)
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
		if _, ok := n.Value.(goYamlAst.MapNode); ok {
			return nil
		}
	} else {
		// Misplaced tag
		if len(*schemaTags) > 0 {
			inferrer.err = yaml.NewNodeError("misplaced schema tag", node.GetComment())
			return nil
		}
	}

	return inferrer
}
