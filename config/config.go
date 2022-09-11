package config

import (
	"github.com/shihanng/tfspace/flag"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Config is the configuration for the whole tfspace.
type Config struct {
	Path string `mapstructure:"config"`
}

// WithConfig add Config related flags to the command.
func WithConfig(cmd *cobra.Command) {
	flag.StringP(cmd.PersistentFlags(), "config", "c", "./tfspace.yml", "path to tfspace.yml file")
}

// GetConfig collects all values from flags and return Config with those values.
func GetConfig() (*Config, error) {
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
