package cmd

import (
	"bytes"
	"github.com/spf13/cobra"
)

func newCmdExecutor(provider func(stderr *bytes.Buffer) *cobra.Command) *cmdExecutor {
	return &cmdExecutor{
		provider: provider,
		stdout:   &bytes.Buffer{},
		stderr:   &bytes.Buffer{},
	}
}

type cmdExecutor struct {
	provider func(stderr *bytes.Buffer) *cobra.Command
	stdout   *bytes.Buffer
	stderr   *bytes.Buffer
}

func (executor *cmdExecutor) execute(args []string) error {
	executor.stdout.Reset()
	executor.stderr.Reset()

	cmd := executor.provider(executor.stderr)
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	cmd.SetOut(executor.stdout)
	cmd.SetErr(executor.stderr)

	cmd.SetArgs(args)

	return cmd.Execute()
}
