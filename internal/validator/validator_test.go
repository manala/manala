package validator

import (
	"github.com/stretchr/testify/suite"
	"io"
	"testing"
)

type ValidatorSuite struct{ suite.Suite }

func TestValidatorSuite(t *testing.T) {
	suite.Run(t, new(ValidatorSuite))
}

func (s *ValidatorSuite) TestValidateFormat() {
	s.True(ValidateFormat("git-repo", "https://github.com/manala/manala-recipes.git"))
	s.True(ValidateFormat("git-repo", "git@github.com:manala/manala.git"))

	s.False(ValidateFormat("git-repo", "foo"))
}

func (s *ValidatorSuite) TestValidate() {
	s.Run("Empty Schema", func() {
		err, errs, ok := Validate("", map[string]interface{}{})
		s.EqualError(err, io.EOF.Error())
		s.Empty(errs)
		s.False(ok)
	})

	s.Run("Invalid Schema", func() {
		err, errs, ok := Validate("foo", map[string]interface{}{})
		s.EqualError(err, "invalid character 'o' in literal false (expecting 'a')")
		s.Empty(errs)
		s.False(ok)
	})

	s.Run("String Schema", func() {
		schema := `{
		  "type": "object",
		  "properties": {
		    "foo": {"type": "string"}
		  }
		}`
		document := map[string]interface{}{
			"foo": "bar",
		}
		err, errs, ok := Validate(schema, document)
		s.NoError(err)
		s.Empty(errs)
		s.True(ok)
	})

	s.Run("Go Schema", func() {
		schema := map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"foo": map[string]interface{}{"type": "string"},
			},
		}
		document := map[string]interface{}{
			"foo": "bar",
		}
		err, errs, ok := Validate(schema, document)
		s.NoError(err)
		s.Empty(errs)
		s.True(ok)
	})

	s.Run("Invalid Document", func() {
		schema := `{
		  "type": "object",
		  "properties": {
		    "string": {"type": "string"},
			"integer": {"type": "integer"},
			"boolean": {"type": "boolean"}
		  }
		}`
		document := map[string]interface{}{
			"string":  "string",
			"integer": "123",
			"boolean": true,
		}
		err, errs, ok := Validate(schema, document)
		s.NoError(err)
		s.Len(errs, 1)

		s.ErrorAs(errs[0], &internalError)
		s.Equal("invalid type", internalError.Message)
		s.Equal("integer", internalError.Fields["expected"])
		s.Equal("string", internalError.Fields["given"])
		s.Equal("integer", internalError.Fields["field"])
		s.False(ok)
	})

	s.Run("Invalid Document With Yaml Content", func() {
		schema := map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"string":  map[string]interface{}{"type": "string"},
				"integer": map[string]interface{}{"type": "integer"},
				"boolean": map[string]interface{}{"type": "boolean"},
			},
		}
		document := map[string]interface{}{
			"string":  "string",
			"integer": "123",
			"boolean": true,
		}
		content := `string: string
integer: "123"
boolean: true`
		err, errs, ok := Validate(schema, document, WithYamlContent([]byte(content)))
		s.NoError(err)
		s.Len(errs, 1)

		s.ErrorAs(errs[0], &internalError)
		s.Equal("invalid type", internalError.Message)
		s.Equal("integer", internalError.Fields["expected"])
		s.Equal("string", internalError.Fields["given"])
		s.Equal("integer", internalError.Fields["field"])
		s.Equal("   1 | \x1b[96mstring\x1b[0m:\x1b[92m string\x1b[0m\n>  2 | \x1b[92m\x1b[0m\x1b[96minteger\x1b[0m:\x1b[92m \"123\"\x1b[0m\n                ^\n   3 | \x1b[96mboolean\x1b[0m:\x1b[95m true\x1b[0m", internalError.Trace)
		s.False(ok)
	})
}
