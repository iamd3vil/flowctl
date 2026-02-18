package messengers

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/cvhariharan/flowctl/internal/config"
	"github.com/google/uuid"
	"github.com/invopop/jsonschema"
)

// WebhookNotifyConfig defines the per-flow webhook configuration rendered in the UI.
type WebhookNotifyConfig struct {
	URL string `json:"url" jsonschema:"title=Webhook URL,description=URL to POST webhook notifications to"`
}

func GetWebhookNotifySchema() interface{} {
	return jsonschema.Reflect(&WebhookNotifyConfig{})
}

// WebhookPayload is the JSON body sent to the webhook endpoint.
type WebhookPayload struct {
	Event     string `json:"event"`
	FlowName  string `json:"flow_name"`
	FlowID    string `json:"flow_id"`
	ExecID    string `json:"exec_id"`
	Status    string `json:"status"`
	Error     string `json:"error,omitempty"`
	Namespace string `json:"namespace"`
}

// WebhookMessenger sends HTTP POST requests using the Standard Webhooks format.
type WebhookMessenger struct {
	secret []byte
	client *http.Client
	logger *slog.Logger
}

// NewWebhookMessenger creates a new WebhookMessenger with the given configuration.
func NewWebhookMessenger(cfg config.WebhookConfig, logger *slog.Logger) (*WebhookMessenger, error) {
	if !cfg.Enabled {
		return nil, fmt.Errorf("webhook messenger is disabled")
	}

	secretStr := strings.TrimPrefix(cfg.Secret, "whsec_")
	secretBytes, err := base64.StdEncoding.DecodeString(secretStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode webhook secret: %w", err)
	}

	timeout := 30 * time.Second
	if cfg.Timeout > 0 {
		timeout = cfg.Timeout
	}

	return &WebhookMessenger{
		secret: secretBytes,
		client: &http.Client{Timeout: timeout},
		logger: logger,
	}, nil
}

// Send posts the message to the URL specified in msg.Config["url"] using Standard Webhooks headers.
func (w *WebhookMessenger) Send(_ context.Context, msg Message) error {
	targetURL, _ := msg.Config["url"].(string)
	if targetURL == "" {
		return fmt.Errorf("webhook messenger requires a url in config")
	}

	evt, ok := msg.Data.(FlowExecutionEvent)
	if !ok {
		return fmt.Errorf("webhook messenger: unsupported event data type %T", msg.Data)
	}

	payload := WebhookPayload{
		Event:     string(msg.Event),
		FlowName:  evt.FlowName,
		FlowID:    evt.FlowID,
		ExecID:    evt.ExecID,
		Status:    evt.Status,
		Error:     evt.Error,
		Namespace: evt.Namespace,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	msgID := "msg_" + uuid.New().String()
	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	toSign := fmt.Sprintf("%s.%s.%s", msgID, timestamp, string(payloadBytes))
	mac := hmac.New(sha256.New, w.secret)
	mac.Write([]byte(toSign))
	signature := "v1," + base64.StdEncoding.EncodeToString(mac.Sum(nil))

	req, err := http.NewRequest(http.MethodPost, targetURL, bytes.NewReader(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create webhook request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("webhook-id", msgID)
	req.Header.Set("webhook-timestamp", timestamp)
	req.Header.Set("webhook-signature", signature)

	resp, err := w.client.Do(req)
	if err != nil {
		w.logger.Error("failed to send webhook", "url", targetURL, "error", err)
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		w.logger.Error("webhook returned non-2xx status", "url", targetURL, "status", resp.StatusCode)
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	w.logger.Debug("webhook sent", "url", targetURL, "msg_id", msgID, "event", msg.Event)
	return nil
}

// Close is a no-op for the webhook messenger.
func (w *WebhookMessenger) Close() {}
