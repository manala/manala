package yaml

import (
	"fmt"
	goYaml "github.com/goccy/go-yaml"
	goYamlAst "github.com/goccy/go-yaml/ast"
	"github.com/xeipuuv/gojsonschema"
	"io"
	"manala/internal/errors/serrors"
	"manala/internal/validation"
	"regexp"
)

/**********/
/* Loader */
/**********/

func NewJsonLoader(node goYamlAst.Node) gojsonschema.JSONLoader {
	return &jsonLoader{
		JSONLoader: gojsonschema.NewRawLoader(node),
	}
}

type jsonLoader struct {
	gojsonschema.JSONLoader
}

func (loader *jsonLoader) LoadJSON() (interface{}, error) {
	var data interface{}

	// Decode node into data
	if err := goYaml.NewDecoder(loader.JsonSource().(goYamlAst.Node)).Decode(&data); err != nil {
		// Nil or empty content
		if err == io.EOF {
			return nil, fmt.Errorf("empty content")
		}
		return nil, NewError(err)
	}

	return gojsonschema.NewGoLoader(data).LoadJSON()
}

/**************/
/* Decorators */
/**************/

func NewNodeValidationResultErrorDecorator(node goYamlAst.Node) *NodeValidationResultErrorDecorator {
	return &NodeValidationResultErrorDecorator{node: node}
}

type NodeValidationResultErrorDecorator struct{ node goYamlAst.Node }

func (decorator *NodeValidationResultErrorDecorator) Decorate(err validation.ResultErrorInterface) validation.ResultErrorInterface {
	return &NodeValidationResultError{
		resultErr: err,
		nodeErr:   NewNodeError(err.Error(), decorator.node),
		Arguments: serrors.NewArguments(),
		Details:   serrors.NewDetails(),
	}
}

func NewNodeValidationResultPathErrorDecorator(node goYamlAst.Node) *NodeValidationResultPathErrorDecorator {
	return &NodeValidationResultPathErrorDecorator{node: node}
}

type NodeValidationResultPathErrorDecorator struct{ node goYamlAst.Node }

func (decorator *NodeValidationResultPathErrorDecorator) Decorate(err validation.ResultErrorInterface) validation.ResultErrorInterface {
	var node goYamlAst.Node

	// Normalize result path
	path := decorator.normalizePath(err.Path())

	// Get yaml path
	if yamlPath, _err := goYaml.PathString(path); _err == nil {
		// Get node
		node, _ = yamlPath.FilterNode(decorator.node)
	}

	return &NodeValidationResultError{
		resultErr: err,
		nodeErr:   NewNodeError(err.Error(), node),
		Arguments: serrors.NewArguments(),
		Details:   serrors.NewDetails(),
	}
}

var resultPathNormalizeRegex = regexp.MustCompile(`\.(\d+)`)

func (decorator *NodeValidationResultPathErrorDecorator) normalizePath(path string) string {
	if path == "(root)" {
		path = ""
	}

	if path == "" {
		path = "$"
	} else {
		path = fmt.Sprintf("$.%s", path)
	}

	// Index
	// $.foo.0 -> $.foo[0]
	path = resultPathNormalizeRegex.ReplaceAllString(path, "[${1}]")

	return path
}

type NodeValidationResultError struct {
	resultErr validation.ResultErrorInterface
	nodeErr   *NodeError
	*serrors.Arguments
	*serrors.Details
}

func (err *NodeValidationResultError) Path() string {
	return err.resultErr.Path()
}

func (err *NodeValidationResultError) Error() string {
	return err.nodeErr.Error()
}

func (err *NodeValidationResultError) ErrorArguments() []any {
	err.AppendArguments(err.resultErr.ErrorArguments()...)
	err.AppendArguments(err.nodeErr.ErrorArguments()...)
	return err.Arguments.ErrorArguments()
}

func (err *NodeValidationResultError) ErrorDetails(ansi bool) string {
	return err.nodeErr.ErrorDetails(ansi)
}
