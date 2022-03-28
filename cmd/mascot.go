package cmd

import (
	_ "embed"
	"fmt"
	"github.com/spf13/cobra"
	"sync"
)

//go:embed embed/mascot.txt
var mascotText string

type MascotCmd struct{}

func (cmd *MascotCmd) Command() *cobra.Command {
	command := &cobra.Command{
		Use:    "duck",
		Hidden: true,
		RunE: func(command *cobra.Command, args []string) error {
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
		},
	}

	return command
}

type mascotFunc = func(cmd *MascotCmd, wg *sync.WaitGroup, errs chan<- error)

func mascotTextRun(cmd *MascotCmd, wg *sync.WaitGroup, errs chan<- error) {
	fmt.Println(mascotText)
	wg.Done()
}

var mascotRun = []mascotFunc{
	mascotTextRun,
}
