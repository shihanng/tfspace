package workspace

import (
	"io/fs"

	"github.com/cockroachdb/errors"
	"github.com/shihanng/tfspace/config"
	"github.com/shihanng/tfspace/space"
	"github.com/shihanng/tfspace/store"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
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
	logger := zap.L()

	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}
	logger = logger.With(zap.String("config_path", cfg.Path))

	logger.Debug("Load spaces")
	spaces, err := store.Load(cfg.Path)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}

		logger.Debug("Config does not exist")
		spaces = space.Spaces{}
	}

	spaces.UnsetWorkspace(args[0])

	if err := store.Save(cfg.Path, spaces); err != nil {
		return err
	}

	return nil
}
