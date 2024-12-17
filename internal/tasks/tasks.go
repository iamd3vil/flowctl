package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/cvhariharan/autopilot/internal/models"
	"github.com/cvhariharan/autopilot/internal/runner"
	"github.com/docker/docker/api/types/mount"
	"github.com/expr-lang/expr"
	"github.com/hashicorp/go-envparse"
	"github.com/hibiken/asynq"
)

const (
	TypeFlowExecution   = "flow_execution"
	TypeActionExecution = "action_execution"
)

type FlowExecutionPayload struct {
	Workflow models.Flow
	Input    map[string]interface{}
	LogID    string
}

func NewFlowExecution(f models.Flow, input map[string]interface{}, logID string) (*asynq.Task, error) {
	payload, err := json.Marshal(FlowExecutionPayload{Workflow: f, Input: input, LogID: logID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeFlowExecution, payload), nil
}

type FlowRunner struct {
	logger          *runner.StreamLogger
	artifactManager runner.ArtifactManager
}

func NewFlowRunner(logger *runner.StreamLogger, artifactManager runner.ArtifactManager) *FlowRunner {
	return &FlowRunner{logger: logger, artifactManager: artifactManager}
}

func (r *FlowRunner) HandleFlowExecution(ctx context.Context, t *asynq.Task) error {
	var payload FlowExecutionPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	// pattern to extract interpolated variables
	pattern := `{{\s*([^}]+)\s*}}`
	re := regexp.MustCompile(pattern)

	streamLogger := r.logger.WithID(payload.LogID)
	defer streamLogger.Close()

	// Create temp file for outputs
	outfile, err := os.CreateTemp("", fmt.Sprintf("output-flow-%s-*", payload.Workflow.Meta.ID))
	if err != nil {
		return fmt.Errorf("could not create tmp file for storing flow %s outputs: %w", payload.Workflow.Meta.ID, err)
	}
	log.Println(outfile.Name())

	for _, action := range payload.Workflow.Actions {
		jobCtx, cancel := context.WithTimeout(ctx, time.Hour)
		defer cancel()

		// Iterate over all the flow variables execute variable interpolation if required
		for i, variable := range action.Variables {
			matches := re.FindAllStringSubmatch(variable.Value(), -1)
			if len(matches) > 0 {
				inputExpr := matches[0][1]
				env := map[string]interface{}{
					"input": payload.Input,
				}

				program, err := expr.Compile(inputExpr, expr.Env(env))
				if err != nil {
					return fmt.Errorf("failed to compile expression: %w", err)
				}

				output, err := expr.Run(program, env)
				if err != nil {
					return fmt.Errorf("failed to run expression: %w", err)
				}

				action.Variables[i][action.Variables[i].Name()] = output
				log.Println(action.Variables[i])
			}
		}

		// Add output env variable
		action.Variables = append(action.Variables, map[string]interface{}{"OUTPUT": "/tmp/flow/output"})

		err := runner.NewDockerRunner(action.ID, r.artifactManager, runner.DockerRunnerOptions{
			ShowImagePull: true,
			Stdout:        streamLogger,
			Stderr:        streamLogger,
		}).CreatesArtifacts(action.Artifacts).
			WithImage(action.Image).
			WithCmd(action.Script).
			WithEnv(action.Variables).
			WithEntrypoint(action.Entrypoint).
			WithSrc(action.Src).
			// Output file
			WithMount(mount.Mount{
				Type:   mount.TypeBind,
				Source: outfile.Name(),
				Target: "/tmp/flow/output",
			}).
			Run(jobCtx)
		if err != nil {
			return fmt.Errorf("failed to run docker runner: %w", err)
		}

		// Parse output file env
		outputTempFile, err := os.Open(outfile.Name())
		if err != nil {
			return fmt.Errorf("error opening output file for reading: %w", err)
		}

		outputEnv, err := envparse.Parse(outputTempFile)
		if err != nil {
			return fmt.Errorf("could not load output env: %w", err)
		}

		// Send a checkpoint to the stream
		if err := streamLogger.Checkpoint(models.ExecutionCheckpoint{ActionID: action.ID, Results: outputEnv}); err != nil {
			return fmt.Errorf("could not save checkpoint for action %s: %w", action.ID, err)
		}

	}

	return nil
}
