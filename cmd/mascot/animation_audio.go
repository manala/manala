//go:build darwin || windows

package mascot

import (
	"bytes"
	"math/rand/v2"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/wav"
)

var animationAudioStreamer beep.StreamSeekCloser

func (model Animation) yellAudio() error {
	// Lazy load streamer and init speaker at the same time
	if animationAudioStreamer == nil {
		var (
			format beep.Format
			err    error
		)

		animationAudioStreamer, format, err = wav.Decode(bytes.NewReader(*model.AudioYell))
		if err != nil {
			return err
		}

		if err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/25)); err != nil {
			return err
		}
	}

	if err := animationAudioStreamer.Seek(0); err != nil {
		return err
	}

	// Random resampling ratio
	streamer := beep.ResampleRatio(4, rand.Float64()+0.5, animationAudioStreamer)

	// Play
	done := make(chan bool)

	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done

	return nil
}
