package heredoc

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Doc(doc string, args ...any) string {
	lines := strings.Split(doc, "\n")

	// Remove first line if trailing
	firstLine := lines[0]
	if firstLine == "" && len(lines) != 2 {
		lines = lines[1:]
	}

	// Skip if empty
	if len(lines) == 0 {
		return ""
	}

	// Remove last line if only made of indentation runes
	lastLine := lines[len(lines)-1]
	if lastLine == indent(lastLine) {
		lines[len(lines)-1] = ""
	}

	// Find shorter line indentation
	shorterIndent := ""
	for _, line := range lines {
		lineIndent := indent(line)
		if lineIndent != "" && (len(lineIndent) < len(shorterIndent) || shorterIndent == "") {
			shorterIndent = lineIndent
		}
	}

	// Remove indentation for all lines
	for i, line := range lines {
		lines[i], _ = strings.CutPrefix(line, shorterIndent)
	}

	return fmt.Sprintf(
		strings.Join(lines, "\n"),
		args...,
	)
}

func indent(s string) string {
	var indent strings.Builder
	for _, r := range s {
		if r == '\t' {
			indent.WriteRune(r)
		} else {
			break
		}
	}

	return indent.String()
}

func Equal(t *testing.T, expected string, actual any, args ...any) {
	strs := strings.Split(fmt.Sprintf("%s", actual), "\n")

	// Trim actual strings
	actuals := make([]string, len(strs))
	for i, str := range strs {
		actuals[i] = strings.TrimRight(str, " ")
	}

	assert.Equal(t,
		Doc(expected, args...),
		strings.Join(actuals, "\n"),
	)
}

func EqualFile(t *testing.T, expected string, path string, args ...any) {
	content, err := os.ReadFile(path)
	require.NoError(t, err)

	// Fix line endings
	content = bytes.ReplaceAll(content, []byte("\r\n"), []byte{'\n'})

	assert.Equal(t,
		Doc(expected, args...),
		string(content),
	)
}
