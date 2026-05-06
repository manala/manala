package output

import (
	"fmt"
	"io"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/x/term"
)

type Output struct {
	Profile

	out io.Writer
}

func New(in, out term.File, env []string) Output {
	return Output{
		out: out,
		Profile: Profile{
			light:   !lipgloss.HasDarkBackground(in, out),
			profile: colorprofile.Detect(out, env),
		},
	}
}

func NewDetached(out io.Writer) Output {
	return Output{
		out: out,
	}
}

func (o Output) Print(a ...any) {
	_, _ = fmt.Fprint(o.out, a...)
}

func (o Output) Println(a ...any) {
	_, _ = fmt.Fprintln(o.out, a...)
}
