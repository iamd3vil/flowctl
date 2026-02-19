package messengers

import (
	"bytes"
	"context"
	"crypto/ed25519"
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

type webhookPayload struct {
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`
	Data      any    `json:"data"`
}

// WebhookMessenger sends HTTP POST requests using the Standard Webhooks format.
type WebhookMessenger struct {
	privateKey ed25519.PrivateKey
	client     *http.Client
	logger     *slog.Logger
}

// NewWebhookMessenger creates a new WebhookMessenger with the given configuration.
func NewWebhookMessenger(cfg config.WebhookConfig, logger *slog.Logger) (*WebhookMessenger, error) {
	if !cfg.Enabled {
		return nil, fmt.Errorf("webhook messenger is disabled")
	}

	seedStr := strings.TrimPrefix(cfg.SigningKey, "whsk_")
	seed, err := base64.StdEncoding.DecodeString(seedStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode webhook signing key: %w", err)
	}

	privateKey := ed25519.NewKeyFromSeed(seed)

	timeout := 30 * time.Second
	if cfg.Timeout > 0 {
		timeout = cfg.Timeout
	}

	return &WebhookMessenger{
		privateKey: privateKey,
		client:     &http.Client{Timeout: timeout},
		logger:     logger,
	}, nil
}

// Send posts the message to the URL specified in msg.Config["url"] using Standard Webhooks headers.
func (w *WebhookMessenger) Send(_ context.Context, msg Message) error {
	targetURL, _ := msg.Config["url"].(string)
	if targetURL == "" {
		return fmt.Errorf("webhook messenger requires a url in config")
	}

	payload := webhookPayload{
		Type:      string(msg.Event),
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		Data:      msg.Data,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	msgID := "msg_" + uuid.New().String()
	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	toSign := fmt.Sprintf("%s.%s.%s", msgID, timestamp, string(payloadBytes))
	sig := ed25519.Sign(w.privateKey, []byte(toSign))
	signature := "v1a," + base64.StdEncoding.EncodeToString(sig)

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
