package yaml

import (
	"github.com/manala/manala/internal/json"
	"github.com/manala/manala/internal/schema"
	"github.com/manala/manala/internal/schema/inferrer"
	"github.com/manala/manala/internal/yaml/annotation"

	goYamlAst "github.com/goccy/go-yaml/ast"
)

type NodeSchemaInferrer struct {
	node   goYamlAst.Node
	schema schema.Schema
	err    error
}

func NewNodeSchemaInferrer(node goYamlAst.Node) *NodeSchemaInferrer {
	return &NodeSchemaInferrer{
		node: node,
	}
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
	// Schema annotation
	var schemaAnnot *annotation.Annotation

	// Get comment
	comment := node.GetComment()
	if comment != nil {
		// Get annotations
		annots, err := annotation.Parse(comment.String())
		if err != nil {
			inf.err = err
			return nil
		}

		// Schema annotation
		schemaAnnot, _ = annots.Lookup("schema")
	}

	n, ok := node.(*goYamlAst.MappingValueNode)
	if !ok {
		if schemaAnnot != nil {
			// Misplaced annotation
			inf.err = NewNodeError("misplaced schema annotation", node.GetComment())
			return nil
		}
		return inf
	}

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
		NewNodeAnnotationSchemaInferrer(n, schemaAnnot),
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

	return inf
}

type NodeTypeSchemaInferrer struct {
	node goYamlAst.Node
}

func NewNodeTypeSchemaInferrer(node goYamlAst.Node) *NodeTypeSchemaInferrer {
	return &NodeTypeSchemaInferrer{
		node: node,
	}
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

type NodeAnnotationSchemaInferrer struct {
	node       goYamlAst.Node
	annotation *annotation.Annotation
}

func NewNodeAnnotationSchemaInferrer(node goYamlAst.Node, annot *annotation.Annotation) *NodeAnnotationSchemaInferrer {
	return &NodeAnnotationSchemaInferrer{
		node:       node,
		annotation: annot,
	}
}

func (inf *NodeAnnotationSchemaInferrer) Infer(schema schema.Schema) error {
	if inf.annotation == nil {
		return nil
	}

	err := json.Unmarshal([]byte(inf.annotation.Value()), &schema)
	if err != nil {
		if inf.node != nil {
			return NewNodeError(err.Error(), inf.node.GetComment())
		}
		return err
	}

	return nil
}
