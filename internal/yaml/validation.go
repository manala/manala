package yaml

import (
	"fmt"
	"github.com/goccy/go-yaml"
	yamlAst "github.com/goccy/go-yaml/ast"
	"github.com/xeipuuv/gojsonschema"
	"io"
	internalReport "manala/internal/report"
)

func NewJsonLoader(node yamlAst.Node) gojsonschema.JSONLoader {
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
	if err := yaml.NewDecoder(loader.JsonSource().(yamlAst.Node)).Decode(&data); err != nil {
		// Nil or empty content
		if err == io.EOF {
			return nil, fmt.Errorf("empty content")
		}
		return nil, NewError(err)
	}

	return gojsonschema.NewGoLoader(data).LoadJSON()
}

func NewValidationReporter(node yamlAst.Node) *ValidationReporter {
	return &ValidationReporter{
		node: node,
	}
}

type ValidationReporter struct {
	node yamlAst.Node
}

func (reporter *ValidationReporter) Report(_ gojsonschema.ResultError, report *internalReport.Report) {
	NewReporter(reporter.node).Report(report)
}

func NewValidationPathReporter(node yamlAst.Node) *ValidationPathReporter {
	return &ValidationPathReporter{
		node: node,
	}
}

type ValidationPathReporter struct {
	node yamlAst.Node
}

func (reporter *ValidationPathReporter) Report(result gojsonschema.ResultError, report *internalReport.Report) {
	// Normalize json path
	path := NewJsonPathNormalizer(result.Field()).Normalize()

	// Get yaml path
	yamlPath, err := yaml.PathString(path)
	if err != nil {
		return
	}

	// Get node
	node, err := yamlPath.FilterNode(reporter.node)
	if err != nil || node == nil {
		return
	}

	NewReporter(node).Report(report)
}
