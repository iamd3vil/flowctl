package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/cvhariharan/autopilot/internal/flow"
	"github.com/cvhariharan/autopilot/internal/runner"
	"github.com/expr-lang/expr"
	"github.com/hibiken/asynq"
)

const (
	TypeFlowExecution = "flow_execution"
)

type FlowExecutionPayload struct {
	Workflow flow.Flow
	Input    map[string]interface{}
}

func NewFlowExecution(f flow.Flow, input map[string]interface{}) (*asynq.Task, error) {
	payload, err := json.Marshal(FlowExecutionPayload{Workflow: f, Input: input})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask("flow_execution", payload), nil
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

		err := runner.NewDockerRunner(action.ID, r.artifactManager, runner.DockerRunnerOptions{
			ShowImagePull: true,
			Stdout:        r.logger.WithID(action.ID),
			Stderr:        r.logger.WithID(action.ID),
		}).CreatesArtifacts(action.Artifacts).
			WithImage(action.Image).
			WithCmd(action.Script).
			WithEnv(action.Variables).
			WithEntrypoint(action.Entrypoint).
			WithSrc(action.Src).
			Run(jobCtx)
		if err != nil {
			return fmt.Errorf("failed to run docker runner: %w", err)
		}
	}

	return nil
}
