package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/fs"
	"sync"
)

type MascotCmd struct {
	Assets fs.ReadFileFS
}

func (cmd *MascotCmd) Command() *cobra.Command {
	command := &cobra.Command{
		Use:    "duck",
		Hidden: true,
		RunE: func(command *cobra.Command, args []string) error {
			return cmd.Run()
		},
	}

	return command
}

func (cmd *MascotCmd) Run() error {
	var wg sync.WaitGroup
	errs := make(chan error)

	for _, run := range mascotRun {
		// Increment the WaitGroup counter
		wg.Add(1)
		// Run
		go run(cmd, &wg, errs)
	}

	// Wait for all runs to complete
	go func() {
		wg.Wait()
		close(errs)
	}()

	// Handle errors
	for err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}

type mascotFunc = func(cmd *MascotCmd, wg *sync.WaitGroup, errs chan<- error)

func mascotRunText(cmd *MascotCmd, wg *sync.WaitGroup, errs chan<- error) {
	text, err := cmd.Assets.ReadFile("assets/mascot.txt")
	if err != nil {
		errs <- err
		return
	}

	fmt.Println(string(text))

	wg.Done()
}

var mascotRun = []mascotFunc{
	mascotRunText,
}
