// Package main is the entrypoint to tfspace.
package main

import (
	"github.com/shihanng/tfspace/cmd"
	"github.com/spf13/viper"
)

func main() {
	viper.SetEnvPrefix("TFSPACE")
	viper.AutomaticEnv()

	cmd.Execute()
}
