// Package backend is the backend subcommand.
package backend

import "github.com/spf13/cobra"

// NewCommand creates new "backend" subcommand.
func NewCommand() *cobra.Command {
	backendCmd := &cobra.Command{ //nolint:exhaustruct
		Use:   "backend",
		Short: "Add/remove Terraform's backend from tfspace",
		Long:  "Add/remove Terraform's backend from tfspace",
	}

	backendCmd.AddCommand(newAddCommand())
	backendCmd.AddCommand(newRmCommand())

	return backendCmd
}
