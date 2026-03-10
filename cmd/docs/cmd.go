package docs

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func NewCommand(root *cobra.Command) *cobra.Command {
	// Command
	command := &cobra.Command{
		Use:    "docs",
		Args:   cobra.MaximumNArgs(1),
		Hidden: true,
		RunE: func(_ *cobra.Command, args []string) error {
			// Args
			dir := filepath.Clean(append(args, "")[0])

			return run(root, dir)
		},
	}

	return command
}

func run(root *cobra.Command, dir string) error {
	return doc.GenMarkdownTree(root, dir)
}
