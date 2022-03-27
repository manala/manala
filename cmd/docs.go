package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

type DocsCmd struct {
	RootCommand *cobra.Command
	Dir         string
}

func (cmd *DocsCmd) Command() *cobra.Command {
	command := &cobra.Command{
		Use:    "docs",
		Hidden: true,
		Args:   cobra.NoArgs,
		RunE: func(command *cobra.Command, args []string) error {
			return doc.GenMarkdownTree(cmd.RootCommand, cmd.Dir)
		},
	}

	return command
}
