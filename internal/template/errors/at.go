package errors

import "strings"

// At creates an Error positioned at the given line and 0-based byte column.
// The byte column is converted to a 1-based rune column by iterating over runes.
// A zero byte column means no column info — column is left at 0.
func At(err error, src string, line, byteColumn int) Error {
	e := Error{
		error:  err,
		line:   line,
		column: 0,
	}

	if src == "" || line <= 0 {
		return e
	}

	// Extract target line, then count runes up to byteColumn (0-based byte offset)
	lines := strings.SplitN(src, "\n", line+1)
	if line > len(lines) {
		return e
	}

	lineContent := lines[line-1]
	e.column = 1
	for range lineContent[:min(byteColumn, len(lineContent))] {
		e.column++
	}

	return e
}
