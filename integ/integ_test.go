package integ_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/shihanng/tfspace/cmd"
	"github.com/spf13/pflag"
	"gotest.tools/v3/golden"
)

func TestMain(m *testing.M) {
	// These are default values. We can override these with flags.
	opts := godog.Options{ //nolint:exhaustruct
		Output:    colors.Colored(os.Stdout),
		Format:    "pretty",
		Randomize: time.Now().UTC().UnixNano(), // randomize scenario execution order
	}

	godog.BindCommandLineFlags("godog.", &opts)

	pflag.Parse()
	opts.Paths = pflag.Args()

	status := godog.TestSuite{ //nolint:exhaustruct
		Name:                "tfspace",
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}.Run()

	// Optional: Run `testing` package's logic besides godog.
	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}

func terraformerRuns(ctx context.Context, args string) error {
	out, err := getOutput(ctx)
	if err != nil {
		return err
	}

	cmd.Execute(cmd.WithArgs(strings.Fields(args)...), cmd.WithOutErr(out))
	return nil
}

func tfspaceShouldPrintContentOnScreen(ctx context.Context, filename string) error {
	out, err := getOutput(ctx)
	if err != nil {
		return err
	}

	return assertWith(func(a *T) {
		golden.Assert(a, out.String(), filename+".txt")
	})
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		return context.WithValue(ctx, outputCtxKey{}, new(bytes.Buffer)), nil
	})

	ctx.Step(`^Terraformer runs "([^"]*)"$`, terraformerRuns)
	ctx.Step(`^tfspace should print "([^"]*)" content on screen$`, tfspaceShouldPrintContentOnScreen)
}

type outputCtxKey struct{}

func getOutput(ctx context.Context) (*bytes.Buffer, error) {
	out, ok := ctx.Value(outputCtxKey{}).(*bytes.Buffer)
	if !ok {
		return nil, errors.New("bytes.Buffer not found in context")
	}
	return out, nil
}

type T struct {
	err error
}

func (t *T) Log(args ...interface{}) {
	t.err = errors.New(fmt.Sprintln(args...))
}

func (t *T) FailNow() {}

func (t *T) Fail() {}

func assertWith(f func(t *T)) error {
	var t T
	f(&t)
	return t.err
}
