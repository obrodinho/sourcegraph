package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/sourcegraph/sourcegraph/enterprise/cmd/batcheshelper/log"
	"github.com/sourcegraph/sourcegraph/enterprise/cmd/batcheshelper/run"
	"github.com/sourcegraph/sourcegraph/enterprise/cmd/batcheshelper/util"
	batcheslib "github.com/sourcegraph/sourcegraph/lib/batches"
	"github.com/sourcegraph/sourcegraph/lib/batches/execution"
	"github.com/sourcegraph/sourcegraph/lib/errors"
)

func main() {
	if err := doMain(); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func doMain() error {
	inputPath := flag.String("input", "input.json", "The input JSON file for the workspace execution. Defaults to \"input.json\".")
	previousPath := flag.String("previousStepPath", "", "The path to the previous step's result file. Defaults to current working directory.")
	workspaceFilesPath := flag.String("workspaceFiles", "/data/workspace-files", "The path to the workspace files. Defaults to \"/data/workspace-files\".")
	flag.Usage = usage
	flag.Parse()

	arguments, err := parseArgs(os.Args[1:])
	if err != nil {
		return err
	}

	executionInput, err := parseInput(*inputPath)
	if err != nil {
		return err
	}

	previousResult, err := parsePreviousStepResult(*previousPath, arguments.step, executionInput.SkippedSteps)
	if err != nil {
		return err
	}

	logger := &log.Logger{Writer: os.Stdout}

	wd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "getting working directory")
	}
	ctx := context.Background()
	switch arguments.mode {
	case "pre":
		return run.Pre(ctx, logger, arguments.step, executionInput, previousResult, wd, *workspaceFilesPath)
	case "post":
		return run.Post(ctx, logger, &util.RealCmdRunner{}, arguments.step, executionInput, previousResult, wd, *workspaceFilesPath)
	default:
		return errors.Newf("invalid mode %q", arguments.mode)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <pre|post> <step index> [OPTIONS]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "OPTIONS:\n")
	flag.PrintDefaults()
}

func parseArgs(arguments []string) (args, error) {
	if len(arguments) < 2 {
		return args{}, errors.New("missing arguments")
	} else if len(arguments) > 2 {
		return args{}, errors.New("too many arguments")
	}

	mode := arguments[0]
	if mode != "pre" && mode != "post" {
		return args{}, errors.Newf("invalid mode %q", mode)
	}

	step, err := strconv.Atoi(arguments[1])
	if err != nil {
		return args{}, errors.Wrap(err, "failed to parse step")
	}

	return args{mode, step}, nil
}

type args struct {
	mode string
	step int
}

func parseInput(inputPath string) (batcheslib.WorkspacesExecutionInput, error) {
	var executionInput batcheslib.WorkspacesExecutionInput

	input, err := os.ReadFile(inputPath)
	if err != nil {
		return executionInput, errors.Wrapf(err, "failed to read execution input file %q", inputPath)
	}

	if err = json.Unmarshal(input, &executionInput); err != nil {
		return executionInput, errors.Wrap(err, "failed to unmarshal execution input")
	}
	return executionInput, nil
}

func parsePreviousStepResult(path string, step int, skippedSteps map[int]struct{}) (execution.AfterStepResult, error) {
	var previousResult execution.AfterStepResult
	if step > 0 {
		// Get the last actually executed step index.
		previousStepIndex := getPreviousStepIndex(step, skippedSteps)
		if previousStepIndex == -1 {
			// No previous step was executed.
			return previousResult, nil
		}

		// Read the previous step's result file.
		stepResultPath := filepath.Join(path, util.StepJSONFile(previousStepIndex))
		stepJSON, err := os.ReadFile(stepResultPath)
		if err != nil {
			return previousResult, errors.Wrap(err, "failed to read step result file")
		}
		if err = json.Unmarshal(stepJSON, &previousResult); err != nil {
			return previousResult, errors.Wrap(err, "failed to unmarshal step result file")
		}
	}
	return previousResult, nil
}

func getPreviousStepIndex(step int, skippedSteps map[int]struct{}) int {
	if skippedSteps == nil {
		return step - 1
	}
	for i := step - 1; i >= 0; i-- {
		if _, ok := skippedSteps[i]; !ok {
			return i
		}
	}
	return -1
}
