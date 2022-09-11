// Package cmd contains (sub)commands of tfspace.
package cmd

import (
	"io"

	"github.com/shihanng/tfspace/cmd/workspace"
	"github.com/shihanng/tfspace/flag"
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
		SilenceErrors:     true,
		PersistentPreRunE: rootPreRun,
		PersistentPostRun: rootPostRun,

		// Disable completion for now.
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	rootCmd.AddCommand(workspace.NewCommand())

	flag.Bool(rootCmd.PersistentFlags(), "debug", false, "emits debug level logs")

	for _, option := range options {
		option(rootCmd)
	}

	return rootCmd.Execute()
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
		return err
	}

	zap.ReplaceGlobals(logger)

	return nil
}

func rootPostRun(_ *cobra.Command, _ []string) {
	_ = zap.L().Sync()
}
