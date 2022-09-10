package main

import (
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/spf13/pflag" // godog v0.11.0 and later
)

var opts = godog.Options{ //nolint:exhaustruct,gochecknoglobals
	Output: colors.Colored(os.Stdout),
	Format: "progress", // can define default values
}

func TestMain(m *testing.M) {
	godog.BindCommandLineFlags("godog.", &opts)

	pflag.Parse()
	opts.Paths = pflag.Args()

	status := godog.TestSuite{ //nolint:exhaustruct
		Name: "godogs",
		// TestSuiteInitializer: InitializeTestSuite,
		// ScenarioInitializer:  InitializeScenario,
		Options: &opts,
	}.Run()

	// Optional: Run `testing` package's logic besides godog.
	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}
