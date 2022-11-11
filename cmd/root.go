// Package cmd contains (sub)commands of tfspace.
package cmd

import (
	"io"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/shihanng/tfspace/cmd/backend"
	"github.com/shihanng/tfspace/cmd/use"
	"github.com/shihanng/tfspace/cmd/varfile"
	"github.com/shihanng/tfspace/cmd/workspace"
	"github.com/shihanng/tfspace/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

		// Disable completion for now.
		CompletionOptions: cobra.CompletionOptions{ //nolint:exhaustruct
			DisableDefaultCmd: true,
		},
	}

	rootCmd.AddCommand(workspace.NewCommand())
	rootCmd.AddCommand(backend.NewCommand())
	rootCmd.AddCommand(varfile.NewCommand())
	rootCmd.AddCommand(use.NewCommand())

	rootCmd.PersistentFlags().Bool("debug", false, "emits debug level logs")

	config.WithConfig(rootCmd)

	for _, option := range options {
		option(rootCmd)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	if err := BindPFlags(rootCmd); err != nil {
		return err
	}

	return rootCmd.Execute() //nolint:wrapcheck
}

func BindPFlags(cmd *cobra.Command) error {
	for _, c := range cmd.Commands() {
		if err := BindPFlags(c); err != nil {
			return err
		}
	}

	if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
		return errors.Wrap(err, "cmd: fail to bind persistent flags")
	}

	return errors.Wrap(viper.BindPFlags(cmd.Flags()), "cmd: fail to bind flags")
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

func rootPostRun(_ *cobra.Command, _ []string) {
	_ = zap.L().Sync()
}
