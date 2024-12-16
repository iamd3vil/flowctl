package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/cvhariharan/autopilot/internal/models"
	"github.com/cvhariharan/autopilot/internal/runner"
	"github.com/expr-lang/expr"
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

type ActionExecutionPayload struct {
	Action          models.Action
	Input           map[string]interface{}
	PreviousOutputs map[string]interface{}
	LogID           string
}

func NewFlowExecution(f models.Flow, input map[string]interface{}, logID string) (*asynq.Task, error) {
	payload, err := json.Marshal(FlowExecutionPayload{Workflow: f, Input: input, LogID: logID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeFlowExecution, payload), nil
}

func NewActionExecution(action models.Action, input map[string]interface{}, logID string) (*asynq.Task, error) {
	payload, err := json.Marshal(ActionExecutionPayload{Action: action, Input: input, LogID: logID})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeActionExecution, payload), nil
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
			Stdout:        streamLogger,
			Stderr:        streamLogger,
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

		// Send a checkpoint to the stream
		if err := streamLogger.Checkpoint(action.ID); err != nil {
			return fmt.Errorf("could not save checkpoint for action %s: %w", action.ID, err)
		}

	}

	return nil
}

func (r *FlowRunner) HandleActionExecution(ctx context.Context, t *asynq.Task) error {
	var payload ActionExecutionPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	// pattern to extract interpolated variables
	pattern := `{{\s*([^}]+)\s*}}`
	re := regexp.MustCompile(pattern)

	streamLogger := r.logger.WithID(payload.LogID)
	defer streamLogger.Close()

	// Iterate over all the flow variables execute variable interpolation if required
	for i, variable := range payload.Action.Variables {
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

			payload.Action.Variables[i][payload.Action.Variables[i].Name()] = output
			log.Println(payload.Action.Variables[i])
		}
	}

	jobCtx, cancel := context.WithTimeout(ctx, 1*time.Hour)
	defer cancel()

	err := runner.NewDockerRunner(payload.Action.ID, r.artifactManager, runner.DockerRunnerOptions{
		ShowImagePull: true,
		Stdout:        streamLogger,
		Stderr:        streamLogger,
	}).CreatesArtifacts(payload.Action.Artifacts).
		WithImage(payload.Action.Image).
		WithCmd(payload.Action.Script).
		WithEnv(payload.Action.Variables).
		WithEntrypoint(payload.Action.Entrypoint).
		WithSrc(payload.Action.Src).
		Run(jobCtx)
	if err != nil {
		return fmt.Errorf("failed to run docker runner: %w", err)
	}

	return nil
}
