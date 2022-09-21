// Package workspace is the workspace subcommand.
package workspace

import "github.com/spf13/cobra"

// NewCommand creates new "workspace" subcommand.
func NewCommand() *cobra.Command {
	workspaceCmd := &cobra.Command{ //nolint:exhaustruct
		Use:   "workspace",
		Short: "Add/remove Terraform's workspace to tfspace",
		Long:  "Add/remove Terraform's workspace to tfspace",
	}

	workspaceCmd.AddCommand(newAddCommand())
	workspaceCmd.AddCommand(newRmCommand())

	return workspaceCmd
}
