// Package cmd contains (sub)commands of tfspace.
package cmd

import (
	"io"
	"os"
	"os/exec"

	"github.com/cockroachdb/errors"
	"github.com/shihanng/tfspace/cmd/backend"
	cmdspace "github.com/shihanng/tfspace/cmd/space"
	"github.com/shihanng/tfspace/cmd/varfile"
	"github.com/shihanng/tfspace/cmd/workspace"
	"github.com/shihanng/tfspace/config"
	"github.com/shihanng/tfspace/space"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/twpayne/go-shell"
	"go.uber.org/zap"
)

// Execute is the entrypoint to tfspace root command.
func Execute(options ...func(*cobra.Command)) error {
	viper.SetEnvPrefix("TFSPACE")
	viper.AutomaticEnv()

	rootCmd := &cobra.Command{ //nolint:exhaustruct
		Use:               "tfspace",
		Short:             "Manage multiple environments in a Terraform project with ease.",
		Long:              "Manage multiple environments in a Terraform project with ease.",
		Args:              cobra.ExactArgs(1),
		SilenceUsage:      true,
		SilenceErrors:     true,
		PersistentPreRunE: rootPreRun,
		PersistentPostRun: rootPostRun,
		RunE:              runRoot,

		// Disable completion for now.
		CompletionOptions: cobra.CompletionOptions{ //nolint:exhaustruct
			DisableDefaultCmd: true,
		},
	}

	rootCmd.AddCommand(workspace.NewCommand())
	rootCmd.AddCommand(backend.NewCommand())
	rootCmd.AddCommand(varfile.NewCommand())

	rootCmd.PersistentFlags().Bool("debug", false, "emits debug level logs")

	config.WithConfig(rootCmd)

	for _, option := range options {
		option(rootCmd)
	}

	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		return errors.Wrap(err, "cmd: fail to bind persistent flags")
	}

	return rootCmd.Execute() //nolint:wrapcheck
}

// WithArgs pass arguments to root command. This is for testing purpose.
func WithArgs(args ...string) func(*cobra.Command) {
	return func(c *cobra.Command) {
		c.SetArgs(args)
	}
}

// WithOutErr sets Stdout and Stderr output to out. This is for testing purpose.
func WithOutErr(out io.Writer) func(*cobra.Command) {
	return func(c *cobra.Command) {
		c.SetOut(out)
		c.SetErr(out)
	}
}

func rootPreRun(_ *cobra.Command, _ []string) error {
	// Setup global logger that can be access via zap.L() or zap.S().
	isDebug := viper.GetBool("debug")

	var (
		logger *zap.Logger
		err    error
	)

	if isDebug {
		logger, err = zap.NewDevelopment()
		logger.Debug("In development mode")
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		return errors.Wrap(err, "root: fail to set logger")
	}

	zap.ReplaceGlobals(logger)

	return nil
}

func runRoot(_ *cobra.Command, args []string) error {
	logger := zap.L()

	shell, found := shell.CurrentUserShell()

	logger = logger.With(zap.String("shell", shell))

	if !found {
		logger.Debug("Failed to get user shell")
	}

	cmd := exec.Command(shell) //nolint:gosec
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmdspace.WithSpace(func(s *space.Spaces) error {
		env, err := s.Env(args[0])
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

func rootPostRun(_ *cobra.Command, _ []string) {
	_ = zap.L().Sync()
}
