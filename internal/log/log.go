package log

import (
	"fmt"
	"io"
	"strings"

	"charm.land/lipgloss/v2"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

const messageWidth = 33

var Discard *Log = &Log{out: io.Discard}

type Log struct {
	out     io.Writer
	verbose int
}

func New(out io.Writer) *Log {
	return &Log{
		out: out,
	}
}

func (l *Log) Verbose(verbose int) {
	l.verbose = verbose
}

func (l *Log) Debug(msg string, args ...any) {
	l.Log(LevelDebug, msg, args...)
}

func (l *Log) Info(msg string, args ...any) {
	l.Log(LevelInfo, msg, args...)
}

func (l *Log) Warn(msg string, args ...any) {
	l.Log(LevelWarn, msg, args...)
}

func (l *Log) Log(level Level, msg string, args ...any) {
	if level == LevelDebug && l.verbose < 2 {
		return
	}
	if level == LevelInfo && l.verbose < 1 {
		return
	}

	// Attrs
	var attrs [][2]any
	for i := 0; i+1 < len(args); i += 2 {
		attrs = append(attrs, [2]any{args[i], args[i+1]})
	}
	_, _ = lipgloss.Fprintln(l.out, l.log(level, msg, attrs))
}

func (l *Log) Error(err error) {
	_, _ = lipgloss.Fprintln(l.out, l.error(err, 0))
}

func (l *Log) error(err error, depth int) string {
	var b strings.Builder

	// Attrs
	var attrs [][2]any
	if a, ok := err.(Attrs); ok {
		attrs = a.Attrs()
	}

	b.WriteString(
		lipgloss.NewStyle().MarginLeft(depth*2).Render(
			l.log(LevelError, err.Error(), attrs),
		) + "\n",
	)

	// Dumper
	if e, ok := err.(Dumper); ok {
		if d := e.Dumper(); d != nil {
			var bd strings.Builder
			d.Dump(&bd)
			if dump := bd.String(); dump != "" {
				b.WriteString(
					lipgloss.NewStyle().Margin(1, 0, 0, 3+depth*2).Render(dump),
				)
			}
		}
	}

	// Children errors
	var children []error
	switch u := err.(type) {
	case interface{ Unwrap() error }:
		if child := u.Unwrap(); child != nil {
			children = append(children, child)
		}
	case interface{ Unwrap() []error }:
		for _, child := range u.Unwrap() {
			if child != nil {
				children = append(children, child)
			}
		}
	}
	for _, child := range children {
		b.WriteString(l.error(child, depth+1))
	}

	return b.String()
}

func (l *Log) log(level Level, msg string, attrs [][2]any) string {
	var b strings.Builder

	bulletStyle := lipgloss.NewStyle().PaddingLeft(1)
	messageStyle := lipgloss.NewStyle().PaddingLeft(1)
	if len(attrs) > 0 {
		msg = lipgloss.Wrap(msg, messageWidth, "")
		messageStyle = messageStyle.Width(messageWidth)
	}

	switch level {
	case LevelError:
		b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top,
			bulletStyle.Foreground(lipgloss.Red).Render("✖"),
			messageStyle.Foreground(lipgloss.Red).Bold(true).Render(msg),
		))
	case LevelWarn:
		b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top,
			bulletStyle.Foreground(lipgloss.Yellow).Render("▲"),
			messageStyle.Foreground(lipgloss.Yellow).Bold(true).Render(msg),
		))
	case LevelInfo:
		b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top,
			bulletStyle.Foreground(lipgloss.Green).Render("•"),
			messageStyle.Render(msg),
		))
	case LevelDebug:
		b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top,
			bulletStyle.Foreground(lipgloss.Blue).Render("◦"),
			messageStyle.Render(msg),
		))
	}

	if len(attrs) > 0 {
		keyStyle := lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Cyan)
		valueStyle := lipgloss.NewStyle().Bold(true)
		for _, attr := range attrs {
			_, _ = fmt.Fprintf(&b, "%s=%s",
				keyStyle.Render(fmt.Sprintf("%v", attr[0])),
				valueStyle.Render(fmt.Sprintf("%v", attr[1])),
			)
		}
	}

	return b.String()
}

type Attrs interface {
	Attrs() [][2]any
}

type Dumper interface {
	Dumper() dumper
}

type dumper = interface {
	Dump(w io.Writer)
}
