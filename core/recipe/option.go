package recipe

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/goccy/go-yaml"
	yamlAst "github.com/goccy/go-yaml/ast"
	"github.com/xeipuuv/gojsonschema"
	"manala/app/interfaces"
	internalValidation "manala/internal/validation"
	internalYaml "manala/internal/yaml"
)

//go:embed resources/option.schema.json
var optionSchema string

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

	validation, err := gojsonschema.Validate(
		gojsonschema.NewStringLoader(optionSchema),
		gojsonschema.NewGoLoader(fields),
	)
	if err != nil {
		return err
	}

	if !validation.Valid() {
		return internalValidation.NewError(
			"invalid option",
			validation,
		).
			WithMessages([]internalValidation.ErrorMessage{
				{Field: "(root)", Type: "required", Message: "missing label field"},
				{Field: "(root)", Type: "additional_property_not_allowed", Message: "don't support additional properties"},
			})
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
	options *[]interfaces.RecipeOption
	err     error
}

func (inferrer *OptionsInferrer) Infer(node yamlAst.Node, options *[]interfaces.RecipeOption) error {
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
			if err := NewSchemaChainInferrer(
				NewSchemaTypeInferrer(),
				NewSchemaTagsInferrer(schemaTags),
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

				// Unmarshall
				if err := json.Unmarshal([]byte(tag.Value), &option); err != nil {
					var _validationError *internalValidation.Error
					switch {
					case errors.As(err, &_validationError):
						inferrer.err = _validationError.
							WithReporter(internalYaml.NewValidationReporter(node.GetComment()))
					default:
						inferrer.err = internalYaml.NewNodeError(err.Error(), node.GetComment())
					}

					return nil
				}

				*inferrer.options = append(*inferrer.options, option)
			}
		} else {
			// Misplaced tag
			inferrer.err = internalYaml.NewNodeError("misplaced option tag", node.GetComment())
			return nil
		}
	}

	return inferrer
}
