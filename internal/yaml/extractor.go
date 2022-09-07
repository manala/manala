package yaml

import (
	"fmt"
	yamlAst "github.com/goccy/go-yaml/ast"
)

func NewExtractor(node *yamlAst.Node) *Extractor {
	return &Extractor{
		node: node,
	}
}

type Extractor struct {
	node *yamlAst.Node
}

func (extractor *Extractor) ExtractRootMap(key string) (yamlAst.Node, error) {
	var subject yamlAst.Node

	switch node := (*extractor.node).(type) {
	case *yamlAst.MappingValueNode:
		if node.Key.GetToken().Value == key {
			subject = node.Value
			*extractor.node = &yamlAst.MappingNode{
				BaseNode: &yamlAst.BaseNode{},
			}
		}
	case *yamlAst.MappingNode:
		for i, n := range node.Values {
			if n.Key.GetToken().Value == key {
				subject = n.Value
				node.Values = append(node.Values[:i], node.Values[i+1:]...)
				if len(node.Values) == 1 {
					*extractor.node = node.Values[0]
				}
				break
			}
		}
	default:
		return nil, NewNodeError("root must be a map", node)
	}

	if subject == nil {
		return nil, fmt.Errorf("unable to find \"%s\" map", key)
	}

	if _, ok := subject.(yamlAst.MapNode); !ok {
		return nil, NewNodeError(
			fmt.Sprintf("\"%s\" is not a map", key),
			subject,
		)
	}

	return subject, nil
}
