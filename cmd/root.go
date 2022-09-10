package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{ //nolint:gochecknoglobals
	SilenceUsage: true,
	Use:          "tfspace",
	Short:        "Manage multiple environments in a Terraform project with ease.",
	Long:         "Manage multiple environments in a Terraform project with ease.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
