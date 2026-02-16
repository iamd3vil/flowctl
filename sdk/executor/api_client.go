package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type TriggerFlowResponse struct {
	ExecID      string  `json:"exec_id"`
	ScheduledAt *string `json:"scheduled_at,omitempty"`
}

type FlowStatusResponse struct {
	ID              string          `json:"id"`
	FlowName        string          `json:"flow_name"`
	FlowID          string          `json:"flow_id"`
	Status          string          `json:"status"`
	TriggerType     string          `json:"trigger_type"`
	Input           json.RawMessage `json:"input,omitempty"`
	TriggeredBy     string          `json:"triggered_by"`
	CurrentActionID string          `json:"current_action_id"`
	CreatedAt       string          `json:"created_at"`
	StartedAt       string          `json:"started_at"`
	CompletedAt     string          `json:"completed_at"`
	ScheduledAt     string          `json:"scheduled_at,omitempty"`
}

// APIClient is an HTTP client for interacting with the server
type APIClient struct {
	baseURL    string
	apiKey     string
	userUUID   string
	httpClient *http.Client
}

func NewAPIClient(baseURL, apiKey, userUUID string) *APIClient {
	return &APIClient{
		baseURL:  strings.TrimRight(baseURL, "/"),
		apiKey:   apiKey,
		userUUID: userUUID,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// do executes an HTTP request with common auth headers, reads the response,
// and returns the body. Returns an error for non-200 status codes.
func (c *APIClient) do(ctx context.Context, method, path string, body io.Reader, contentType string) ([]byte, error) {
	endpoint := fmt.Sprintf("%s%s", c.baseURL, path)

	req, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	if c.userUUID != "" {
		req.Header.Set("X-User-UUID", c.userUUID)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// TriggerFlow triggers a flow execution via the HTTP API. Params are sent as form-encoded data
func (c *APIClient) TriggerFlow(ctx context.Context, namespace, flowID string, params map[string]any) (TriggerFlowResponse, error) {
	form := url.Values{}
	for k, v := range params {
		form.Set(k, fmt.Sprintf("%v", v))
	}

	path := fmt.Sprintf("/api/v1/%s/trigger/%s", namespace, flowID)
	body, err := c.do(ctx, http.MethodPost, path, strings.NewReader(form.Encode()), "application/x-www-form-urlencoded")
	if err != nil {
		return TriggerFlowResponse{}, fmt.Errorf("trigger flow: %w", err)
	}

	var result TriggerFlowResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return TriggerFlowResponse{}, fmt.Errorf("failed to decode trigger response: %w", err)
	}

	return result, nil
}

// GetFlowStatus retrieves the execution status of a flow.
func (c *APIClient) GetFlowStatus(ctx context.Context, namespace, execID string) (FlowStatusResponse, error) {
	path := fmt.Sprintf("/api/v1/%s/flows/executions/%s", namespace, execID)
	body, err := c.do(ctx, http.MethodGet, path, nil, "")
	if err != nil {
		return FlowStatusResponse{}, fmt.Errorf("get flow status: %w", err)
	}

	var result FlowStatusResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return FlowStatusResponse{}, fmt.Errorf("failed to decode status response: %w", err)
	}

	return result, nil
}
