package parsing

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/v2/quick"
)

const Context = 3

type ErrorDumper struct {
	Err   *Error
	File  string
	Src   string
	Lexer string
}

func (d ErrorDumper) Dump(w io.Writer) {
	src := d.Src

	// Highlighting
	var h strings.Builder
	if err := quick.Highlight(&h, src, d.Lexer, "terminal16m", "native"); err == nil {
		src = h.String()
	}

	lines := strings.Split(src, "\n")

	minLine := max(d.Err.Line-Context, 1)
	maxLine := min(d.Err.Line+Context, len(lines))

	for minLine < d.Err.Line && lines[minLine-1] == "" {
		minLine++
	}

	for maxLine > d.Err.Line && lines[maxLine-1] == "" {
		maxLine--
	}

	width := len(strconv.Itoa(maxLine))

	if d.File != "" {
		_, _ = fmt.Fprintf(w, "in %s:%d", d.File, d.Err.Line)
		if d.Err.Column != 0 {
			_, _ = fmt.Fprintf(w, ":%d", d.Err.Column)
		}
		_, _ = fmt.Fprintf(w, "\n")
	}

	for i := minLine; i <= maxLine; i++ {
		if i == d.Err.Line {
			_, _ = fmt.Fprintf(w, "> %*d | %s\n", width, i, lines[i-1])
			if d.Err.Column > 0 {
				_, _ = fmt.Fprintf(w, "  %*s   %s^\n", width, "", strings.Repeat(" ", d.Err.Column-1))
			}
		} else {
			_, _ = fmt.Fprintf(w, "  %*d | %s\n", width, i, lines[i-1])
		}
	}

	_, _ = fmt.Fprintf(w, "* %v\n", d.Err.Err)
}
