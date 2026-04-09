package yaml

import (
	"strings"

	goYamlAst "github.com/goccy/go-yaml/ast"
	"github.com/manala/manala/internal/path"
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
