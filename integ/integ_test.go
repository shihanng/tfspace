package integ_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/spf13/pflag"
	"gotest.tools/v3/fs"
	"gotest.tools/v3/golden"
	"gotest.tools/v3/icmd"
)

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
}

func (s *stepDefinition) terraformerRuns(ctx context.Context, args string) (context.Context, error) {
	cmd := icmd.Command(s.binPath, strings.Fields(args)...)
	res := icmd.RunCmd(cmd)

	return withcCmdResultCtx(ctx, res), nil
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

	return assertWith(func(a *T) {
		result.Assert(a, icmd.Expected{ExitCode: exitCode})
		golden.Assert(a, output, normalizeFilename(filename))
	})
}

func (s *stepDefinition) aProjectWithoutTfspaceyml(ctx context.Context) (context.Context, error) {
	if err := assertWith(func(a *T) {
		dir := fs.NewDir(a, "new_space")
		ctx = withConfigPath(ctx, dir.Path())
	}); err != nil {
		return ctx, err
	}

	return ctx, nil
}

func InitializeScenario(binPath string) func(ctx *godog.ScenarioContext) {
	return func(ctx *godog.ScenarioContext) {
		sd := stepDefinition{
			binPath: binPath,
		}

		ctx.Step(`^Terraformer runs "tfspace ([^"]*)"$`, sd.terraformerRuns)
		ctx.Step(`^tfspace should print "([^"]*)" (error|content) on screen$`, sd.tfspaceShouldPrintOnScreen)
		ctx.Step(`^a project without tfspace\.yml$`, sd.aProjectWithoutTfspaceyml)
	}
}

type configPathCtxKey struct{}

func withConfigPath(ctx context.Context, path string) context.Context {
	return context.WithValue(ctx, configPathCtxKey{}, path)
}

func getConfigPath(ctx context.Context) (string, error) {
	path, ok := ctx.Value(configPathCtxKey{}).(string)
	if !ok {
		return "", errors.New("config path not found in context")
	}
	return path, nil
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
	err error
}

func (t *T) Log(args ...interface{}) {
	t.err = errors.New(fmt.Sprintln(args...))
}

func (t *T) FailNow() {}

func (t *T) Fail() {}

func (t *T) Cleanup(f func()) {
	f()
}

func assertWith(f func(t *T)) error {
	var t T
	f(&t)
	return t.err
}

func normalizeFilename(filename string) string {
	return strings.ReplaceAll(filename, " ", "_") + ".txt"
}
