package yaml

import (
	"strings"

	"github.com/manala/manala/internal/path"
	"github.com/manala/manala/internal/serrors"

	goYaml "github.com/goccy/go-yaml"
	goYamlAst "github.com/goccy/go-yaml/ast"
)

func NewNodePath(node goYamlAst.Node) path.Path {
	nodePath := node.GetPath()

	if nodePath == "$" {
		nodePath = ""
	} else {
		nodePath, _ = strings.CutPrefix(nodePath, "$.")
	}

	return path.Path(nodePath)
}

type NodePathAccessor struct {
	path path.Path
}

func NewNodePathAccessor(path path.Path) NodePathAccessor {
	return NodePathAccessor{
		path: path,
	}
}

func (accessor NodePathAccessor) Get(node goYamlAst.Node) (goYamlAst.Node, error) {
	path := accessor.path.String()

	// Compute path
	if path == "" {
		path = "$"
	} else {
		path = "$." + path
	}

	// Get yaml path
	yamlPath, err := goYaml.PathString(path)
	if err != nil {
		return nil, err
	}

	// Filter node
	node, err = yamlPath.FilterNode(node)
	if err != nil {
		return nil, err
	}

	if node == nil {
		return nil, serrors.New("unable to access yaml path").
			WithArguments("path", accessor.path.String())
	}

	return node, nil
}
