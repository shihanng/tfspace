// Package cmd contains (sub)commands of tfspace.
package cmd

import (
	"io"

	"github.com/shihanng/tfspace/cmd/workspace"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Execute is the entrypoint to tfspace root command.
func Execute(options ...func(*cobra.Command)) error {
	viper.SetEnvPrefix("TFSPACE")
	viper.AutomaticEnv()

	rootCmd := &cobra.Command{ //nolint:exhaustruct
		SilenceUsage: true,
		Use:          "tfspace",
		Short:        "Manage multiple environments in a Terraform project with ease.",
		Long:         "Manage multiple environments in a Terraform project with ease.",
	}

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
