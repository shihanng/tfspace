// Package varfile is the varfile subcommand.
package varfile

import "github.com/spf13/cobra"

// NewCommand creates new "varfile" subcommand.
func NewCommand() *cobra.Command {
	varfileCmd := &cobra.Command{ //nolint:exhaustruct
		Use:   "varfile",
		Short: "Add/remove Terraform's varfile from tfspace",
		Long:  "Add/remove Terraform's varfile from tfspace",
	}

	varfileCmd.AddCommand(newAddCommand())
	varfileCmd.AddCommand(newRmCommand())

	return varfileCmd
}
