package yaml

import (
	"fmt"
	"regexp"
)

func NewJsonPathNormalizer(path string) *JsonPathNormalizer {
	return &JsonPathNormalizer{
		path: path,
	}
}

var jsonPathNormalizerIndexRegex = regexp.MustCompile(`\.(\d+)`)

type JsonPathNormalizer struct {
	path string
}

func (normalizer *JsonPathNormalizer) Normalize() string {
	path := normalizer.path

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
	path = jsonPathNormalizerIndexRegex.ReplaceAllString(path, "[${1}]")

	return path
}
