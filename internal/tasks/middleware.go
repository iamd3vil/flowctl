package tasks

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cvhariharan/flowctl/internal/repo"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

type StatusTrackerDB struct {
	store repo.Store
}

func NewStatusTracker(s repo.Store) *StatusTrackerDB {
	return &StatusTrackerDB{store: s}
}

func (s *StatusTrackerDB) SetStatus(ctx context.Context, execID string, status repo.ExecutionStatus, namespaceID string, err error) error {
	var errMsg sql.NullString
	if err != nil {
		errMsg = sql.NullString{String: err.Error(), Valid: true}
	}
	namespaceUUID, parseErr := uuid.Parse(namespaceID)
	if parseErr != nil {
		return fmt.Errorf("invalid namespace ID: %w", parseErr)
	}
	_, err = s.store.UpdateExecutionStatus(ctx, repo.UpdateExecutionStatusParams{
		Status: status,
		Error:  errMsg,
		ExecID: execID,
		Uuid:   namespaceUUID,
	})
	if err != nil {
		return fmt.Errorf("could not update error execution status: %w", err)
	}

	return nil
}

// createExecutionLog creates an execution log entry for all task executions
// Uses the user UUID from the payload and sets appropriate trigger type
func (s *StatusTrackerDB) createExecutionLog(ctx context.Context, payload FlowExecutionPayload) error {
	namespaceUUID, err := uuid.Parse(payload.NamespaceID)
	if err != nil {
		return fmt.Errorf("invalid namespace UUID: %w", err)
	}

	userUUID, err := uuid.Parse(payload.UserUUID)
	if err != nil {
		return fmt.Errorf("invalid user UUID: %w", err)
	}

	inputB, err := json.Marshal(payload.Input)
	if err != nil {
		return fmt.Errorf("could not marshal input to json: %w", err)
	}

	// Default to manual trigger type
	triggerType := repo.TriggerTypeManual
	if payload.TriggerType == TriggerTypeScheduled {
		triggerType = repo.TriggerTypeScheduled
	}

	_, err = s.store.AddExecutionLog(ctx, repo.AddExecutionLogParams{
		ExecID:      payload.ExecID,
		FlowID:      payload.Workflow.Meta.DBID,
		Input:       inputB,
		TriggerType: triggerType,
		Uuid:        userUUID,
		Uuid_2:      namespaceUUID,
	})
	if err != nil {
		return fmt.Errorf("could not add execution log entry: %w", err)
	}

	return nil
}

func (s *StatusTrackerDB) TrackerMiddleware(next func(context.Context, *asynq.Task) error) func(context.Context, *asynq.Task) error {
	return func(ctx context.Context, t *asynq.Task) error {
		var payload FlowExecutionPayload

		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			return fmt.Errorf("payload could not be deserialized: %w", err)
		}

		// Create execution log for scheduled executions only (manual ones are created in core)
		if payload.TriggerType == TriggerTypeScheduled {
			if err := s.createExecutionLog(ctx, payload); err != nil {
				return fmt.Errorf("failed to create execution log for scheduled task: %w", err)
			}
		}

		if err := s.SetStatus(ctx, payload.ExecID, repo.ExecutionStatusRunning, payload.NamespaceID, nil); err != nil {
			return err
		}

		if err := next(ctx, t); err != nil {
			if errors.Is(err, ErrPendingApproval) {
				return s.SetStatus(ctx, payload.ExecID, repo.ExecutionStatusPendingApproval, payload.NamespaceID, nil)
			}
			if errors.Is(err, ErrExecutionCancelled) {
				return s.SetStatus(ctx, payload.ExecID, repo.ExecutionStatusCancelled, payload.NamespaceID, nil)
			}
			return s.SetStatus(ctx, payload.ExecID, repo.ExecutionStatusErrored, payload.NamespaceID, err)
		}

		return s.SetStatus(ctx, payload.ExecID, repo.ExecutionStatusCompleted, payload.NamespaceID, nil)
	}
}
