//go:build darwin || windows

package cmd

import (
	"bytes"
	_ "embed"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"sync"
	"time"
)

//go:embed embed/mascot.wav
var mascotAudio []byte

func mascotAudioRun(cmd *MascotCmd, wg *sync.WaitGroup, errs chan<- error) {
	streamer, format, err := wav.Decode(bytes.NewReader(mascotAudio))
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
		[]mascotFunc{mascotAudioRun},
		mascotRun...,
	)
}
