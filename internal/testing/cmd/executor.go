package cmd

import (
	"bytes"
	"github.com/spf13/cobra"
)

func NewExecutor(provider func(stdout *bytes.Buffer, stderr *bytes.Buffer) *cobra.Command) *Executor {
	return &Executor{
		provider: provider,
		Stdout:   &bytes.Buffer{},
		Stderr:   &bytes.Buffer{},
	}
}

type Executor struct {
	provider func(stdout *bytes.Buffer, stderr *bytes.Buffer) *cobra.Command
	Stdout   *bytes.Buffer
	Stderr   *bytes.Buffer
}

func (executor *Executor) Execute(args []string) error {
	executor.Stdout.Reset()
	executor.Stderr.Reset()

	cmd := executor.provider(executor.Stdout, executor.Stderr)
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	cmd.SetOut(executor.Stdout)
	cmd.SetErr(executor.Stderr)

	cmd.SetArgs(args)

	return cmd.Execute()
}
