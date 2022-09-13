package flag

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Bool defines a bool flag with specified name, default value, and usage string.
// The return value is the address of a bool variable that stores the value of the flag.
func Bool(pf *pflag.FlagSet, name string, value bool, usage string) {
	pf.Bool(name, value, usage)
	_ = viper.BindPFlag(name, pf.Lookup(name))
}

// StringP is like String, but accepts a shorthand letter that can be used after a single dash.
func StringP(pf *pflag.FlagSet, name, shorthand string, value string, usage string) {
	pf.StringP(name, shorthand, value, usage)
	_ = viper.BindPFlag(name, pf.Lookup(name))
}
