package lipgloss

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"io"
)

var (
	color = lipgloss.AdaptiveColor{Light: "#1A1A1A", Dark: "#DDDDDD"}
)

func New(out io.Writer, err io.Writer) *Output {
	return &Output{
		outRenderer: lipgloss.NewRenderer(out),
		errRenderer: lipgloss.NewRenderer(err),
	}
}

type Output struct {
	outRenderer *lipgloss.Renderer
	errRenderer *lipgloss.Renderer
}

func (output *Output) writeOutString(s string) {
	_, _ = output.outRenderer.Output().WriteString(s)
}

func (output *Output) writeErrString(s string) {
	_, _ = output.errRenderer.Output().WriteString(s)
}

func (output *Output) outStyle() lipgloss.Style {
	return output.outRenderer.NewStyle()
}

func (output *Output) errStyle() lipgloss.Style {
	return output.errRenderer.NewStyle()
}

func (output *Output) errAnsi() bool {
	return output.errRenderer.ColorProfile() != termenv.Ascii
}
