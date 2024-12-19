package tasks

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/cvhariharan/autopilot/internal/repo"
	"github.com/hibiken/asynq"
)

type StatusTrackerDB struct {
	store repo.Store
}

func NewStatusTracker(s repo.Store) *StatusTrackerDB {
	return &StatusTrackerDB{store: s}
}

func (s *StatusTrackerDB) SetStatus(ctx context.Context, execID string, res map[string]string, err error) error {
	if err != nil {
		_, err = s.store.UpdateExecutionStatus(ctx, repo.UpdateExecutionStatusParams{
			Status: repo.ExecutionStatusErrored,
			Error:  sql.NullString{String: err.Error(), Valid: true},
			ExecID: execID,
		})
		if err != nil {
			return fmt.Errorf("could not update error execution status: %w", err)
		}
	} else {
		resB, err := json.Marshal(res)
		if err != nil {
			return fmt.Errorf("could not marshal results for status tracking: %w", err)
		}

		_, err = s.store.UpdateExecutionStatus(ctx, repo.UpdateExecutionStatusParams{
			Status: repo.ExecutionStatusErrored,
			Output: resB,
			ExecID: execID,
		})
		if err != nil {
			return fmt.Errorf("could not update completed execution status: %w", err)
		}
	}

	return nil
}

func (s *StatusTrackerDB) TrackerMiddleware(next func(context.Context, *asynq.Task) error) func(context.Context, *asynq.Task) error {
	return func(ctx context.Context, t *asynq.Task) error {

	}
}
