// Package cmd contains (sub)commands of tfspace.
package cmd

import (
	"github.com/spf13/cobra"
)

// Execute is the entrypoint to tfspace root command.
func Execute(options ...func(*cobra.Command)) {
	rootCmd := &cobra.Command{ //nolint:exhaustruct
		SilenceUsage: true,
		Use:          "tfspace",
		Short:        "Manage multiple environments in a Terraform project with ease.",
		Long:         "Manage multiple environments in a Terraform project with ease.",
	}

	for _, option := range options {
		option(rootCmd)
	}

	cobra.CheckErr(rootCmd.Execute())
}

// WithArgs pass arguments to root command. This is for testing purpose.
func WithArgs(args ...string) func(*cobra.Command) {
	return func(c *cobra.Command) {
		c.SetArgs(args)
	}
}
