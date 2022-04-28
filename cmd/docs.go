package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"path/filepath"
)

func newDocsCmd(rootCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "docs",
		Args:   cobra.NoArgs,
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return doc.GenMarkdownTree(rootCmd, filepath.Join("docs", "commands"))
		},
	}

	return cmd
}
