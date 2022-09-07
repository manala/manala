package recipe

import (
	"encoding/json"
	"fmt"
	"github.com/goccy/go-yaml"
	yamlAst "github.com/goccy/go-yaml/ast"
	"manala/core"
	internalYaml "manala/internal/yaml"
)

type option struct {
	label  string
	schema map[string]interface{}
	node   *yamlAst.MappingValueNode
}

func (option *option) Label() string {
	return option.label
}

func (option *option) Schema() map[string]interface{} {
	return option.schema
}

func (option *option) Set(value interface{}) error {
	// Are float actually int ?
	// Coming from json, every number is a float...
	switch v := value.(type) {
	case float32:
		vv := uint32(v)
		if v == float32(vv) {
			value = vv
		}
	case float64:
		vv := uint64(v)
		if v == float64(vv) {
			value = vv
		}
	}

	switch v := value.(type) {
	case
		nil,
		bool,
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64,
		string:
		node, err := yaml.ValueToNode(v)
		if err != nil {
			return err
		}

		return option.node.Replace(node)
	}

	return fmt.Errorf("unsupported option value type: %s", value)
}

func (option *option) UnmarshalJSON(data []byte) error {
	var fields map[string]interface{}
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	if _, ok := fields["label"]; ok {
		option.label = fields["label"].(string)
	}

	return nil
}

func NewOptionsInferrer() *OptionsInferrer {
	return &OptionsInferrer{}
}

type OptionsInferrer struct {
	options *[]core.RecipeOption
	err     error
}

func (inferrer *OptionsInferrer) Infer(node yamlAst.Node, options *[]core.RecipeOption) error {
	if _, ok := interface{}(node).(yamlAst.MapNode); !ok {
		return internalYaml.NewNodeError("unable to infer options type", node)
	}

	inferrer.options = options

	yamlAst.Walk(inferrer, node)

	return inferrer.err
}

func (inferrer *OptionsInferrer) Visit(node yamlAst.Node) yamlAst.Visitor {
	optionTags := &internalYaml.Tags{}
	schemaTags := &internalYaml.Tags{}

	// Get schema comment tags
	comment := node.GetComment()
	if comment != nil {
		var tags internalYaml.Tags
		internalYaml.ParseCommentTags(comment.String(), &tags)
		optionTags = tags.Filter("option")
		schemaTags = tags.Filter("schema")
	}

	if len(*optionTags) > 0 {
		if n, ok := node.(*yamlAst.MappingValueNode); ok {
			// Infer schema
			schema := map[string]interface{}{}
			if err := internalYaml.NewSchemaChainInferrer(
				internalYaml.NewSchemaTypeInferrer(),
				internalYaml.NewSchemaTagsInferrer(schemaTags),
			).Infer(n, schema); err != nil {
				inferrer.err = err

				return nil
			}

			// Handle option tags
			for _, tag := range *optionTags {
				option := &option{
					schema: schema,
					node:   n,
				}
				if err := json.Unmarshal([]byte(tag.Value), &option); err != nil {
					inferrer.err = err

					return nil
				}

				*inferrer.options = append(*inferrer.options, option)
			}
		} else {
			// Misplaced tag
			inferrer.err = internalYaml.NewNodeError("misplaced option tag", node)
			return nil
		}
	}

	return inferrer
}
