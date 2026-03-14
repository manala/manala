package mascot

import (
	"bytes"
	_ "embed"
	"image"
	"image/draw"
	"image/png"
	"io"
	"math/rand/v2"
	"time"

	"github.com/manala/manala/cmd"
	"github.com/manala/manala/internal/image/asciify"
	"github.com/manala/manala/internal/image/scale"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/term"
)

//go:embed assets/frames.png
var framesFile []byte

func RunMascot(in io.Reader, out io.Writer, repeat int) error {
	var (
		mascot Mascot
		err    error
	)

	// Ensure we're running in a terminal before bubbletea tries to use it
	if f, ok := out.(term.File); !ok || !term.IsTerminal(f.Fd()) {
		return &cmd.TerminalNotFoundError{}
	}

	mascot, err = NewMascot(repeat)
	if err != nil {
		return err
	}

	// Audio
	if m, ok := mascot.audio(); ok {
		defer m.audioStop()
	}

	p := tea.NewProgram(
		mascot,
		tea.WithInput(in),
		tea.WithOutput(out),
	)

	m, err := p.Run()
	if err != nil {
		return err
	}

	mascot = m.(Mascot)
	if mascot.err != nil {
		return mascot.err
	}

	return nil
}

type audio interface {
	audioInit() error
	audioPlay() tea.Cmd
	audioStop()
}

type (
	idleDoneMsg struct{}
	yellDoneMsg struct {
		err error
	}
)

type Mascot struct {
	framesImg  []*image.NRGBA
	framesStr  []string
	frameIndex int
	duration   int
	repeat     int
	quitting   bool
	err        error
}

func NewMascot(repeat int) (Mascot, error) {
	mascot := Mascot{
		duration: 345,
		repeat:   repeat,
	}

	// Decode frames
	framesSrc, err := png.Decode(bytes.NewReader(framesFile))
	if err != nil {
		return mascot, err
	}

	// NRGBA only
	var frames *image.NRGBA
	if nrgba, ok := framesSrc.(*image.NRGBA); ok {
		frames = nrgba
	} else {
		bounds := framesSrc.Bounds()
		frames = image.NewNRGBA(bounds)
		draw.Draw(frames, bounds, framesSrc, bounds.Min, draw.Src)
	}

	framesBounds := frames.Bounds()
	framesHeight := framesBounds.Dy() / 2

	mascot.framesImg = make([]*image.NRGBA, 2)
	for i := range mascot.framesImg {
		frameRect := image.Rect(framesBounds.Min.X, framesBounds.Min.Y+i*framesHeight, framesBounds.Max.X, framesBounds.Min.Y+(i+1)*framesHeight)
		mascot.framesImg[i] = frames.SubImage(frameRect).(*image.NRGBA)
	}

	// Init audio
	if m, ok := mascot.audio(); ok {
		if err = m.audioInit(); err != nil {
			return mascot, err
		}
	}

	return mascot, nil
}

func (mascot Mascot) Init() tea.Cmd {
	return mascot.idle()
}

func (mascot Mascot) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		// Quit
		case "ctrl+c", "q", "esc":
			return mascot, tea.Quit
		}
	// Scale
	case tea.WindowSizeMsg:
		mascot.framesStr = mascot.scale(msg.Width, msg.Height)
		return mascot, nil
	// Yell
	case idleDoneMsg:
		if mascot.quitting {
			return mascot, tea.Quit
		}

		mascot.frameIndex = 1
		return mascot, mascot.yell()
	// Idle
	case yellDoneMsg:
		if msg.err != nil {
			mascot.err = msg.err
			return mascot, tea.Quit
		}

		// Repeat counter
		if mascot.repeat > 0 {
			mascot.repeat--
		}
		if mascot.repeat == 0 {
			mascot.quitting = true
		}

		mascot.frameIndex = 0
		return mascot, mascot.idle()
	}

	return mascot, nil
}

func (mascot Mascot) View() tea.View {
	view := tea.NewView("")
	view.AltScreen = true

	if len(mascot.framesStr) > 0 {
		view.SetContent(mascot.framesStr[mascot.frameIndex])
	}

	return view
}

func (mascot Mascot) audio() (audio, bool) {
	m, ok := any(mascot).(audio)
	return m, ok
}

func (mascot Mascot) idle() tea.Cmd {
	duration := (mascot.duration / 2) + rand.IntN(mascot.duration)
	return tea.Tick(
		time.Duration(duration)*time.Millisecond,
		func(time.Time) tea.Msg { return idleDoneMsg{} },
	)
}

func (mascot Mascot) yell() tea.Cmd {
	// Play audio and wait for completion
	if m, ok := mascot.audio(); ok {
		return m.audioPlay()
	}

	// No audio, just hold the yell frame
	duration := (mascot.duration / 2) + rand.IntN(mascot.duration)

	return tea.Tick(
		time.Duration(duration)*time.Millisecond,
		func(time.Time) tea.Msg { return yellDoneMsg{} },
	)
}

func (mascot Mascot) scale(width, height int) []string {
	// Frames images
	framesImg := make([]*image.NRGBA, len(mascot.framesImg))
	for i, img := range mascot.framesImg {
		framesImg[i] = scale.Scale(img, width, height*2)
	}

	// Frames strings
	framesStr := make([]string, len(framesImg))
	for i, frame := range framesImg {
		framesStr[i] = asciify.Asciify(frame)
	}

	return framesStr
}
