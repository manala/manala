package path

import (
	"strings"

	"github.com/manala/manala/internal/path"

	"github.com/goccy/go-yaml/ast"
)

func NewNodePath(node ast.Node) path.Path {
	nodePath := node.GetPath()

	if nodePath == "$" {
		nodePath = ""
	} else {
		nodePath, _ = strings.CutPrefix(nodePath, "$.")
	}

	return path.Path(nodePath)
}
