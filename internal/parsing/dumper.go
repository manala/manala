package parsing

import (
	"fmt"
	"io"
	"strconv"
	"strings"

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

func (d ErrorDumper) Dump(w io.Writer) {
	// File
	atStyle := lipgloss.NewStyle().Faint(true)
	fileStyle := lipgloss.NewStyle().Foreground(lipgloss.Cyan)
	if d.File != "" {
		file := fmt.Sprintf("%s:%d", d.File, d.Err.Line)
		if d.Err.Column != 0 {
			file += fmt.Sprintf(":%d", d.Err.Column)
		}
		_, _ = fmt.Fprintf(w, "%s %s\n\n",
			atStyle.Render("at"),
			fileStyle.Render(file),
		)
	}

	// Source
	src := d.Src

	// Highlight
	var highlight strings.Builder
	if err := quick.Highlight(&highlight, src, d.Lexer, "terminal256", "dracula"); err == nil {
		src = highlight.String()
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

	caretStyle := lipgloss.NewStyle().Foreground(lipgloss.Red).Bold(true)
	gutterStyle := lipgloss.NewStyle().Faint(true)
	cursorStyle := lipgloss.NewStyle()
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Red).Bold(true).Italic(true)

	// Snippet
	for i := lineMin; i <= lineMax; i++ {
		gutter := fmt.Sprintf("%*d │", gutterWidth, i)

		if i == d.Err.Line {
			_, _ = fmt.Fprint(w, caretStyle.Render("▶ "))
			_, _ = fmt.Fprint(w, cursorStyle.Render(gutter))
		} else {
			_, _ = fmt.Fprint(w, "  ")
			_, _ = fmt.Fprint(w, gutterStyle.Render(gutter))
		}

		if line := lines[i-1]; lipgloss.Width(line) > 0 {
			_, _ = fmt.Fprint(w, " "+line)
		}
		_, _ = fmt.Fprint(w, "\n")

		if i == d.Err.Line {
			_, _ = fmt.Fprintf(w, "  %*s", gutterWidth, "")
			if d.Err.Column == 0 {
				_, _ = fmt.Fprint(w, cursorStyle.Render(" ├ "))
			} else {
				_, _ = fmt.Fprint(w, cursorStyle.Render(
					" ├"+
						strings.Repeat("─", d.Err.Column)+
						"╯ ",
				))
			}
			_, _ = fmt.Fprint(w, errorStyle.Render(d.Err.Err.Error())+"\n")
		}
	}
}
