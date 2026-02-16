package flow

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cvhariharan/flowctl/sdk/executor"
	"github.com/invopop/jsonschema"
	"gopkg.in/yaml.v3"
)

type FlowWithConfig struct {
	FlowID    string `yaml:"flow_id" json:"flow_id" jsonschema:"title=flow_id,description=ID of the flow to execute,required" jsonschema_extras:"placeholder=my-flow-id"`
	Params    string `yaml:"params,omitempty" json:"params,omitempty" jsonschema:"title=params,description=JSON parameters to pass to the flow" jsonschema_extras:"widget=codeeditor"`
	Wait      bool   `yaml:"wait,omitempty" json:"wait,omitempty" jsonschema:"title=wait,description=Wait for the flow execution to complete" jsonschema_extras:"type=checkbox"`
	Namespace string `yaml:"namespace,omitempty" json:"namespace,omitempty" jsonschema:"title=namespace,description=Target namespace (defaults to current namespace)"`
}

// Terminal execution statuses
const (
	statusCompleted = "completed"
	statusErrored   = "errored"
	statusCancelled = "cancelled"
)

type FlowExecutor struct {
	name   string
	execID string
}

func GetSchema() interface{} {
	return jsonschema.Reflect(&FlowWithConfig{})
}

// NewFlowExecutor returns a closure that creates FlowExecutor instances.
func NewFlowExecutor(name string, driver executor.NodeDriver, execID string) (executor.Executor, error) {
	return &FlowExecutor{name: name, execID: execID}, nil
}

func (j *FlowExecutor) GetArtifactsDir() string {
	return ""
}

func (j *FlowExecutor) Execute(ctx context.Context, execCtx executor.ExecutionContext) (map[string]string, error) {
	if execCtx.APIKey == "" || execCtx.APIBaseURL == "" {
		return nil, fmt.Errorf("flow executor %s requires API credentials (APIKey and APIBaseURL)", j.name)
	}

	var config FlowWithConfig
	if err := yaml.Unmarshal(execCtx.WithConfig, &config); err != nil {
		return nil, fmt.Errorf("could not read config for Flow executor %s: %w", j.name, err)
	}

	if config.FlowID == "" {
		return nil, fmt.Errorf("flow_id is required for Flow executor %s", j.name)
	}

	// Parse params JSON into map
	params := make(map[string]interface{})
	if config.Params != "" {
		if err := json.Unmarshal([]byte(config.Params), &params); err != nil {
			return nil, fmt.Errorf("failed to parse params JSON: %w", err)
		}
	}

	// Use configured namespace or default to current namespace
	namespace := config.Namespace
	if namespace == "" {
		namespace = execCtx.NamespaceName
	}

	client := executor.NewAPIClient(execCtx.APIBaseURL, execCtx.APIKey, execCtx.UserUUID)

	// Trigger the child flow via API
	triggerResp, err := client.TriggerFlow(ctx, namespace, config.FlowID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to trigger flow %q: %w", config.FlowID, err)
	}

	fmt.Fprintf(execCtx.Stdout, "queued flow %s with exec_id %s\n", config.FlowID, triggerResp.ExecID)

	outputs := map[string]string{
		"exec_id": triggerResp.ExecID,
	}

	if config.Wait {
		status, err := j.waitForCompletion(ctx, client, namespace, triggerResp.ExecID)
		if err != nil {
			return nil, err
		}
		outputs["status"] = status

		if status == statusErrored {
			return outputs, fmt.Errorf("child flow execution %s errored", triggerResp.ExecID)
		}
		if status == statusCancelled {
			return outputs, fmt.Errorf("child flow execution %s was cancelled", triggerResp.ExecID)
		}

		fmt.Fprintf(execCtx.Stdout, "flow %s completed with status %s\n", config.FlowID, status)
	}

	return outputs, nil
}

func (j *FlowExecutor) waitForCompletion(ctx context.Context, client *executor.APIClient, namespace, execID string) (string, error) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("context cancelled while waiting for execution %s: %w", execID, ctx.Err())
		case <-ticker.C:
			status, err := client.GetFlowStatus(ctx, namespace, execID)
			if err != nil {
				return "", fmt.Errorf("failed to get execution status for %s: %w", execID, err)
			}

			switch status.Status {
			case statusCompleted, statusErrored, statusCancelled:
				return status.Status, nil
			}
		}
	}
}
