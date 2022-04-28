package validator

import (
	"fmt"
	"github.com/goccy/go-yaml"
	"github.com/xeipuuv/gojsonschema"
	internalErrors "manala/internal/errors"
)

func ResultError(err gojsonschema.ResultError) *internalErrors.InternalError {
	_err := internalErrors.New(err.Description())

	field := err.Field()
	if field == "(root)" {
		field = ""
	}

	// Special error type treatments
	switch err.(type) {
	case *gojsonschema.InvalidTypeError:
		_ = _err.
			With("invalid type").
			WithField("expected", err.Details()["expected"]).
			WithField("given", err.Details()["given"])
	case *gojsonschema.AdditionalPropertyNotAllowedError:
		_ = _err.
			With("additional property is not allowed").
			WithField("property", err.Details()[gojsonschema.STRING_PROPERTY])
		if field != "" {
			field = fmt.Sprintf("%s.%s", field, err.Details()[gojsonschema.STRING_PROPERTY])
		}
	}

	if field != "" {
		_ = _err.WithField("field", field)
	}

	return _err
}

func ResultYamlContentError(err gojsonschema.ResultError, content []byte) *internalErrors.InternalError {
	_err := ResultError(err)

	if field := _err.Fields["field"]; field != "" {
		// Try to locate source
		pathString := fmt.Sprintf("$.%s", field)
		if path, __err := yaml.PathString(pathString); __err == nil {
			if source, __err := path.AnnotateSource(content, true); __err == nil {
				_ = _err.WithTrace(string(source))
			}
		}
	}

	return _err
}
