package yaml

import (
	goYamlAst "github.com/goccy/go-yaml/ast"
	"manala/internal/serrors"
)

func NewExtractor(node *goYamlAst.Node) *Extractor {
	return &Extractor{
		node: node,
	}
}

type Extractor struct {
	node *goYamlAst.Node
}

func (extractor *Extractor) ExtractRootMap(key string) (goYamlAst.Node, error) {
	var subject goYamlAst.Node

	switch node := (*extractor.node).(type) {
	case *goYamlAst.MappingValueNode:
		if node.Key.GetToken().Value == key {
			subject = node.Value
			*extractor.node = &goYamlAst.MappingNode{
				BaseNode: &goYamlAst.BaseNode{},
			}
		}
	case *goYamlAst.MappingNode:
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
		return nil, serrors.New("unable to find map").
			WithArguments("key", key)
	}

	if _, ok := subject.(goYamlAst.MapNode); !ok {
		return nil, NewNodeError("key is not a map", subject).
			WithArguments("key", key)
	}

	return subject, nil
}
