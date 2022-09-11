// Package main is the entrypoint to tfspace.
package main

import (
	"github.com/shihanng/tfspace/cmd"
	"github.com/spf13/cobra"
)

func main() {
	cobra.CheckErr(cmd.Execute())
}
