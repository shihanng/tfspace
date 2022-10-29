package backend

import (
	cmdspace "github.com/shihanng/tfspace/cmd/space"
	"github.com/shihanng/tfspace/space"
	"github.com/spf13/cobra"
)

func newRmCommand() *cobra.Command {
	rmCmd := &cobra.Command{ //nolint:exhaustruct
		Use:           "rm <space> <value>",
		Short:         "Remove Terraform's backend from tfspace",
		Long:          "Remove Terraform's backend <value> from tfspace's <space>",
		Args:          cobra.ExactArgs(2), //nolint:gomnd
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          runRm,
	}

	return rmCmd
}

func runRm(_ *cobra.Command, args []string) error {
	err := cmdspace.WithSpace(func(s *space.Spaces) error {
		s.RemoveBackend(args[0], args[1])

		return nil
	})

	return err
}
