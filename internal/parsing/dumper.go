package parsing

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/manala/manala/internal/output"

	"charm.land/lipgloss/v2"
	"github.com/alecthomas/chroma/v2/quick"
)

const Context = 3

type ErrorDumper struct {
	Err   *Error
	File  string
	Src   string
	Lexer string
}

func (d ErrorDumper) Dump(p output.Profile) string {
	var b strings.Builder

	// File
	if d.File != "" {
		file := fmt.Sprintf("%s:%d", d.File, d.Err.Line)
		if d.Err.Column != 0 {
			file += fmt.Sprintf(":%d", d.Err.Column)
		}
		_, _ = fmt.Fprintf(&b, "%s %s\n\n",
			p.MutedStyle().Render("at"),
			p.LitteralStyle().Render(file),
		)
	}

	// Source
	src := d.Src

	// Highlight
	if p.Rich() {
		// Formatter
		formatter := "terminal16"
		switch {
		case p.Extended():
			formatter = "terminal256"
		case p.True():
			formatter = "terminal16m"
		}
		// Style
		style := "catppuccin-mocha"
		if p.Light() {
			style = "catppuccin-latte"
		}
		var h strings.Builder
		if err := quick.Highlight(&h, src, d.Lexer, formatter, style); err == nil {
			src = h.String()
		}
	}

	// Lines
	lines := strings.Split(src, "\n")
	lineMin := max(d.Err.Line-Context, 1)
	lineMax := min(d.Err.Line+Context, len(lines))
	for lineMin < d.Err.Line && lipgloss.Width(lines[lineMin-1]) == 0 {
		lineMin++
	}
	for lineMax > d.Err.Line && lipgloss.Width(lines[lineMax-1]) == 0 {
		lineMax--
	}

	gutterWidth := len(strconv.Itoa(lineMax))

	// Snippet
	for i := lineMin; i <= lineMax; i++ {
		gutter := fmt.Sprintf("%*d │", gutterWidth, i)

		if i == d.Err.Line {
			b.WriteString(p.ErrorStyle().Render("▶ "))
			b.WriteString(p.ErrorStyle().Render(gutter))
		} else {
			b.WriteString("  ")
			b.WriteString(p.MutedStyle().Render(gutter))
		}

		if line := lines[i-1]; lipgloss.Width(line) > 0 {
			b.WriteString(" " + line)
		}
		b.WriteString("\n")

		if i == d.Err.Line {
			_, _ = fmt.Fprintf(&b, "  %*s", gutterWidth, "")
			if d.Err.Column == 0 {
				b.WriteString(p.ErrorStyle().Render(" ├ "))
			} else {
				b.WriteString(p.ErrorStyle().Render(
					" ├" +
						strings.Repeat("─", d.Err.Column) +
						"╯ ",
				))
			}
			b.WriteString(p.ErrorStyle().Render(d.Err.Err.Error()) + "\n")
		}
	}

	return b.String()
}
