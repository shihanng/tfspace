package backend

import (
	cmdspace "github.com/shihanng/tfspace/cmd/space"
	"github.com/shihanng/tfspace/space"
	"github.com/spf13/cobra"
)

func newAddCommand() *cobra.Command {
	addCmd := &cobra.Command{ //nolint:exhaustruct
		Use:           "add <space> <value>",
		Short:         "Add Terraform's backend to tfspace",
		Long:          "Add Terraform's backend <value> to tfspace's <value>",
		Args:          cobra.ExactArgs(2), //nolint:gomnd
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          runAdd,
	}

	return addCmd
}

func runAdd(_ *cobra.Command, args []string) error {
	err := cmdspace.WithSpace(func(s *space.Spaces) {
		s.AddBackend(args[0], args[1])
	})

	return err
}
