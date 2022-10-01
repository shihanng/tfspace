package workspace

import (
	cmdspace "github.com/shihanng/tfspace/cmd/space"
	"github.com/shihanng/tfspace/space"
	"github.com/spf13/cobra"
)

func newRmCommand() *cobra.Command {
	rmCmd := &cobra.Command{ //nolint:exhaustruct
		Use:           "rm <space>",
		Short:         "Remove Terraform's workspace from tfspace",
		Long:          "Remove Terraform's workspace from tfspace's <space>",
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          runRm,
	}

	return rmCmd
}

func runRm(_ *cobra.Command, args []string) error {
	err := cmdspace.WithSpace(func(s *space.Spaces) {
		s.UnsetWorkspace(args[0])
	})

	return err
}
