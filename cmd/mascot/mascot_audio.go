//go:build darwin || windows

package mascot

import (
	"bytes"
	_ "embed"
	"math/rand/v2"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/wav"
)

var (
	//go:embed assets/audio.wav
	audioFile     []byte
	audioStreamer beep.StreamSeekCloser
)

func (mascot Mascot) audioInit() error {
	var (
		format beep.Format
		err    error
	)

	audioStreamer, format, err = wav.Decode(bytes.NewReader(audioFile))
	if err != nil {
		return err
	}

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/25))
	if err != nil {
		return err
	}

	return nil
}

func (mascot Mascot) audioPlay() tea.Cmd {
	msg := yellDoneMsg{}

	return func() tea.Msg {
		if err := audioStreamer.Seek(0); err != nil {
			msg.err = err
			return msg
		}

		ratio := rand.Float64() + 0.5
		resampled := beep.ResampleRatio(4, ratio, audioStreamer)
		speaker.PlayAndWait(resampled)

		return msg
	}
}

func (mascot Mascot) audioStop() {
	_ = audioStreamer.Close()
}
