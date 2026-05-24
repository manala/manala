package source

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/manala/manala/internal/output"

	"charm.land/lipgloss/v2"
	"github.com/alecthomas/chroma/v2/quick"
)

const Context = 3

type Origin struct {
	Source   string
	Language string
	File     string // optional — only set when the source is a file on disk
}

type Error struct {
	Origin

	Position Position // deepest Position in the chain — carries the original error message
	Line     int      // absolute line (accumulated from all relative offsets in the chain)
	Column   int
}

func (e Error) Error() string {
	return e.Position.Error()
}

func (e Error) Attrs() [][2]any {
	if e.File != "" {
		return [][2]any{
			{"file", e.File},
		}
	}
	return nil
}

// Unwrap continues the descent from below this Error's Position, looking for deeper ones.
// Each deeper Position found is returned as a new Error with updated accumulated offsets.
// Returns nil when there are no more Positions — this Error is then the terminal node.
func (e Error) Unwrap() []error {
	return e.expand(e.Position.Unwrap(), e.Origin, e.Line, e.Column)
}

func (e Error) Render(p output.Profile) string {
	var b strings.Builder

	// File
	if e.File != "" {
		file := e.File
		if e.Line > 0 {
			file += fmt.Sprintf(":%d", e.Line)
			if e.Column > 0 {
				file += fmt.Sprintf(":%d", e.Column)
			}
		}
		_, _ = fmt.Fprintf(&b, "\n%s %s\n",
			p.MutedStyle().Render("at"),
			p.LitteralStyle().Render(file),
		)
	}
	b.WriteByte('\n')

	src := e.Source

	// Highlight
	if p.Rich() {
		formatter := "terminal16"
		switch {
		case p.Extended():
			formatter = "terminal256"
		case p.True():
			formatter = "terminal16m"
		}
		style := "catppuccin-mocha"
		if p.Light() {
			style = "catppuccin-latte"
		}
		var h strings.Builder
		if err := quick.Highlight(&h, src, e.Language, formatter, style); err == nil {
			src = h.String()
		}
	}

	lines := strings.Split(src, "\n")

	lineMin := max(e.Line-Context, 1)
	lineMax := min(e.Line+Context, len(lines))
	for lineMin < lineMax && lineMin < e.Line && lipgloss.Width(lines[lineMin-1]) == 0 {
		lineMin++
	}
	for lineMax > lineMin && lineMax > e.Line && lipgloss.Width(lines[lineMax-1]) == 0 {
		lineMax--
	}

	gutterWidth := len(strconv.Itoa(lineMax))

	for i := lineMin; i <= lineMax; i++ {
		gutter := fmt.Sprintf("%*d │", gutterWidth, i)

		if i == e.Line {
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

		if i == e.Line {
			_, _ = fmt.Fprintf(&b, "  %*s", gutterWidth, "")
			if e.Column == 0 {
				b.WriteString(p.ErrorStyle().Render(" ├ "))
			} else {
				b.WriteString(p.ErrorStyle().Render(
					" ├" +
						strings.Repeat("─", e.Column) +
						"╯ ",
				))
			}
			b.WriteString(p.ErrorStyle().Render(e.Position.Error()) + "\n")
		}
	}

	// Message footer when the focal line falls outside the rendered window
	if e.Line < lineMin || e.Line > lineMax {
		b.WriteByte('\n')
		b.WriteString(p.ErrorStyle().Render(e.Position.Error()) + "\n")
	}

	return b.String()
}

// expand follows the same traversal as from, but wraps each Position found into an Error
// instead of collecting it. It stops at the first Position per branch — the new Error's
// own Unwrap() will continue the descent when needed.
func (e Error) expand(err error, origin Origin, absLine, absCol int) []error {
	for err != nil {
		if pos, ok := err.(Position); ok {
			l, c := pos.Position()
			return []error{Error{
				Origin:   origin,
				Position: pos,
				Line:     absLine + max(l-1, 0),
				Column:   absCol + max(c-1, 0),
			}}
		}
		switch x := err.(type) {
		case interface{ Unwrap() []error }:
			// fan out: each branch may independently contain a deeper Position
			var result []error
			for _, child := range x.Unwrap() {
				result = append(result, e.expand(child, origin, absLine, absCol)...)
			}
			return result
		case interface{ Unwrap() error }:
			// not a Position — keep descending
			err = x.Unwrap()
		default:
			// true leaf with no deeper Position
			return nil
		}
	}
	return nil
}

type Position interface {
	error
	Position() (int, int)
	Unwrap() error
}
