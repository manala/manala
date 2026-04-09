package option

import (
	"strings"

	"github.com/manala/manala/app"
	"github.com/manala/manala/internal/schema"
	"github.com/manala/manala/internal/schema/inferrer"
	"github.com/manala/manala/internal/yaml"
	"github.com/manala/manala/internal/yaml/annotation"

	"github.com/goccy/go-yaml/ast"
)

type Inferrer struct {
	options *[]app.RecipeOption
	err     error
}

func NewInferrer() *Inferrer {
	return &Inferrer{}
}

func (inf *Inferrer) Infer(node ast.Node, options *[]app.RecipeOption) error {
	if _, ok := any(node).(ast.MapNode); !ok {
		return yaml.NewNodeError("unable to infer options type", node)
	}

	inf.options = options

	ast.Walk(inf, node)

	return inf.err
}

func (inf *Inferrer) Visit(node ast.Node) ast.Visitor {
	// Get comment
	comment := node.GetComment()
	if comment == nil {
		return inf
	}

	// Get annotations
	annots, err := annotation.Parse(comment.String())
	if err != nil {
		inf.err = err
		return nil
	}

	// Option annotation
	optionAnnot, ok := annots.Lookup("option")
	if !ok {
		return inf
	}

	// Misplaced annotation
	n, ok := node.(*ast.MappingValueNode)
	if !ok {
		inf.err = yaml.NewNodeError("misplaced recipe option annotation", node.GetComment())
		return nil
	}

	// Schema annotation
	schemaAnnot, _ := annots.Lookup("schema")

	// Infer schema
	schema := schema.Schema{}
	if err := inferrer.NewChain(
		yaml.NewNodeAnnotationSchemaInferrer(n, schemaAnnot),
		yaml.NewNodeTypeSchemaInferrer(n),
	).Infer(schema); err != nil {
		inf.err = err
		return nil
	}

	// Read from annotation
	option, err := New(strings.NewReader(optionAnnot.Value()), schema, yaml.NewNodePath(n))
	if err != nil {
		inf.err = yaml.NewNodeError("unable to read recipe option", n.GetComment()).
			WithErrors(err)
		return nil
	}

	*inf.options = append(*inf.options, option)

	return inf
}
