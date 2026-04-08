package parsing

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/v2/quick"
)

const Context = 3

type Dumper struct {
	Err   *Error
	Src   string
	Lexer string
}

func (d *Dumper) Dump(ansi bool) string {
	source := d.Src

	if ansi {
		var h strings.Builder
		if err := quick.Highlight(&h, source, d.Lexer, "terminal16m", "native"); err == nil {
			source = h.String()
		}
	}

	lines := strings.Split(source, "\n")

	minLine := max(d.Err.Line-Context, 1)
	maxLine := min(d.Err.Line+Context, len(lines))

	for minLine < d.Err.Line && lines[minLine-1] == "" {
		minLine++
	}

	for maxLine > d.Err.Line && lines[maxLine-1] == "" {
		maxLine--
	}

	width := len(strconv.Itoa(maxLine))

	var b strings.Builder

	for i := minLine; i <= maxLine; i++ {
		if i == d.Err.Line {
			_, _ = fmt.Fprintf(&b, "> %*d | %s\n", width, i, lines[i-1])
			_, _ = fmt.Fprintf(&b, "  %*s   %s^\n", width, "", strings.Repeat(" ", d.Err.Column-1))
		} else {
			_, _ = fmt.Fprintf(&b, "  %*d | %s\n", width, i, lines[i-1])
		}
	}

	_, _ = fmt.Fprintf(&b, "* %v\n", d.Err.Err)

	return b.String()
}
