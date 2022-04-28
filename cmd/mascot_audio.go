//go:build darwin || windows

package cmd

import (
	"bytes"
	_ "embed"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"math/rand"
	"time"
)

//go:embed mascot/mascot_yell.wav
var mascotAudioYell []byte

var mascotAudioStreamer beep.StreamSeekCloser

func (model *mascotModel) yellAudio() error {
	// Lazy load streamer and init speaker at the same time
	if mascotAudioStreamer == nil {
		var format beep.Format
		var err error
		mascotAudioStreamer, format, err = wav.Decode(bytes.NewReader(mascotAudioYell))
		if err != nil {
			return err
		}

		if err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/25)); err != nil {
			return err
		}
	}

	if err := mascotAudioStreamer.Seek(0); err != nil {
		return err
	}

	// Random resampling ratio
	streamer := beep.ResampleRatio(4, rand.Float64()+0.5, mascotAudioStreamer)

	// Play
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	<-done

	return nil
}
