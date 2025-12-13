package storage

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

// Job represents a job in the queue
type Job struct {
	ID          int64     `json:"id" db:"id"`
	ExecID      string    `json:"exec_id" db:"exec_id"`
	PayloadType string    `json:"payload_type" db:"payload_type"`
	Payload     []byte    `json:"payload" db:"payload"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

var (
	ErrNoJobs = errors.New("no jobs available")
)

// Storage interface for job queue storage backends
type Storage interface {
	// Initialize sets up the storage backend (creates tables, etc.)
	Initialize(ctx context.Context) error

	// Put adds a job to the queue
	Put(ctx context.Context, job Job) error

	// GetByPayloadType retrieves and locks a job of specific payload type from the queue
	// The job remains locked until the done channel is closed
	// Returns ErrNoJobs if no jobs are available
	GetByPayloadType(ctx context.Context, payloadType string, done chan struct{}) (Job, error)

	// Delete removes a job from the queue
	Delete(ctx context.Context, jobID int64) error

	// CancelByExecID removes all jobs with the given execution ID
	CancelByExecID(ctx context.Context, execID string) error

	// Close closes the storage backend
	Close() error
}

// NewJob creates a new job with the given execution ID, payload type, and payload
func NewJob(execID string, payloadType string, payload any) (Job, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return Job{}, err
	}

	return Job{
		ExecID:      execID,
		PayloadType: payloadType,
		Payload:     payloadBytes,
		CreatedAt:   time.Now(),
	}, nil
}
