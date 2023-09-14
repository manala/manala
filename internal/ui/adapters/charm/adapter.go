package charm

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"io"
	"slices"
)

func New(in io.Reader, out io.Writer, err io.Writer) *Adapter {
	return &Adapter{
		in:          in,
		outRenderer: lipgloss.NewRenderer(out, termenv.WithColorCache(true)),
		errRenderer: lipgloss.NewRenderer(err, termenv.WithColorCache(true)),
	}
}

type Adapter struct {
	in          io.Reader
	outRenderer *lipgloss.Renderer
	errRenderer *lipgloss.Renderer
}

/****************/
/* Models Index */
/****************/

func newModelsIndex(models *[]tea.Model) *modelsIndex {
	return &modelsIndex{
		models: models,
	}
}

type modelsIndex struct {
	models   *[]tea.Model
	index    int
	Circular bool
}

func (i *modelsIndex) Get() int {
	return i.index
}

func (i *modelsIndex) Set(index int) {
	i.index = index
}

func (i *modelsIndex) Next() int {
	if i.Circular {
		return (i.index + 1) % len(*i.models)
	}
	index := i.index + 1
	if index == len(*i.models) {
		index--
	}
	return index
}

func (i *modelsIndex) SetNext() {
	i.Set(i.Next())
}

func (i *modelsIndex) Previous() int {
	index := i.index - 1
	if index < 0 {
		if i.Circular {
			index = len(*i.models) - 1
		} else {
			index = 0
		}
	}
	return index
}

func (i *modelsIndex) SetPrevious() {
	i.Set(i.Previous())
}

func (i *modelsIndex) Reset() {
	i.Set(i.First())
}

func (i *modelsIndex) First() int {
	return 0
}

func (i *modelsIndex) Last() int {
	return len(*i.models) - 1
}

/********/
/* Cmds */
/********/

func newCmds() *cmds {
	cmds := make(cmds, 0)
	return &cmds
}

type cmds []tea.Cmd

func (c *cmds) Add(cmds ...tea.Cmd) *cmds {
	*c = append(*c,
		slices.DeleteFunc(cmds, func(cmd tea.Cmd) bool {
			return cmd == nil
		})...,
	)
	return c
}

func (c *cmds) AddSequence(cmds ...tea.Cmd) *cmds {
	c.Add(tea.Sequence(cmds...))
	return c
}

func (c *cmds) Init(models ...tea.Model) *cmds {
	for _, model := range models {
		c.Add(model.Init())
	}
	return c
}

func (c *cmds) Update(model tea.Model, msgs ...tea.Msg) tea.Model {
	var cmd tea.Cmd
	for _, msg := range msgs {
		model, cmd = model.Update(msg)
		c.Add(cmd)
	}
	return model
}

func (c *cmds) Batch() tea.Cmd {
	switch len(*c) {
	case 0:
		return nil
	case 1:
		return (*c)[0]
	}

	return tea.Batch(*c...)
}

func (c *cmds) Sequence() tea.Cmd {
	switch len(*c) {
	case 0:
		return nil
	case 1:
		return (*c)[0]
	}

	return tea.Sequence(*c...)
}

/************/
/* Messages */
/************/

type focusMsg bool
type hoverMsg bool
