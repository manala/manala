package yaml

import (
	goYamlAst "github.com/goccy/go-yaml/ast"
	"manala/internal/json"
	"manala/internal/schema"
	"manala/internal/schema/inferrer"
)

func NewNodeSchemaInferrer(node goYamlAst.Node) *NodeSchemaInferrer {
	return &NodeSchemaInferrer{
		node: node,
	}
}

type NodeSchemaInferrer struct {
	node   goYamlAst.Node
	schema schema.Schema
	err    error
}

func (inf *NodeSchemaInferrer) Infer(schema schema.Schema) error {
	if _, ok := any(inf.node).(goYamlAst.MapNode); !ok {
		return NewNodeError("unable to infer schema type", inf.node)
	}

	inf.schema = schema

	goYamlAst.Walk(inf, inf.node)

	return inf.err
}

func (inf *NodeSchemaInferrer) Visit(node goYamlAst.Node) goYamlAst.Visitor {
	schemaTags := &Tags{}

	// Get schema comment tags
	comment := node.GetComment()
	if comment != nil {
		var tags Tags
		ParseCommentTags(comment.String(), &tags)
		schemaTags = tags.Filter("schema")
	}

	if n, ok := node.(*goYamlAst.MappingValueNode); ok {
		// Get property key
		propertyKey := n.Key.GetToken().Value

		// Infer property schema
		propertySchema := schema.Schema{}
		if err := inferrer.NewChain(
			inferrer.NewFunc(func(schema schema.Schema) error {
				// Only mapping value
				if n, ok := node.(*goYamlAst.MappingValueNode); ok {
					if _, ok := n.Value.(goYamlAst.MapNode); ok {
						return NewNodeSchemaInferrer(n.Value).Infer(schema)
					}

					return nil
				}

				return NewNodeError("unable to infer schema type", node)
			}),
			NewNodeTagsSchemaInferrer(n, schemaTags),
			NewNodeTypeSchemaInferrer(n),
		).Infer(propertySchema); err != nil {
			inf.err = err

			return nil
		}

		// Ensure schema is set
		inf.schema["type"] = "object"
		inf.schema["additionalProperties"] = false
		if _, ok := inf.schema["properties"]; !ok {
			inf.schema["properties"] = map[string]any{}
		}

		// Set schema property
		inf.schema["properties"].(map[string]any)[propertyKey] = map[string]any(propertySchema)

		// Stop visiting when map nodes
		if _, ok := n.Value.(goYamlAst.MapNode); ok {
			return nil
		}
	} else {
		// Misplaced tag
		if len(*schemaTags) > 0 {
			inf.err = NewNodeError("misplaced schema tag", node.GetComment())
			return nil
		}
	}

	return inf
}

func NewNodeTypeSchemaInferrer(node goYamlAst.Node) *NodeTypeSchemaInferrer {
	return &NodeTypeSchemaInferrer{
		node: node,
	}
}

type NodeTypeSchemaInferrer struct {
	node goYamlAst.Node
}

func (inf *NodeTypeSchemaInferrer) Infer(schema schema.Schema) error {
	// Type already set, don't overwrite it
	if _, ok := schema["type"]; ok {
		return nil
	}

	// In case of an enum, don't infer the type
	if _, ok := schema["enum"]; ok {
		return nil
	}

	if n, ok := inf.node.(*goYamlAst.MappingValueNode); ok {
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
			return NewNodeError("unable to infer schema value type", v)
		}
	} else {
		return NewNodeError("unable to infer schema type", inf.node)
	}

	return nil
}

func NewNodeTagsSchemaInferrer(node goYamlAst.Node, tags *Tags) *NodeTagsSchemaInferrer {
	return &NodeTagsSchemaInferrer{
		node: node,
		tags: tags,
	}
}

type NodeTagsSchemaInferrer struct {
	node goYamlAst.Node
	tags *Tags
}

func (inf *NodeTagsSchemaInferrer) Infer(schema schema.Schema) error {
	for _, tag := range *inf.tags {
		if err := json.Unmarshal([]byte(tag.Value), &schema); err != nil {
			if inf.node != nil {
				return NewNodeError(err.Error(), inf.node.GetComment())
			}

			return err
		}
	}

	return nil
}
