package integ_test

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/spf13/pflag"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/env"
	"gotest.tools/v3/fs"
	"gotest.tools/v3/golden"
	"gotest.tools/v3/icmd"
)

const tmpDirPrefix = "tfspace_integtest"

func TestMain(m *testing.M) {
	// These are default values. We can override these with flags.
	opts := godog.Options{ //nolint:exhaustruct
		Output:    colors.Colored(os.Stdout),
		Format:    "pretty",
		Randomize: time.Now().UTC().UnixNano(), // randomize scenario execution order
	}

	godog.BindCommandLineFlags("godog.", &opts)

	path := flag.String("path", "../tfspace", "path to the tfspace binary")
	binPath, err := filepath.Abs(*path)
	if err != nil {
		panic(err)
	}

	flag.Parse()
	pflag.Parse()
	opts.Paths = pflag.Args()

	status := godog.TestSuite{ //nolint:exhaustruct
		Name:                "tfspace",
		ScenarioInitializer: InitializeScenario(binPath),
		Options:             &opts,
	}.Run()

	// Optional: Run `testing` package's logic besides godog.
	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}

type stepDefinition struct {
	binPath string
	t       *T
}

func (s *stepDefinition) terraformerRuns(ctx context.Context, args string) (context.Context, error) {
	cmd := icmd.Command(s.binPath, strings.Fields(args)...)
	res := icmd.RunCmd(cmd)

	return withcCmdResultCtx(ctx, res), nil
}

func (s *stepDefinition) tfspaceShouldRunWithoutError(ctx context.Context) error {
	result, err := cmdResult(ctx)
	if err != nil {
		return err
	}

	return s.assertWith(func(a *T) {
		result.Assert(a, icmd.Expected{ExitCode: 0})
	})
}

func (s *stepDefinition) tfspaceShouldPrintOnScreen(ctx context.Context, filename, resultType string) error {
	result, err := cmdResult(ctx)
	if err != nil {
		return err
	}

	var (
		output   string
		exitCode int
	)

	switch resultType {
	case "content":
		exitCode = 0
		output = result.Stdout()
	case "error":
		exitCode = 1
		output = result.Stderr()
	}

	return s.assertWith(func(a *T) {
		result.Assert(a, icmd.Expected{ExitCode: exitCode})
		golden.Assert(a, output, normalizeFilename(filename))
	})
}

func (s *stepDefinition) aProjectWithoutTfspaceyml(ctx context.Context) (context.Context, error) {
	if err := s.assertWith(func(a *T) {
		dir := fs.NewDir(a, tmpDirPrefix)
		env.ChangeWorkingDir(a, dir.Path())
	}); err != nil {
		return ctx, err
	}

	return ctx, nil
}

func (s *stepDefinition) theTfspaceymlShouldContain(expected *godog.DocString) error {
	actual, err := os.ReadFile("./tfspace.yml")
	if err != nil {
		return err
	}

	return s.assertWith(func(a *T) {
		assert.Equal(a, string(actual), expected.Content)
	})
}

func (s *stepDefinition) assertWith(f func(t *T)) error {
	f(s.t)
	return s.t.err
}

func InitializeScenario(binPath string) func(ctx *godog.ScenarioContext) {
	return func(ctx *godog.ScenarioContext) {
		sd := stepDefinition{
			binPath: binPath,
			t:       &T{},
		}

		ctx.After(func(ctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
			sd.t.runCleanup()
			return ctx, nil
		})

		ctx.Step(`^Terraformer runs "tfspace ([^"]*)"$`, sd.terraformerRuns)
		ctx.Step(`^tfspace should print "([^"]*)" (error|content) on screen$`, sd.tfspaceShouldPrintOnScreen)
		ctx.Step(`^a project without tfspace\.yml$`, sd.aProjectWithoutTfspaceyml)
		ctx.Step(`^tfspace should run without error$`, sd.tfspaceShouldRunWithoutError)
		ctx.Step(`^the tfspace\.yml should contain:$`, sd.theTfspaceymlShouldContain)
	}
}

type cmdResultCtxKey struct{}

func withcCmdResultCtx(ctx context.Context, result *icmd.Result) context.Context {
	return context.WithValue(ctx, cmdResultCtxKey{}, result)
}

func cmdResult(ctx context.Context) (*icmd.Result, error) {
	result, ok := ctx.Value(cmdResultCtxKey{}).(*icmd.Result)
	if !ok {
		return nil, errors.New("cmdResult not in context")
	}
	return result, nil
}

type T struct {
	err          error
	cleanupFuncs []func()
}

func (t *T) Log(args ...interface{}) {
	fmt.Println(args...)
}

func (t *T) FailNow() {
	t.err = errors.New("integ_test: fail now")
}

func (t *T) Fail() {
	t.err = errors.New("integ_test: fail")
}

func (t *T) Cleanup(f func()) {
	t.cleanupFuncs = append(t.cleanupFuncs, f)
}

func (t *T) runCleanup() {
	for _, f := range t.cleanupFuncs {
		defer f()
	}
	t.cleanupFuncs = nil
}

func normalizeFilename(filename string) string {
	return strings.ReplaceAll(filename, " ", "_") + ".txt"
}
