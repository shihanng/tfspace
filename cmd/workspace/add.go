package workspace

import (
	"github.com/spf13/cobra"
)

func newAddCommand() *cobra.Command {
	addCmd := &cobra.Command{ //nolint:exhaustruct
		Use:           "add <space> <value>",
		Short:         "Add Terraform's workspace to tfspace",
		Long:          "Add Terraform's workspace <value> to tfspace's <space>",
		Args:          cobra.ExactArgs(2),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          runAdd,
	}

	return addCmd
}

func runAdd(_ *cobra.Command, _ []string) error {
	return nil
}
