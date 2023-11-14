//go:build darwin || windows

package charm

import (
	"bytes"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"math/rand"
	"time"
)

/*********/
/* Model */
/*********/

var animationAudioStreamer beep.StreamSeekCloser

func (model animationModel) yellAudio() error {
	// Lazy load streamer and init speaker at the same time
	if animationAudioStreamer == nil {
		var format beep.Format
		var err error
		animationAudioStreamer, format, err = wav.Decode(bytes.NewReader(*model.animation.AudioYell))
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
