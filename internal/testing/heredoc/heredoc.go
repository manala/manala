package heredoc

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
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

func Equal(s *assert.Assertions, expected string, actual string, args ...any) {
	// Trim actual
	var actuals []string
	for _, _actual := range strings.Split(actual, "\n") {
		actuals = append(actuals, strings.TrimRight(_actual, " "))
	}

	s.Equal(
		Doc(expected, args...),
		strings.Join(actuals, "\n"),
	)
}

func EqualFile(s *assert.Assertions, expected string, path string, args ...any) {
	content, err := os.ReadFile(path)
	s.NoError(err)

	s.Equal(
		Doc(expected, args...),
		string(content),
	)
}
