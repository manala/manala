package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

type DocsCmd struct{}

func (cmd *DocsCmd) Command(rootCommand *cobra.Command, dir string) *cobra.Command {
	command := &cobra.Command{
		Use:    "docs",
		Hidden: true,
		Args:   cobra.NoArgs,
		RunE: func(command *cobra.Command, args []string) error {
			return doc.GenMarkdownTree(rootCommand, dir)
		},
	}

	return command
}
