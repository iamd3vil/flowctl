package scheduler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/cvhariharan/flowctl/internal/metrics"
	"github.com/cvhariharan/flowctl/internal/repo"
	"github.com/cvhariharan/flowctl/internal/streamlogger"
	"github.com/cvhariharan/flowctl/sdk/executor"
	"github.com/expr-lang/expr"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

const PayloadTypeFlowExecution PayloadType = "flow_execution"

// FlowExecutionHandler handles flow execution jobs
type FlowExecutionHandler struct {
	store           repo.Store
	secretsProvider SecretsProviderFn
	logmanager      streamlogger.LogManager
	logger          *slog.Logger
	metrics         *metrics.Manager
}

// FlowHandlerConfig holds configuration for FlowExecutionHandler
type FlowHandlerConfig struct {
	Store           repo.Store
	SecretsProvider SecretsProviderFn
	LogManager      streamlogger.LogManager
	Logger          *slog.Logger
	Metrics         *metrics.Manager
}

// NewFlowExecutionHandler creates a new flow execution handler
func NewFlowExecutionHandler(cfg FlowHandlerConfig) *FlowExecutionHandler {
	return &FlowExecutionHandler{
		store:           cfg.Store,
		secretsProvider: cfg.SecretsProvider,
		logmanager:      cfg.LogManager,
		logger:          cfg.Logger,
		metrics:         cfg.Metrics,
	}
}

// SetSecretsProvider allows updating secrets provider after creation
func (h *FlowExecutionHandler) SetSecretsProvider(sp SecretsProviderFn) {
	h.secretsProvider = sp
}

// Type returns the payload type this handler processes
func (h *FlowExecutionHandler) Type() PayloadType {
	return PayloadTypeFlowExecution
}

// Handle processes a flow execution job
func (h *FlowExecutionHandler) Handle(ctx context.Context, job Job) error {
	var payload FlowExecutionPayload
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		return fmt.Errorf("failed to unmarshal flow payload: %w", err)
	}

	// Create execution log for scheduled executions only (manual ones are created in core)
	if payload.TriggerType == TriggerTypeScheduled {
		if err := h.createExecutionLog(ctx, job.ExecID, payload); err != nil {
			return fmt.Errorf("failed to create execution log: %w", err)
		}
	}

	// Set status to Running
	if err := h.setStatus(ctx, job.ExecID, repo.ExecutionStatusRunning, payload.NamespaceID, nil); err != nil {
		return fmt.Errorf("could not update execution_log status: %w", err)
	}

	if h.metrics != nil {
		h.metrics.IncExecutionsRunning(payload.NamespaceID, payload.Workflow.Meta.ID)
	}

	// Execute the flow
	if err := h.executeFlow(ctx, job.ExecID, payload); err != nil {
		h.logger.Error("error executing flow", "flow", payload.Workflow.Meta.ID, "error", err)
		if errors.Is(err, ErrPendingApproval) {
			return h.setStatusWithMetrics(ctx, job.ExecID, repo.ExecutionStatusPendingApproval, payload.NamespaceID, payload.Workflow.Meta.ID, nil)
		}
		if errors.Is(err, ErrExecutionCancelled) {
			return h.setStatusWithMetrics(ctx, job.ExecID, repo.ExecutionStatusCancelled, payload.NamespaceID, payload.Workflow.Meta.ID, nil)
		}
		return h.setStatusWithMetrics(ctx, job.ExecID, repo.ExecutionStatusErrored, payload.NamespaceID, payload.Workflow.Meta.ID, err)
	}

	if h.metrics != nil {
		h.metrics.DecExecutionsRunning(payload.NamespaceID, payload.Workflow.Meta.ID)
	}

	return h.setStatusWithMetrics(ctx, job.ExecID, repo.ExecutionStatusCompleted, payload.NamespaceID, payload.Workflow.Meta.ID, nil)
}

// executeFlow executes a flow
func (h *FlowExecutionHandler) executeFlow(ctx context.Context, execID string, payload FlowExecutionPayload) error {
	if payload.StartingActionIdx < 0 {
		payload.StartingActionIdx = 0
	}
	if payload.StartingActionIdx > len(payload.Workflow.Actions) {
		payload.StartingActionIdx = len(payload.Workflow.Actions)
	}

	// Create temporary directory for artifacts shared across all actions in this flow
	artifactDir := filepath.Join(os.TempDir(), fmt.Sprintf("artifacts-store-%s", execID))
	if err := os.MkdirAll(artifactDir, 0700); err != nil {
		return fmt.Errorf("failed to create artifact directory: %w", err)
	}
	h.logger.Debug("artifact directory creation", "dir", artifactDir)

	// Copy files from flow directory to artifacts if flow directory is specified
	if payload.FlowDirectory != "" {
		if err := h.copyFlowFilesToArtifacts(payload.FlowDirectory, artifactDir); err != nil {
			return fmt.Errorf("failed to copy flow files to artifacts: %w", err)
		}
	}

	streamID := execID

	streamLogger, err := h.logmanager.NewLogger(streamID)
	if err != nil {
		return err
	}
	defer streamLogger.Close()

	// Get flow-specific secrets
	flowSecrets := h.getFlowSecrets(ctx, payload.Workflow.Meta.ID, payload.NamespaceID, execID)

	// Initialize outputs map to accumulate results from all previous actions
	outputs := make(map[string]any)

	for i := payload.StartingActionIdx; i < len(payload.Workflow.Actions); i++ {
		action := payload.Workflow.Actions[i]

		res, err := h.executeSingleAction(ctx, action, payload.Workflow.Meta.SrcDir, payload.Input, streamLogger, artifactDir, flowSecrets, outputs, execID, payload.NamespaceID)
		if err != nil {
			return err
		}

		h.logger.Debug("Action results", "results", res)
		processActionResults(res, outputs)
		h.logger.Debug("outputs", "results", outputs)
	}

	// Only remove the artifact store when all actions have been executed
	// This is to account for approval actions that could be run later
	os.RemoveAll(artifactDir)
	return nil
}

// getFlowSecrets retrieves flow-specific secrets or returns an empty map if unavailable
func (h *FlowExecutionHandler) getFlowSecrets(ctx context.Context, flowID string, namespaceID string, execID string) map[string]string {
	if h.secretsProvider == nil {
		return make(map[string]string)
	}

	secrets, err := h.secretsProvider(ctx, flowID, namespaceID)
	if err != nil {
		h.logger.Error("failed to get flow secrets", "execID", execID, "error", err)
		return make(map[string]string)
	}

	return secrets
}

// copyFlowFilesToArtifacts copies top-level files from the flow directory to the artifacts directory
func (h *FlowExecutionHandler) copyFlowFilesToArtifacts(flowDir string, artifactDir string) error {
	localArtifactDir := filepath.Join(artifactDir, "local")
	if err := os.MkdirAll(localArtifactDir, 0755); err != nil {
		return fmt.Errorf("failed to create local artifact directory: %w", err)
	}

	entries, err := os.ReadDir(flowDir)
	if err != nil {
		return fmt.Errorf("failed to read flow directory: %w", err)
	}

	for _, entry := range entries {
		// Skip directories, only copy top-level files
		if entry.IsDir() {
			continue
		}

		srcPath := filepath.Join(flowDir, entry.Name())
		destPath := filepath.Join(localArtifactDir, entry.Name())

		srcFile, err := os.Open(srcPath)
		if err != nil {
			return fmt.Errorf("failed to open source file %s: %w", srcPath, err)
		}
		defer srcFile.Close()

		destFile, err := os.Create(destPath)
		if err != nil {
			return fmt.Errorf("failed to create destination file %s: %w", destPath, err)
		}
		defer destFile.Close()

		if _, err := io.Copy(destFile, srcFile); err != nil {
			return fmt.Errorf("failed to copy file %s to %s: %w", srcPath, destPath, err)
		}

		h.logger.Debug("copied flow file to artifacts", "src", srcPath, "dest", destPath)
	}

	return nil
}

// executeSingleAction executes a single action within a flow, handling approval and error checkpointing
func (h *FlowExecutionHandler) executeSingleAction(ctx context.Context, action Action, srcDir string, input map[string]any, streamLogger streamlogger.Logger, artifactDir string, secrets map[string]string, outputs map[string]any, execID string, namespaceID string) (map[string]string, error) {
	// Check for context cancellation
	if ctx.Err() != nil {
		if err := streamLogger.Checkpoint("", "", "execution cancelled", streamlogger.CancelledMessageType); err != nil {
			h.logger.Error("failed to send cancellation message", "error", err)
		}
		return nil, ErrExecutionCancelled
	}

	// Check for approval requests
	if err := h.checkApproval(ctx, execID, action, namespaceID); err != nil {
		return nil, err
	}

	// Run the action
	res, err := h.runAction(ctx, execID, action, input, streamLogger, artifactDir, secrets, outputs)
	if err != nil {
		// Check if the error is due to context cancellation
		if errors.Is(err, context.Canceled) {
			if streamErr := streamLogger.Checkpoint(action.ID, "", "execution cancelled", streamlogger.CancelledMessageType); streamErr != nil {
				h.logger.Error("failed to send cancelled message", "execID", execID, "actionID", action.ID, "error", streamErr)
			}
			return nil, ErrExecutionCancelled
		}
		streamLogger.Checkpoint(action.ID, "", err.Error(), streamlogger.ErrMessageType)
		return nil, err
	}

	// Checkpoint successful result
	if err := streamLogger.Checkpoint(action.ID, "", res, streamlogger.ResultMessageType); err != nil {
		return nil, err
	}

	return res, nil
}

// processActionResults processes action results and updates the outputs map
func processActionResults(results map[string]string, outputs map[string]any) {
	for k, v := range results {
		parts := strings.SplitN(k, "@", 2)
		// node suffixed output
		if len(parts) == 2 {
			keyName := parts[0]
			nodeName := parts[1]

			if _, exists := outputs[nodeName]; !exists {
				outputs[nodeName] = make(map[string]any)
			}
			outputs[nodeName].(map[string]any)[keyName] = v
		} else {
			outputs[k] = v
		}
	}
}

// executeOnNode executes an action on a single node and returns the results
func (h *FlowExecutionHandler) executeOnNode(ctx context.Context, execID string, node Node, action Action, streamLogger streamlogger.Logger, inputVars map[string]any, withConfig []byte, artifactDir string) ExecResults {
	nodeLogger := streamlogger.NewNodeContextLogger(streamLogger, action.ID, node.Name)

	// Create a separate executor instance for each node
	var exec executor.Executor
	nodeExecutorID := fmt.Sprintf("%s-%s", action.ID, node.Name)
	if node.Name == "" {
		nodeExecutorID = action.ID
	}

	// Check if node is accessible
	// Ignore local node
	if node.Name != "" {
		if err := node.CheckConnectivity(); err != nil {
			h.logger.Debug("node connectivity", "error", err)
			return ExecResults{
				result: nil,
				err:    fmt.Errorf("failed to connect to node %s", node.Name),
			}
		}
	}

	// Convert task node to executor node
	execNode := executor.Node{
		Hostname:       node.Hostname,
		Port:           node.Port,
		Username:       node.Username,
		ConnectionType: node.ConnectionType,
		OSFamily:       node.OSFamily,
		Auth: executor.NodeAuth{
			Method: string(node.Auth.Method),
			Key:    node.Auth.Key,
		},
	}

	driver, err := executor.NewNodeDriver(ctx, execNode)
	if err != nil {
		return ExecResults{
			result: nil,
			err:    fmt.Errorf("failed to create node driver: %w", err),
		}
	}
	defer driver.Close()

	ef, err := executor.GetNewExecutorFunc(action.Executor)
	if err != nil {
		return ExecResults{
			result: nil,
			err:    fmt.Errorf("failed to get executor for %s: %w", action.ID, err),
		}
	}
	exec, err = ef(nodeExecutorID, driver)
	if err != nil {
		return ExecResults{
			result: nil,
			err:    fmt.Errorf("failed to create executor for %s: %w", action.ID, err),
		}
	}

	// Push existing artifacts to this node's executor before execution
	if err := h.pushArtifactsWithDriver(ctx, driver, artifactDir, execID); err != nil {
		return ExecResults{
			result: nil,
			err:    fmt.Errorf("failed to push artifacts to node %s: %w", node.Name, err),
		}
	}

	res, err := exec.Execute(ctx, executor.ExecutionContext{
		ExecID:     execID,
		Inputs:     inputVars,
		WithConfig: withConfig,
		Stdout:     nodeLogger,
		Stderr:     nodeLogger,
	})

	// Pull all artifacts from this node after execution
	if err == nil {
		if pullErr := h.pullArtifactsWithDriver(ctx, driver, artifactDir, execID, node.Name); pullErr != nil {
			err = fmt.Errorf("execution succeeded but failed to pull artifacts: %w", pullErr)
		}
	}

	// Add node.Name suffix to result keys
	prefixedRes := prefixResultKeys(res, node.Name)

	return ExecResults{
		result: prefixedRes,
		err:    err,
	}
}

// prefixResultKeys adds node name suffix to result keys for node-specific outputs
func prefixResultKeys(results map[string]string, nodeName string) map[string]string {
	prefixedRes := make(map[string]string)
	for key, value := range results {
		// Format key as valid environment variable (replace special chars with _)
		prefixedKey := regexp.MustCompile(`[^a-zA-Z0-9_]+`).ReplaceAllString(key, "_")
		if nodeName != "" {
			// example key@hostname
			prefixedKey = prefixedKey + "@" + nodeName
		}
		prefixedRes[prefixedKey] = value
	}
	return prefixedRes
}

// interpolateVariables processes action variables and replaces templated values with evaluated expressions
func (h *FlowExecutionHandler) interpolateVariables(action Action, input map[string]any, secrets map[string]string, outputs map[string]any) (map[string]any, error) {
	// pattern to extract interpolated variables
	pattern := `{{\s*([^}]+)\s*}}`
	re := regexp.MustCompile(pattern)

	h.logger.Debug("scheduler variables", "input", input)

	inputVars := make(map[string]any)
	for _, variable := range action.Variables {
		matches := re.FindAllStringSubmatch(variable.Value(), -1)
		if len(matches) > 0 {
			// Interpolated variable, needs evaluation
			inputExpr := matches[0][1]
			env := map[string]any{
				"inputs":  input,
				"secrets": secrets,
				"outputs": outputs,
			}

			program, err := expr.Compile(inputExpr, expr.Env(env))
			if err != nil {
				return nil, fmt.Errorf("failed to compile expression: %w", err)
			}

			output, err := expr.Run(program, env)
			if err != nil {
				return nil, fmt.Errorf("failed to run expression: %w", err)
			}

			inputVars[variable.Name()] = ""
			if output != nil {
				inputVars[variable.Name()] = output
			}
		} else {
			// Normal variable, no evaluation
			inputVars[variable.Name()] = variable.Value()
		}
	}

	return inputVars, nil
}

// runAction executes a single action
func (h *FlowExecutionHandler) runAction(ctx context.Context, execID string, action Action, input map[string]any, streamLogger streamlogger.Logger, artifactDir string, secrets map[string]string, outputs map[string]any) (map[string]string, error) {
	streamLogger.SetActionID(action.ID)

	jobCtx, cancel := context.WithTimeout(ctx, JobTimeout)
	defer cancel()

	// Interpolate variables
	inputVars, err := h.interpolateVariables(action, input, secrets, outputs)
	if err != nil {
		return nil, err
	}

	withConfig, err := yaml.Marshal(action.With)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal 'with' config: %w", err)
	}

	if len(action.On) == 0 {
		action.On = append(action.On, Node{})
	}

	var wg sync.WaitGroup
	resChan := make(chan ExecResults, len(action.On))

	for _, node := range action.On {
		wg.Add(1)
		go func(node Node) {
			defer wg.Done()
			result := h.executeOnNode(jobCtx, execID, node, action, streamLogger, inputVars, withConfig, artifactDir)
			resChan <- result
		}(node)
	}

	wg.Wait()
	close(resChan)

	// Merge all results into a single map
	mergedResults := make(map[string]string)
	for res := range resChan {
		if res.err != nil {
			// Check if any executor returned a context cancellation error
			if errors.Is(res.err, context.Canceled) {
				return nil, context.Canceled
			}
			return nil, res.err
		}
		maps.Copy(mergedResults, res.result)
	}

	return mergedResults, nil
}

// pushArtifactsWithDriver pushes files from the local artifact directory to the remote artifacts directory
// Only pushes direct child files of top-level directories (one level deep)
func (h *FlowExecutionHandler) pushArtifactsWithDriver(ctx context.Context, driver executor.NodeDriver, artifactDir string, execID string) error {
	remoteArtifactsDir := driver.Join(driver.TempDir(), fmt.Sprintf("artifacts-%s", execID))
	h.logger.Debug("remote artifacts directory", "pushdir", remoteArtifactsDir)

	// Read top-level entries in artifact directory
	entries, err := os.ReadDir(artifactDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			dirPath := filepath.Join(artifactDir, entry.Name())
			h.logger.Debug("processing top-level directory", "pushdirentry", dirPath)

			childEntries, err := os.ReadDir(dirPath)
			if err != nil {
				return err
			}

			for _, child := range childEntries {
				if !child.IsDir() {
					info, _ := child.Info()
					h.logger.Debug("file size", "filesize", info.Size())
					localPath := filepath.Join(dirPath, child.Name())
					remotePath := driver.Join(remoteArtifactsDir, entry.Name(), child.Name())
					h.logger.Debug("pushing artifact file", "localPath", localPath, "remotePath", remotePath)
					if err := driver.Upload(ctx, localPath, remotePath); err != nil {
						return fmt.Errorf("failed to push artifact %s: %w", localPath, err)
					}
				}
			}
		}
	}

	return nil
}

// pullArtifactsWithDriver downloads all files from the remote artifacts directory to the local artifact directory
func (h *FlowExecutionHandler) pullArtifactsWithDriver(ctx context.Context, driver executor.NodeDriver, artifactDir string, execID string, nodeName string) error {
	remoteArtifactsDir := driver.Join(driver.TempDir(), fmt.Sprintf("artifacts-%s", execID))
	h.logger.Debug("remote artifacts directory", "pulldir", remoteArtifactsDir)
	files, err := driver.ListFiles(ctx, remoteArtifactsDir)
	if err != nil {
		// If the directory doesn't exist, there are no artifacts to pull
		h.logger.Debug("no artifacts to pull", "remoteDir", remoteArtifactsDir, "error", err)
		return nil
	}

	for _, file := range files {
		remotePath := driver.Join(remoteArtifactsDir, file)

		var localPath string
		if driver.IsRemote() {
			// Remote execution then store in nodeName subdirectory
			localPath = filepath.Join(artifactDir, nodeName, file)
		} else {
			// Local execution then store in local subdirectory
			localPath = filepath.Join(artifactDir, "local", file)
		}

		if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for artifact %s: %w", file, err)
		}

		if err := driver.Download(ctx, remotePath, localPath); err != nil {
			return fmt.Errorf("failed to pull artifact %s from node %s: %w", file, nodeName, err)
		}
	}
	return nil
}

func (h *FlowExecutionHandler) checkApproval(ctx context.Context, execID string, action Action, namespaceID string) error {
	namespaceUUID, err := uuid.Parse(namespaceID)
	if err != nil {
		return fmt.Errorf("invalid namespace UUID: %w", err)
	}

	// Set the current action ID
	h.logger.Debug("current action", "actionID", action.ID)
	if _, err := h.store.UpdateExecutionActionID(ctx, repo.UpdateExecutionActionIDParams{
		CurrentActionID: sql.NullString{String: action.ID, Valid: action.ID != ""},
		ExecID:          execID,
		Uuid:            namespaceUUID,
	}); err != nil {
		return fmt.Errorf("could not update current action ID in exec %s: %w", execID, err)
	}

	if !action.Approval {
		return nil
	}

	// check if pending approval, exit if not approved
	a, err := h.store.GetApprovalRequestForActionAndExec(ctx, repo.GetApprovalRequestForActionAndExecParams{
		ExecID:   execID,
		ActionID: action.ID,
		Uuid:     namespaceUUID,
	})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	// continue execution if approved
	if a.Status == repo.ApprovalStatusApproved {
		return nil
	}

	if a.Status == repo.ApprovalStatusRejected {
		return fmt.Errorf("request for running action %q is rejected", action.Name)
	}

	if a.Status == "" {
		_, err = h.store.RequestApprovalTx(ctx, execID, namespaceUUID, repo.RequestApprovalParam{
			ID: action.ID,
		})
		if err != nil {
			return err
		}
	}

	return ErrPendingApproval
}

// setStatus updates the execution status in the execution_log table
func (h *FlowExecutionHandler) setStatus(ctx context.Context, execID string, status repo.ExecutionStatus, namespaceID string, err error) error {
	var errMsg sql.NullString
	if err != nil {
		errMsg = sql.NullString{String: err.Error(), Valid: true}
	}
	namespaceUUID, parseErr := uuid.Parse(namespaceID)
	if parseErr != nil {
		return fmt.Errorf("invalid namespace ID: %w", parseErr)
	}
	_, updateErr := h.store.UpdateExecutionStatus(ctx, repo.UpdateExecutionStatusParams{
		Status: status,
		Error:  errMsg,
		ExecID: execID,
		Uuid:   namespaceUUID,
	})
	if updateErr != nil {
		return fmt.Errorf("could not update error execution status: %w", updateErr)
	}

	return nil
}

// createExecutionLog creates an execution log entry for scheduled jobs
func (h *FlowExecutionHandler) createExecutionLog(ctx context.Context, execID string, payload FlowExecutionPayload) error {
	namespaceUUID, err := uuid.Parse(payload.NamespaceID)
	if err != nil {
		return fmt.Errorf("invalid namespace UUID: %w", err)
	}

	userUUID, err := uuid.Parse(payload.UserUUID)
	if err != nil {
		return fmt.Errorf("invalid user UUID: %w", err)
	}

	inputJSON, err := json.Marshal(payload.Input)
	if err != nil {
		return fmt.Errorf("failed to marshal input: %w", err)
	}

	triggerType := repo.TriggerTypeManual
	if payload.TriggerType == TriggerTypeScheduled {
		triggerType = repo.TriggerTypeScheduled
	}

	_, err = h.store.AddExecutionLog(ctx, repo.AddExecutionLogParams{
		ExecID:      execID,
		FlowID:      payload.Workflow.Meta.DBID,
		Input:       inputJSON,
		TriggerType: triggerType,
		Uuid:        userUUID,
		Uuid_2:      namespaceUUID,
	})
	if err != nil {
		return fmt.Errorf("failed to add execution log: %w", err)
	}

	return nil
}

// setStatusWithMetrics updates the execution status and tracks metrics
func (h *FlowExecutionHandler) setStatusWithMetrics(ctx context.Context, execID string, status repo.ExecutionStatus, namespaceID string, flowID string, err error) error {
	if err := h.setStatus(ctx, execID, status, namespaceID, err); err != nil {
		return err
	}

	if h.metrics != nil {
		switch status {
		case repo.ExecutionStatusCompleted:
			h.metrics.IncrementExecutionCount(namespaceID, flowID, "completed")
		case repo.ExecutionStatusErrored:
			h.metrics.IncrementExecutionCount(namespaceID, flowID, "errored")
		case repo.ExecutionStatusCancelled:
			h.metrics.IncrementExecutionCount(namespaceID, flowID, "cancelled")
		case repo.ExecutionStatusPendingApproval:
			h.metrics.IncExecutionsWaiting(namespaceID, flowID)
		}
	}

	return nil
}
