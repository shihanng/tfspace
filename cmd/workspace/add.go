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

func runAdd(_ *cobra.Command, args []string) error {
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

	spaces.SetWorkspace(args[0], args[1])

	if err := store.Save(cfg.Path, spaces); err != nil {
		return err
	}

	return nil
}
