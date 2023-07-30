package file

import (
	"github.com/stretchr/testify/assert"
	"os"
)

func EqualContent(s *assert.Assertions, expected string, path string) {
	if !s.FileExists(path) {
		return
	}

	content, err := os.ReadFile(path)
	s.NoError(err)

	s.Equal(expected, string(content))
}
