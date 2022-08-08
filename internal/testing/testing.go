package testing

import (
	"github.com/stretchr/testify/suite"
	"path/filepath"
)

func Path(suite suite.TestingSuite, path ...string) string {
	return filepath.Join(
		filepath.FromSlash(suite.T().Name()),
		filepath.Join(path...),
	)
}

func DataPath(suite suite.TestingSuite, path ...string) string {
	return filepath.Join(
		"testdata",
		Path(suite, path...),
	)
}
