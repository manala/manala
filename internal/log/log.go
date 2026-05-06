package log

import (
	"fmt"
	"strings"

	"github.com/manala/manala/internal/output"

	"charm.land/lipgloss/v2"
)

type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
)

const messageWidth = 33

type Log struct {
	out     output.Output
	verbose int
}

func New(out output.Output) *Log {
	return &Log{
		out: out,
	}
}

func (l *Log) Verbose(verbose int) {
	l.verbose = verbose
}

func (l *Log) Debug(msg string, args ...any) {
	l.Log(Debug, msg, args...)
}

func (l *Log) Info(msg string, args ...any) {
	l.Log(Info, msg, args...)
}

func (l *Log) Warn(msg string, args ...any) {
	l.Log(Warn, msg, args...)
}

func (l *Log) Log(level Level, msg string, args ...any) {
	if level == Debug && l.verbose < 2 {
		return
	}
	if level == Info && l.verbose < 1 {
		return
	}

	// Attrs
	var attrs [][2]any
	for i := 0; i+1 < len(args); i += 2 {
		attrs = append(attrs, [2]any{args[i], args[i+1]})
	}
	l.out.Println(l.log(level, msg, attrs))
}

func (l *Log) log(level Level, msg string, attrs [][2]any) string {
	var b strings.Builder

	messageStyle := lipgloss.NewStyle().PaddingLeft(1)
	if len(attrs) > 0 {
		messageStyle = messageStyle.Width(messageWidth)
	}

	switch level {
	case Error:
		b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top,
			l.out.ErrorStyle().PaddingLeft(1).Render("✖"),
			messageStyle.Inherit(l.out.ErrorStyle()).Render(msg),
		))
	case Warn:
		b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top,
			l.out.WarnStyle().PaddingLeft(1).Render("▲"),
			messageStyle.Inherit(l.out.WarnStyle()).Render(msg),
		))
	case Info:
		b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top,
			l.out.InfoStyle().PaddingLeft(1).Render("●"),
			messageStyle.Inherit(l.out.Style()).Render(msg),
		))
	case Debug:
		b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top,
			l.out.DebugStyle().PaddingLeft(1).Render("○"),
			messageStyle.Inherit(l.out.Style()).Render(msg),
		))
	}

	if len(attrs) > 0 {
		keyStyle := l.out.MutedStyle().PaddingLeft(1)
		valueStyle := l.out.LitteralStyle()
		for _, attr := range attrs {
			_, _ = fmt.Fprintf(&b, "%s%s",
				keyStyle.Render(fmt.Sprintf("%v=", attr[0])),
				valueStyle.Render(fmt.Sprintf("%v", attr[1])),
			)
		}
	}

	return b.String()
}

func (l *Log) indent(s string, n int) string {
	pad := strings.Repeat(" ", n)
	return pad + strings.ReplaceAll(s, "\n", "\n"+pad)
}

func (l *Log) block(s string) string {
	return " │ " + strings.ReplaceAll(s, "\n", "\n │ ")
}

type Attrs interface {
	Attrs() [][2]any
}
