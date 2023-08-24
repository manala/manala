package recipe

import (
	_ "embed"
	"encoding/json"
	"errors"
	goYaml "github.com/goccy/go-yaml"
	goYamlAst "github.com/goccy/go-yaml/ast"
	"github.com/xeipuuv/gojsonschema"
	"manala/app/interfaces"
	"manala/internal/errors/serrors"
	"manala/internal/validation"
	"manala/internal/yaml"
)

//go:embed resources/option.schema.json
var optionSchema string

type option struct {
	label  string
	schema map[string]interface{}
	node   *goYamlAst.MappingValueNode
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
		node, err := goYaml.ValueToNode(v)
		if err != nil {
			return err
		}

		return option.node.Replace(node)
	}

	return serrors.New("unsupported option value type").
		WithArguments("value", value)
}

func (option *option) UnmarshalJSON(data []byte) error {
	var fields map[string]interface{}
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	val, err := gojsonschema.Validate(
		gojsonschema.NewStringLoader(optionSchema),
		gojsonschema.NewGoLoader(fields),
	)
	if err != nil {
		return err
	}

	if !val.Valid() {
		return validation.NewError(
			"invalid option",
			val,
		).
			WithMessages([]validation.ErrorMessage{
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

func (inferrer *OptionsInferrer) Infer(node goYamlAst.Node, options *[]interfaces.RecipeOption) error {
	if _, ok := interface{}(node).(goYamlAst.MapNode); !ok {
		return yaml.NewNodeError("unable to infer options type", node)
	}

	inferrer.options = options

	goYamlAst.Walk(inferrer, node)

	return inferrer.err
}

func (inferrer *OptionsInferrer) Visit(node goYamlAst.Node) goYamlAst.Visitor {
	optionTags := &yaml.Tags{}
	schemaTags := &yaml.Tags{}

	// Get schema comment tags
	comment := node.GetComment()
	if comment != nil {
		var tags yaml.Tags
		yaml.ParseCommentTags(comment.String(), &tags)
		optionTags = tags.Filter("option")
		schemaTags = tags.Filter("schema")
	}

	if len(*optionTags) > 0 {
		if n, ok := node.(*goYamlAst.MappingValueNode); ok {
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
					var _validationError *validation.Error
					switch {
					case errors.As(err, &_validationError):
						inferrer.err = _validationError.
							WithResultErrorDecorator(yaml.NewNodeValidationResultErrorDecorator(node.GetComment()))
					default:
						inferrer.err = yaml.NewNodeError(err.Error(), node.GetComment())
					}

					return nil
				}

				*inferrer.options = append(*inferrer.options, option)
			}
		} else {
			// Misplaced tag
			inferrer.err = yaml.NewNodeError("misplaced option tag", node.GetComment())
			return nil
		}
	}

	return inferrer
}
