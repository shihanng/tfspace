package flag

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Flager interface {
	PersistentFlags() *pflag.FlagSet
	Flags() *pflag.FlagSet
}

func Bool(pf *pflag.FlagSet, name string, value bool, usage string) {
	pf.Bool(name, value, usage)
	viper.BindPFlag(name, pf.Lookup(name))
}
