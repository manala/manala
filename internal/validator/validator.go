package validator

import (
	"github.com/xeipuuv/gojsonschema"
	internalErrors "manala/internal/errors"
)

func ValidateFormat(format string, document interface{}) bool {
	result, _ := gojsonschema.Validate(
		gojsonschema.NewGoLoader(map[string]interface{}{"type": "string", "format": format}),
		gojsonschema.NewGoLoader(document),
	)

	return result.Valid()
}

func Validate(schema interface{}, document interface{}, options ...ValidateOption) (error, []*internalErrors.InternalError, bool) {
	var schemaLoader gojsonschema.JSONLoader
	switch _schema := schema.(type) {
	case string:
		schemaLoader = gojsonschema.NewStringLoader(_schema)
	default:
		schemaLoader = gojsonschema.NewGoLoader(schema)
	}

	ok := true

	result, err := gojsonschema.Validate(
		schemaLoader,
		gojsonschema.NewGoLoader(document),
	)

	if err != nil || !result.Valid() {
		ok = false
	}

	// Parameters
	parameters := &validateParameters{}
	for _, option := range options {
		option(parameters)
	}

	var errs []*internalErrors.InternalError
	if result != nil {
		for _, err := range result.Errors() {
			if parameters.yamlContent != nil {
				errs = append(errs, ResultYamlContentError(err, parameters.yamlContent))
			} else {
				errs = append(errs, ResultError(err))
			}
		}
	}

	return err, errs, ok
}

type validateParameters struct {
	yamlContent []byte
}

type ValidateOption func(parameters *validateParameters)

func WithYamlContent(content []byte) ValidateOption {
	return func(parameters *validateParameters) {
		parameters.yamlContent = content
	}
}

func init() {
	gojsonschema.FormatCheckers.
		Add("git-repo", gitRepoFormatChecker{}).
		Add("file-path", filePathFormatChecker{}).
		Add("domain", domainFormatChecker{})
}
