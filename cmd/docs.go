package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"path/filepath"
)

func newDocsCmd(rootCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "docs",
		Args:   cobra.MaximumNArgs(1),
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get args
			dir := filepath.Clean(append(args, "")[0])

			return doc.GenMarkdownTree(rootCmd, dir)
		},
	}

	return cmd
}
