package validator

import (
	"github.com/stretchr/testify/suite"
	"github.com/xeipuuv/gojsonschema"
	internalErrors "manala/internal/errors"
	"testing"
)

var internalError *internalErrors.InternalError

type ErrorsSuite struct{ suite.Suite }

func TestErrorsSuite(t *testing.T) {
	suite.Run(t, new(ErrorsSuite))
}

func (s *ErrorsSuite) TestResultError() {
	s.Run("Root Invalid Type", func() {
		result, _ := gojsonschema.Validate(
			gojsonschema.NewStringLoader(`{"type": "string"}`),
			gojsonschema.NewStringLoader(`123`),
		)
		err := ResultError(result.Errors()[0])

		s.ErrorAs(err, &internalError)
		s.Equal("invalid type", internalError.Message)
		s.Equal("string", internalError.Fields["expected"])
		s.Equal("integer", internalError.Fields["given"])
		s.NotContains(internalError.Fields, "field")
		s.Empty(internalError.Trace)
	})

	s.Run("Invalid Type", func() {
		result, _ := gojsonschema.Validate(
			gojsonschema.NewStringLoader(`{
	"type": "object",
	"properties": {
		"string": {"type": "string"}
	}
}`),
			gojsonschema.NewStringLoader(`{"string": 123}`),
		)
		err := ResultError(result.Errors()[0])

		s.ErrorAs(err, &internalError)
		s.Equal("invalid type", internalError.Message)
		s.Equal("string", internalError.Fields["expected"])
		s.Equal("integer", internalError.Fields["given"])
		s.Equal("string", internalError.Fields["field"])
		s.Empty(internalError.Trace)
	})
}

func (s *ErrorsSuite) TestResultYamlContentError() {
	s.Run("Root Invalid Type", func() {
		result, _ := gojsonschema.Validate(
			gojsonschema.NewStringLoader(`{"type": "string"}`),
			gojsonschema.NewStringLoader(`123`),
		)
		err := ResultYamlContentError(result.Errors()[0], []byte(`123`))

		s.ErrorAs(err, &internalError)
		s.Equal("invalid type", internalError.Message)
		s.Equal("string", internalError.Fields["expected"])
		s.Equal("integer", internalError.Fields["given"])
		s.NotContains(internalError.Fields, "field")
		s.Empty(internalError.Trace)
	})

	s.Run("Invalid Type", func() {
		result, _ := gojsonschema.Validate(
			gojsonschema.NewStringLoader(`{
	"type": "object",
	"properties": {
		"string": {"type": "string"}
	}
}`),
			gojsonschema.NewStringLoader(`{"string": 123}`),
		)
		err := ResultYamlContentError(result.Errors()[0], []byte(`string: 123`))

		s.ErrorAs(err, &internalError)
		s.Equal("invalid type", internalError.Message)
		s.Equal("string", internalError.Fields["expected"])
		s.Equal("integer", internalError.Fields["given"])
		s.Equal("string", internalError.Fields["field"])
		s.Equal(">  1 | \x1b[96mstring\x1b[0m:\x1b[95m 123\x1b[0m\n               ^\n", internalError.Trace)
	})

	s.Run("Root Additional Property Not Allowed", func() {
		result, _ := gojsonschema.Validate(
			gojsonschema.NewStringLoader(`{
	"type": "object",
	"properties": {
		"string": {"type": "string"}
	},
	"additionalProperties": false
}`),
			gojsonschema.NewStringLoader(`{"string": "string", "integer": 123}`),
		)
		err := ResultYamlContentError(result.Errors()[0], []byte(`string: string
integer: 123`))

		s.ErrorAs(err, &internalError)
		s.Equal("additional property is not allowed", internalError.Message)
		s.Equal("integer", internalError.Fields["property"])
		s.NotContains(internalError.Fields, "field")
		s.Empty(internalError.Trace)
	})

	s.Run("Additional Property Not Allowed", func() {
		result, _ := gojsonschema.Validate(
			gojsonschema.NewStringLoader(`{
	"type": "object",
	"properties": {
		"object": {
			"type": "object",
			"properties": {
				"string": {
					"type": "string"
				}
			},
			"additionalProperties": false
		}
	}
}`),
			gojsonschema.NewStringLoader(`{"object": {
	"string": "string",
	"integer": 123
}}`),
		)
		err := ResultYamlContentError(result.Errors()[0], []byte(`object:
    string: string
    integer: 123`))

		s.ErrorAs(err, &internalError)
		s.Equal("additional property is not allowed", internalError.Message)
		s.Equal("integer", internalError.Fields["property"])
		s.Equal("object.integer", internalError.Fields["field"])
		s.Equal("   1 | \x1b[96mobject\x1b[0m:\x1b[96m\x1b[0m\n   2 | \x1b[96m    string\x1b[0m:\x1b[92m string\x1b[0m\n>  3 | \x1b[92m    \x1b[0m\x1b[96minteger\x1b[0m:\x1b[95m 123\x1b[0m\n                    ^\n", internalError.Trace)
	})
}
