package heredoc

import (
	"fmt"
	"strings"
)

func Doc(doc string) string {
	lines := strings.Split(doc, "\n")

	// Remove first line if trailing
	firstLine := lines[0]
	if firstLine == "" {
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

	return strings.Join(lines, "\n")
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

func Docf(doc string, args ...interface{}) string {
	return fmt.Sprintf(Doc(doc), args...)
}
