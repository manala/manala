//go:build darwin || windows

package mascot

import (
	"bytes"
	_ "embed"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"math/rand/v2"
	"time"
)

var (
	//go:embed assets/audio_yell.wav
	audioYell     []byte
	audioStreamer beep.StreamSeekCloser
)

func (model animation) yellAudio() error {
	// Lazy load streamer and init speaker at the same time
	if audioStreamer == nil {
		var format beep.Format
		var err error
		audioStreamer, format, err = wav.Decode(bytes.NewReader(audioYell))
		if err != nil {
			return err
		}

		if err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/25)); err != nil {
			return err
		}
	}

	if err := audioStreamer.Seek(0); err != nil {
		return err
	}

	// Random resampling ratio
	streamer := beep.ResampleRatio(4, rand.Float64()+0.5, audioStreamer)

	// Play
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	<-done

	return nil
}
