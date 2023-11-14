package option

import (
	goYamlAst "github.com/goccy/go-yaml/ast"
	"manala/app"
	"manala/internal/schema"
	"manala/internal/schema/inferrer"
	"manala/internal/yaml"
	"strings"
)

func NewInferrer() *Inferrer {
	return &Inferrer{}
}

type Inferrer struct {
	options *[]app.RecipeOption
	err     error
}

func (inf *Inferrer) Infer(node goYamlAst.Node, options *[]app.RecipeOption) error {
	if _, ok := any(node).(goYamlAst.MapNode); !ok {
		return yaml.NewNodeError("unable to infer options type", node)
	}

	inf.options = options

	goYamlAst.Walk(inf, node)

	return inf.err
}

func (inf *Inferrer) Visit(node goYamlAst.Node) goYamlAst.Visitor {
	optionTags := &yaml.Tags{}
	schemaTags := &yaml.Tags{}

	// Get comment tags
	comment := node.GetComment()
	if comment != nil {
		var tags yaml.Tags
		yaml.ParseCommentTags(comment.String(), &tags)
		optionTags = tags.Filter("option")
		schemaTags = tags.Filter("schema")
	}

	if len(*optionTags) > 0 {
		if _node, ok := node.(*goYamlAst.MappingValueNode); ok {
			// Infer schema
			schema := schema.Schema{}
			if err := inferrer.NewChain(
				yaml.NewNodeTagsSchemaInferrer(_node, schemaTags),
				yaml.NewNodeTypeSchemaInferrer(_node),
			).Infer(schema); err != nil {
				inf.err = err

				return nil
			}

			// Handle option tags
			for _, tag := range *optionTags {
				// Read from tag
				option, err := New(strings.NewReader(tag.Value), schema, yaml.NewNodePath(_node))
				if err != nil {
					inf.err = yaml.NewNodeError("unable to read recipe option", _node.GetComment()).
						WithErrors(err)
					return nil
				}

				*inf.options = append(*inf.options, option)
			}
		} else {
			// Misplaced tag
			inf.err = yaml.NewNodeError("misplaced recipe option tag", node.GetComment())
			return nil
		}
	}

	return inf
}
