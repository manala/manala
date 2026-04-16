package log

import (
	"fmt"
	"io"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

const messageWidth = 33

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
	var attrs [][2]any
	for i := 0; i+1 < len(args); i += 2 {
		attrs = append(attrs, [2]any{args[i], args[i+1]})
	}
	var b strings.Builder
	l.log(&b, level, msg, attrs)
	_, _ = lipgloss.Fprint(l.out, b.String())
}

func (l *Log) Error(err error) {
	var b strings.Builder
	l.error(&b, err, 0)
	_, _ = lipgloss.Fprint(l.out, b.String())
}

func (l *Log) error(w io.Writer, err error, depth int) {
	// Attrs
	var attrs [][2]any
	if a, ok := err.(Attrs); ok {
		attrs = a.Attrs()
	}

	l.log(w, LevelError, err.Error(), attrs)

	// Dumper
	if d, ok := err.(Dumper); ok {
		if dumper := d.Dumper(); dumper != nil {
			var b strings.Builder
			dumper.Dump(&b)
			if dump := b.String(); dump != "" {
				dumpStyle := lipgloss.NewStyle().Margin(1, 0, 0, 3)
				_, _ = io.WriteString(w, dumpStyle.Render(dump))
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
		l.error(w, child, depth+1)
	}
}

func (l *Log) log(w io.Writer, level Level, msg string, attrs [][2]any) {
	bulletStyle := lipgloss.NewStyle().PaddingLeft(1)
	messageStyle := lipgloss.NewStyle().PaddingLeft(1)
	if len(attrs) > 0 {
		msg = ansi.Truncate(msg, messageWidth, "")
		messageStyle = messageStyle.Width(messageWidth)
	}
	switch level {
	case LevelError:
		_, _ = io.WriteString(w, bulletStyle.Foreground(lipgloss.Red).Render("✖"))
		_, _ = io.WriteString(w, messageStyle.Foreground(lipgloss.Red).Bold(true).Render(msg))
	case LevelWarn:
		_, _ = io.WriteString(w, bulletStyle.Foreground(lipgloss.Yellow).Render("▲"))
		_, _ = io.WriteString(w, messageStyle.Foreground(lipgloss.Yellow).Bold(true).Render(msg))
	case LevelInfo:
		_, _ = io.WriteString(w, bulletStyle.Foreground(lipgloss.Green).Render("•"))
		_, _ = io.WriteString(w, messageStyle.Render(msg))
	case LevelDebug:
		_, _ = io.WriteString(w, bulletStyle.Foreground(lipgloss.Blue).Render("◦"))
		_, _ = io.WriteString(w, messageStyle.Render(msg))
	}

	if len(attrs) > 0 {
		keyStyle := lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Cyan)
		valueStyle := lipgloss.NewStyle().Bold(true)
		for _, attr := range attrs {
			_, _ = io.WriteString(w, keyStyle.Render(fmt.Sprintf("%v", attr[0])))
			_, _ = io.WriteString(w, "=")
			_, _ = io.WriteString(w, valueStyle.Render(fmt.Sprintf("%v", attr[1])))
		}
	}

	_, _ = io.WriteString(w, "\n")
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
