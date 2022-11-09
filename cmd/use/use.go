// Package use is the use subcommand.
package use

import (
	"os"
	"os/exec"

	"github.com/cockroachdb/errors"
	cmdspace "github.com/shihanng/tfspace/cmd/space"
	"github.com/shihanng/tfspace/space"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/twpayne/go-shell"
	"go.uber.org/zap"
)

// NewCommand creates new "use" subcommand.
func NewCommand() *cobra.Command {
	useCmd := &cobra.Command{ //nolint:exhaustruct
		Use:   "use",
		Short: "Start a new shell with the environment for a specific space.",
		Long:  "Start a new shell with the environment for a specific space.",
		RunE:  useRoot,
	}

	useCmd.Flags().BoolP("with-apply", "a", false, "export TF_CLI_ARGS_apply environment variable")

	if err := viper.BindPFlags(useCmd.Flags()); err != nil {
		return nil
	}

	return useCmd
}

type useConfig struct {
	WithApply bool `mapstructure:"with-apply"`
}

func useRoot(_ *cobra.Command, args []string) error {
	logger := zap.L()

	shell, found := shell.CurrentUserShell()

	logger = logger.With(zap.String("shell", shell))

	if !found {
		logger.Debug("Failed to get user shell")
	}

	var config useConfig
	if err := viper.Unmarshal(&config); err != nil {
		return errors.Wrap(err, "use: fail to get config")
	}

	cmd := exec.Command(shell) //nolint:gosec
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmdspace.WithSpace(func(s *space.Spaces) error {
		env, err := s.Env(args[0], config.WithApply)
		if err != nil {
			return err
		}

		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, env...)

		return nil
	})
	if err != nil {
		return err
	}

	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "cmd: fail to run %s", shell)
	}

	return nil
}
