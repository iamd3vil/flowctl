package scheduler

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"strings"

	"github.com/cvhariharan/flowctl/internal/messengers"
	"github.com/cvhariharan/flowctl/internal/repo"
	"github.com/google/uuid"
)

//go:embed templates/*.html
var templateFS embed.FS

const PayloadTypeNotification PayloadType = "notification"

type NotificationPayload struct {
	FlowID      string   `json:"flow_id"`
	FlowName    string   `json:"flow_name"`
	ExecID      string   `json:"exec_id"`
	Status      string   `json:"status"`
	Error       string   `json:"error,omitempty"`
	Receivers   []string `json:"receivers"`
	NamespaceID string   `json:"namespace_id"`
	Channel     string   `json:"channel"`
}

// NotificationHandler processes notification jobs
type NotificationHandler struct {
	messengers map[string]messengers.Messenger
	store      repo.Store
	logger     *slog.Logger
	templates  *template.Template
	rootURL    string
}

func NewNotificationHandler(messengers map[string]messengers.Messenger, store repo.Store, logger *slog.Logger, rootURL string) (*NotificationHandler, error) {
	tmpl, err := template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse notification templates: %w", err)
	}

	return &NotificationHandler{
		messengers: messengers,
		store:      store,
		logger:     logger,
		templates:  tmpl,
		rootURL:    rootURL,
	}, nil
}

func (h *NotificationHandler) Type() PayloadType {
	return PayloadTypeNotification
}

func (h *NotificationHandler) Handle(ctx context.Context, job Job) error {
	var payload NotificationPayload
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		return fmt.Errorf("failed to unmarshal notification payload: %w", err)
	}

	h.logger.Debug("processing notification", "flow_id", payload.FlowID, "exec_id", payload.ExecID, "status", payload.Status, "channel", payload.Channel)

	// Route to messenger by channel name
	messenger, ok := h.messengers[payload.Channel]
	if !ok {
		h.logger.Warn("no messenger configured for channel", "channel", payload.Channel)
		return nil
	}

	// Resolve receivers to actual email addresses
	recipients, err := h.resolveReceivers(ctx, payload.Receivers)
	if err != nil {
		// don't return, continue with resolved receivers
		h.logger.Error("failed to resolve receivers", "error", err)
	}

	if len(recipients) == 0 {
		h.logger.Debug("no recipients to notify", "exec_id", payload.ExecID)
		return nil
	}

	namespaceUUID, err := uuid.Parse(payload.NamespaceID)
	if err != nil {
		return err
	}
	namespace, err := h.store.GetNamespaceByUUID(ctx, namespaceUUID)
	if err != nil {
		return fmt.Errorf("could not get namespace name for %s: %w", payload.NamespaceID, err)
	}

	msg := messengers.Message{
		Title:      h.buildSubject(payload),
		Body:       h.buildBody(payload, namespace.Name),
		Recipients: recipients,
	}

	if err := messenger.Send(msg); err != nil {
		return fmt.Errorf("failed to send notification via %s: %w", payload.Channel, err)
	}

	h.logger.Info("notification sent", "flow_id", payload.FlowID, "exec_id", payload.ExecID, "channel", payload.Channel, "recipient_count", len(recipients))

	return nil
}

// resolveReceivers resolves receiver strings to Recipient structs
// Formats: "email@example.com" or "group:groupname"
func (h *NotificationHandler) resolveReceivers(ctx context.Context, receivers []string) ([]messengers.Recipient, error) {
	var recipients []messengers.Recipient

	for _, r := range receivers {
		// Check if this is a group reference
		if groupName, isGroup := strings.CutPrefix(r, "group:"); isGroup {
			if groupName == "" {
				h.logger.Warn("empty group name", "receiver", r)
				continue
			}

			// Resolve group members from database
			members, err := h.store.GetGroupMembersByName(ctx, groupName)
			if err != nil {
				h.logger.Error("failed to get group members", "group", groupName, "error", err)
				continue
			}

			for _, member := range members {
				recipients = append(recipients, messengers.Recipient{
					UUID:  member.Uuid.String(),
					Email: member.Username,
				})
			}
		} else {
			// Treat as direct email address
			if r == "" {
				h.logger.Warn("empty receiver", "receiver", r)
				continue
			}
			recipients = append(recipients, messengers.Recipient{
				Email: r,
			})
		}
	}

	return recipients, nil
}

// buildSubject creates the email subject line
func (h *NotificationHandler) buildSubject(payload NotificationPayload) string {
	var status string
	switch payload.Status {
	case "completed":
		status = "[Success]"
	case "errored":
		status = "[Failed]"
	case "cancelled":
		status = "[Cancelled]"
	case "pending_approval":
		status = "[Waiting]"
	default:
		status = "[Update]"
	}

	return fmt.Sprintf("%s Flow %s - %s", status, payload.FlowName, payload.ExecID[:8])
}

// buildBody creates the email body
func (h *NotificationHandler) buildBody(payload NotificationPayload, namespace string) string {
	var statusMsg string
	switch payload.Status {
	case "completed":
		statusMsg = "has completed successfully"
	case "errored":
		statusMsg = "has failed with an error"
	case "cancelled":
		statusMsg = "was cancelled"
	case "pending_approval":
		statusMsg = "is waiting for approval"
	default:
		statusMsg = "status changed to " + payload.Status
	}

	data := struct {
		FlowName  string
		FlowID    string
		ExecID    string
		Status    string
		Namespace string
		StatusMsg string
		Error     string
		RootURL   string
	}{
		FlowName:  payload.FlowName,
		FlowID:    payload.FlowID,
		ExecID:    payload.ExecID,
		Status:    payload.Status,
		StatusMsg: statusMsg,
		Namespace: namespace,
		Error:     payload.Error,
		RootURL:   h.rootURL,
	}

	var buf bytes.Buffer
	if err := h.templates.ExecuteTemplate(&buf, "notification.html", data); err != nil {
		h.logger.Error("failed to execute template", "template", "notification.html", "error", err)
		// fallback to simple text message
		return fmt.Sprintf("Flow %s has %s", payload.FlowName, statusMsg)
	}

	return buf.String()
}
