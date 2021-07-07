// +build darwin windows

package cmd

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"sync"
	"time"
)

func mascotRunAudio(cmd *MascotCmd, wg *sync.WaitGroup, errs chan<- error) {
	audio, err := cmd.Assets.Open("assets/mascot.wav")
	if err != nil {
		errs <- err
		return
	}

	streamer, format, err := wav.Decode(audio)
	if err != nil {
		errs <- err
		return
	}
	defer streamer.Close()

	if err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/25)); err != nil {
		errs <- err
		return
	}

	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		wg.Done()
	})))
}

func init() {
	mascotRun = append(
		[]mascotFunc{mascotRunAudio},
		mascotRun...
	)
}
