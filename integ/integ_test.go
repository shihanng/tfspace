package integ_test

import (
	"bytes"
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
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
	"gotest.tools/v3/env"
	"gotest.tools/v3/fs"
	"gotest.tools/v3/golden"
	"gotest.tools/v3/icmd"
)

const tmpDirPrefix = "tfspace_integtest"

var path = flag.String("path", "../tfspace", "path to the tfspace binary") //nolint:gochecknoglobals

func TestFeatures(t *testing.T) {
	t.Parallel()

	// These are default values. We can override these with flags.
	opts := godog.Options{ //nolint:exhaustruct
		Output:    colors.Colored(os.Stdout),
		Format:    "pretty",
		Randomize: time.Now().UTC().UnixNano(), // randomize scenario execution order
	}

	flag.Parse()

	binPath, err := filepath.Abs(*path)
	if err != nil {
		t.Fatal(err)
	}

	suite := godog.TestSuite{ //nolint:exhaustruct
		Name:                "tfspace",
		ScenarioInitializer: InitializeScenario(binPath),
		Options:             &opts,
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
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

func (s *stepDefinition) terraformerRunsAndThenEnv(ctx context.Context, args string) (context.Context, error) {
	buf := bytes.NewBufferString("env")
	cmd := icmd.Command(s.binPath, strings.Fields(args)...)
	res := icmd.RunCmd(cmd, icmd.WithStdin(buf))

	return withcCmdResultCtx(ctx, res), nil
}

func (s *stepDefinition) tfspaceShouldRunWithoutError(ctx context.Context) error {
	result, err := cmdResult(ctx)
	if err != nil {
		return err
	}

	return s.assertWith(func(a *T) {
		result.Assert(a, icmd.Expected{ExitCode: 0}) //nolint:exhaustruct
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
		result.Assert(a, icmd.Expected{ExitCode: exitCode}) //nolint:exhaustruct
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
		return errors.Wrap(err, "integ: read tfspace file")
	}

	return s.assertWith(func(a *T) {
		assert.Equal(a, string(actual), expected.Content)
	})
}

func (s *stepDefinition) shouldSetEnvironmentVariables(ctx context.Context, table *godog.Table) error {
	result, err := cmdResult(ctx)
	if err != nil {
		return err
	}

	return s.assertWith(func(a *T) {
		for _, row := range table.Rows {
			envRow := fmt.Sprintf("%s=%s", row.Cells[0].Value, row.Cells[1].Value)
			assert.Assert(a, cmp.Contains(result.String(), envRow))
		}
	})
}

func (s *stepDefinition) assertWith(f func(t *T)) error {
	f(s.t)

	return s.t.err
}

func InitializeScenario(binPath string) func(ctx *godog.ScenarioContext) {
	return func(ctx *godog.ScenarioContext) {
		stepDef := stepDefinition{
			binPath: binPath,
			t:       &T{}, //nolint:exhaustruct
		}

		ctx.After(func(ctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
			stepDef.t.runCleanup()

			return ctx, nil
		})

		ctx.Step(`^Terraformer runs "tfspace ([^"]*)"$`, stepDef.terraformerRuns)
		ctx.Step(`^tfspace should print "([^"]*)" (error|content) on screen$`, stepDef.tfspaceShouldPrintOnScreen)
		ctx.Step(`^a project without tfspace\.yml$`, stepDef.aProjectWithoutTfspaceyml)
		ctx.Step(`^tfspace should run without error$`, stepDef.tfspaceShouldRunWithoutError)
		ctx.Step(`^the tfspace\.yml should contain:$`, stepDef.theTfspaceymlShouldContain)
		ctx.Step(`^should set environment variables:$`, stepDef.shouldSetEnvironmentVariables)
		ctx.Step(`^Terraformer runs "tfspace ([^"]*)" and then env$`, stepDef.terraformerRunsAndThenEnv)
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
	fmt.Println(args...) //nolint:forbidigo
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
